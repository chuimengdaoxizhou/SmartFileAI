package main

import (
	"RAG/models"
	"RAG/routes"
	"RAG/services"
	"RAG/utils"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//func main() {
//
//	// Initialize configuration
//	utils.InitConfig()
//	// Initialize logger
//	utils.InitLogger()
//	fmt.Println("Logger initialized")
//
//	prompt := "哪些学术成果可以申请激励？"
//	an := allgrpc.Getdata(prompt)
//	fmt.Println(an)
//
//	//services.InitProducer()
//	//fmt.Println("Kafka producer initialized")
//	//services.SendMessage("file_info_topic", "/home/chenyun/下载/train.json", "chenyun")
//	//fmt.Println("Message sent to Kafka topic")
//	//services.CloseProducer()
//	//fmt.Println("Kafka producer closed")
//	//
//	//prompt1 := "哪些学术成果可以申请激励？"
//	//an1 := allgrpc.Getdata(prompt1)
//	//fmt.Println(an1)
//}

//func main() {
//
//	question := "/home/chenyun/下载/train.json"
//	result := allgrpc.Updata(question)
//	fmt.Println(result)
//
//	prompt := "哪些学术成果可以申请激励？"
//	an := allgrpc.Getdata(prompt)
//	fmt.Println(an)
//}

func main() {
	// Initialize configuration
	utils.InitConfig()

	// Initialize logger
	utils.InitLogger()

	// Initialize database
	//models.InitDB()

	// Initialize Redis
	//services.InitRedis()
	//services.InitProducer()
	// Setup router
	router := routes.SetupRouter()
	port := utils.Config.GetString("server.port")

	// Setup HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine to allow graceful shutdown
	go func() {
		utils.Logger.Println("Starting server on port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	// Graceful shutdown logic
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create a deadline to wait for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	utils.Logger.Println("Shutting down gracefully...")

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		utils.Logger.Fatalf("Server shutdown failed: %v", err)
	}

	// Close resources
	models.CloseDB()
	services.CloseRedisClient()

	utils.Logger.Println("Server gracefully stopped")
}
