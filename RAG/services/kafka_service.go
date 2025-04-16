package services

import (
	"RAG/utils"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

var KafkaProducer *kafka.Producer

// InitProducer 初始化 Kafka 生产者
func InitProducer() {
	fmt.Println("获取 Kafka 地址")
	//addr := utils.Config.Get("Kafka.address")
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Kafka 集群地址
	}
	var err error
	fmt.Println("创建 Kafka 生产者")
	KafkaProducer, err = kafka.NewProducer(config)
	if err != nil {
		utils.Logger.Fatalf("Failed to create kafka producer: %v", err)
	}
	utils.Logger.Println("Successed to create kafka producer")
}

type data struct {
	FilePath string
	Name     string
}
type data2 struct {
	bucket   string
	filename string
}

// SendMessage 发送消息
func SendMessage(topic, userid, fileName string) {
	// 创建消息结构体
	message := data2{
		bucket:   userid,
		filename: fileName,
	}

	// 将结构体转换为 JSON 字符串
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	// 发送消息到 Kafka
	err = KafkaProducer.Produce(&kafka.Message{
		//TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: 0},
		Value:          messageBytes,
	}, nil)

	if err != nil {
		log.Fatalf("Failed to produce message: %v", err)
	}

}

// CloseProducer 关闭 Kafka 生产者
func CloseProducer() {
	// 调用 Kafka 生产者的 Close 方法来优雅地关闭生产者
	KafkaProducer.Close()
	utils.Logger.Println("Kafka producer closed successfully.")
}
