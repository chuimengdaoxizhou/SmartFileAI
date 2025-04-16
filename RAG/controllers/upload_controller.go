package controllers

import (
	"RAG/minio_client"
	"RAG/utils"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// UploadRequest 上传请求结构体
type UploadRequest struct {
	FileName string `form:"fileName" binding:"required"` // 文件名
	Chunk    int    `form:"chunk" binding:"required"`    // 当前分片索引
	Total    int    `form:"total" binding:"required"`    // 总分片数
	MD5      string `form:"md5" binding:"required"`      // 整个文件的 MD5
	ChunkMD5 string `form:"chunkMD5" binding:"required"` // 当前分片的 MD5
}

// ResumableUpload 断点续传上传处理
func ResumableUpload(c *gin.Context) {
	var req UploadRequest
	// 获取 POST 表单中的所有参数
	req.FileName = c.PostForm("fileName")
	chunk, _ := strconv.Atoi(c.PostForm("chunk"))
	req.Chunk = chunk

	total, _ := strconv.Atoi(c.PostForm("total"))
	req.Total = total

	req.MD5 = c.PostForm("md5")
	req.ChunkMD5 = c.PostForm("chunkMD5")

	// 绑定其他字段
	//if err := c.ShouldBind(&req); err != nil {
	//	utils.Logger.Printf("绑定参数错误 (Binding error): %v", err)
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error":   "请求参数错误",
	//		"details": err.Error(),
	//	})
	//	return
	//}

	utils.Logger.Println("开始上传文件:", req.FileName, "总分片数:", req.Total)
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "上传文件失败"})
		return
	}

	// 打开上传的文件
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败"})
		return
	}
	defer f.Close()

	// 创建临时文件用于存储分片，同时计算 MD5
	tempDir := fmt.Sprintf("temp/%s", req.MD5)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建临时目录失败"})
		return
	}
	tempFilePath := fmt.Sprintf("%s/part%d", tempDir, req.Chunk)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建临时分片文件失败"})
		return
	}
	defer tempFile.Close()

	// 计算分片 MD5 并保存到临时文件
	hash := md5.New()
	multiWriter := io.MultiWriter(tempFile, hash)
	if _, err := io.Copy(multiWriter, f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存分片或计算 MD5 失败"})
		return
	}
	chunkMD5 := hex.EncodeToString(hash.Sum(nil))

	// 验证分片 MD5
	if chunkMD5 != req.ChunkMD5 {
		// 如果 MD5 不匹配，删除临时文件
		os.Remove(tempFilePath)
		c.JSON(http.StatusBadRequest, gin.H{"error": "分片 MD5 验证失败"})
		return
	}

	// 检查是否为最后一个分片
	if req.Chunk == req.Total-1 {
		// 合并所有分片
		finalFilePath := fmt.Sprintf("uploads/%s%s", req.MD5, filepath.Ext(req.FileName))
		finalFile, err := os.Create(finalFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建最终文件失败"})
			return
		}
		defer finalFile.Close()

		// 计算合并后文件的 MD5
		hash = md5.New()
		multiWriter = io.MultiWriter(finalFile, hash)

		// 按顺序合并分片
		for i := 0; i < req.Total; i++ {
			partPath := fmt.Sprintf("%s/part%d", tempDir, i)
			partFile, err := os.Open(partPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("打开分片 %d 失败", i)})
				return
			}
			if _, err := io.Copy(multiWriter, partFile); err != nil {
				partFile.Close()
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("合并分片 %d 失败", i)})
				return
			}
			partFile.Close()
		}

		// 验证合并后文件的 MD5
		finalMD5 := hex.EncodeToString(hash.Sum(nil))
		if finalMD5 != req.MD5 {
			os.Remove(finalFilePath)
			c.JSON(http.StatusBadRequest, gin.H{"error": "合并文件 MD5 验证失败"})
			return
		}

		// 上传合并后的文件到 Minio
		ctx := context.Background()
		_, err = minio_client.MinioClient.FPutObject(
			ctx,
			minio_client.BucketName,
			fmt.Sprintf("%s%s", req.MD5, filepath.Ext(req.FileName)),
			finalFilePath,
			minio.PutObjectOptions{},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "上传到 Minio 失败"})
			return
		}

		// 清理临时文件和目录
		if err := os.RemoveAll(tempDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "清理临时文件失败"})
			return
		}
		if err := os.Remove(finalFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "清理最终文件失败"})
			return
		}
		utils.Logger.Println("文件上传成功:", req.FileName)

		c.JSON(http.StatusOK, gin.H{"message": "文件上传成功"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("分片 %d 上传成功", req.Chunk)})
}

// CheckChunk 检查分片状态
func CheckChunk(c *gin.Context) {
	md5 := c.Query("md5")
	chunkStr := c.Query("chunk")
	if md5 == "" || chunkStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 md5 或 chunk 参数"})
		return
	}

	chunk, err := strconv.Atoi(chunkStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chunk 参数无效"})
		return
	}

	// 检查分片文件是否存在
	tempFilePath := fmt.Sprintf("temp/%s/part%d", md5, chunk)
	if _, err := os.Stat(tempFilePath); err == nil {
		c.JSON(http.StatusOK, gin.H{"exists": true})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": false})
}
