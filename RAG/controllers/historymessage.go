package controllers

import (
	"RAG/models"
	"RAG/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var log = utils.Logger

func StoreHistoryToMongo(userHistory *models.UserHistoryMessage) error {
	// 连接 MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // 替换为你的 MongoDB URI
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("连接 MongoDB 失败: %v", err)
	}
	defer client.Disconnect(context.Background())

	// 动态生成集合名称，例如：user_history_user123_1234567890
	collectionName := fmt.Sprintf("user_history_%s_%d", userHistory.UserID, userHistory.CreateTime)
	collection := client.Database("chat_db").Collection(collectionName)

	// 将 UserHistoryMessage 转换为 BSON 格式并插入到 MongoDB
	_, err = collection.InsertOne(context.Background(), userHistory)
	if err != nil {
		return fmt.Errorf("插入历史记录到 MongoDB 失败: %v", err)
	}

	log.Printf("历史记录成功存储到集合 %s", collectionName)
	return nil
}

// 根据 UserID 和 CreateTime 读取历史记录
func GetHistoryFromMongo(userID string, createTime int64) (*models.UserHistoryMessage, error) {
	// 连接 MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // 替换为你的 MongoDB URI
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("连接 MongoDB 失败: %v", err)
	}
	defer client.Disconnect(context.Background())

	// 动态生成集合名称
	collectionName := fmt.Sprintf("user_history_%s_%d", userID, createTime)
	collection := client.Database("chat_db").Collection(collectionName)

	// 查询历史记录
	var result models.UserHistoryMessage
	err = collection.FindOne(context.Background(), bson.M{}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("未找到用户 %s 在时间 %d 的历史记录", userID, createTime)
		}
		return nil, fmt.Errorf("查询历史记录失败: %v", err)
	}

	log.Printf("成功查找历史记录，集合名称为 %s", collectionName)
	return &result, nil
}
