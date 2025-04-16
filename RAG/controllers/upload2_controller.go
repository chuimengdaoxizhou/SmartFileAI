package controllers

import (
	"RAG/services"
	"RAG/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const uploadPath = "./uploads" // Directory to save uploaded files
const maxUploadSize = 10 << 20 // 10 MB

// 文件上传处理函数
func UploadFileHandler(c *gin.Context) {
	// 获取用户名
	userid := c.GetString("userID")
	// 限制请求体的大小
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	// 解析 multipart 表单
	file, handler, err := c.Request.FormFile("myFile")
	if err != nil {
		utils.Logger.Printf("获取文件时出错: %v", err)
		c.JSON(400, gin.H{
			"error": "获取文件时出错",
		})
		return
	}
	defer file.Close() // 确保文件关闭

	// 获取上传文件的文件名
	fileName := handler.Filename // 使用 FormFile 返回的文件信息中的文件名
	newfileName := GetNewFileName(fileName)
	utils.Logger.Printf("接收到文件: %s", fileName)
	utils.Logger.Printf("将 %s 修改文件名为：%s", fileName, newfileName)

	// 如果上传目录不存在，则创建它
	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		utils.Logger.Printf("创建上传目录时出错: %v", err)
		c.JSON(500, gin.H{
			"error": "无法创建上传目录",
		})
		return
	}

	// 在上传目录中创建一个新文件
	dstPath := filepath.Join(uploadPath, newfileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		utils.Logger.Printf("创建文件时出错: %v", err)
		c.JSON(500, gin.H{
			"error": "无法保存文件",
		})
		return
	}
	defer dst.Close()

	// 将上传的文件内容复制到目标文件中
	_, err = io.Copy(dst, file)
	if err != nil {
		utils.Logger.Printf("保存文件内容时出错: %v", err)
		c.JSON(500, gin.H{
			"error": "保存文件内容时出错",
		})
		return
	}
	utils.Logger.Printf("文件 '%s' 上传成功！", fileName)

	err = services.UploadFileToMinIO(dstPath, userid, newfileName)
	if err != nil {
		utils.Logger.Printf(" %s 上传的文件 %s 存储到 minio 失败", userid, fileName)
	}
	utils.Logger.Printf(" %s 上传的文件 %s 存储到 minio 成功", userid, fileName)
	// 将文件名传递给 kafka
	topic := utils.Config.GetString("topic")
	services.SendMessage(topic, dstPath, newfileName)
	// 返回成功响应（JSON 格式）
	c.JSON(200, gin.H{
		"message": fmt.Sprintf("文件 '%s' 上传成功！", fileName),
	})
}

// GetFileName 从文件路径中提取文件名
func GetFileName(filePath string) string {
	// 使用 filepath.Base 获取文件名
	return filepath.Base(filePath)
}

// GetNewFileName 在原文件名的基础上追加当前时间
func GetNewFileName(filePath string) string {
	// 提取文件的目录和文件名
	dir := filepath.Dir(filePath)
	ext := filepath.Ext(filePath) // 提取文件扩展名

	// 获取当前时间，格式为 YYYYMMDD_HHMMSS
	currentTime := time.Now().Format("20060102_150405")

	// 提取原始文件名（去除扩展名）
	baseName := filepath.Base(filePath)
	fileNameWithoutExt := baseName[:len(baseName)-len(ext)]

	// 创建新的文件名（包括扩展名）
	newFileName := fileNameWithoutExt + "_" + currentTime + ext

	// 组合新的文件路径
	return filepath.Join(dir, newFileName)
}
