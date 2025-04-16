package services

import (
	"RAG/utils"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

//   sudo docker run -d -p 9000:9000 -e "MINIO_ALLOW_ANONYMOUS_ACCESS=yes" 400c20c8aac0 server /data

// MinioClient UploadFileToMinIO 上传本地文件到 MinIO
// 参数:
// - filePath: 本地文件路径
// - minioClient: MinIO 客户端实例
// - bucketName: MinIO 存储桶名称
// - objectName: 存储在 MinIO 中的对象名称

var MinioClient *minio.Client

func InitMinioClient() {
	endpoint := utils.Config.GetString("endpoint")
	accessKeyID := utils.Config.GetString("access_key")
	secretAccessKey := utils.Config.GetString("secret_key")
	useSSL := utils.Config.GetBool("use_ssl")

	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}
	utils.Logger.Println("InitMinioClient successed")
}
func UploadFileToMinIO(filePath string, bucketName, objectName string) error {
	// 打开本地文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	// 获取文件信息以确定大小
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %v", filePath, err)
	}

	// 上传文件到 MinIO
	_, err = MinioClient.PutObject(context.Background(), bucketName, objectName, file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "application/octet-stream", // 默认 Content-Type，可根据需要修改
	})
	if err != nil {
		return fmt.Errorf("failed to upload file %s to MinIO: %v", filePath, err)
	}

	return nil
}

func mainn() {
	// 初始化 MinIO 客户端
	endpoint := "localhost:9000"
	accessKeyID := "your-access-key"
	secretAccessKey := "your-secret-key"
	useSSL := false
	bucketName := "my-bucket"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 示例：上传文件
	filePath := "/path/to/file.txt" // 替换为实际文件路径
	objectName := "file.txt"        // 存储在 MinIO 中的对象名称

	err = UploadFileToMinIO(filePath, bucketName, objectName)
	if err != nil {
		log.Printf("Upload failed: %v", err)
		return
	}
	// 列出桶中的所有对象
	objectCh := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Recursive: true, // 是否递归列出所有对象
	})
	fmt.Println(objectCh)
	// 使用 StatObject 查找文件是否存在
	_, err = minioClient.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			// 文件不存在
			fmt.Printf("文件 %s 不存在。\n", objectName)
		} else {
			// 其他错误
			log.Fatalf("查找文件时出错: %v", err)
		}
	} else {
		// 文件存在
		fmt.Printf("文件 %s 存在。\n", objectName)
	}
	log.Printf("File %s uploaded successfully to %s/%s", filePath, bucketName, objectName)
}
