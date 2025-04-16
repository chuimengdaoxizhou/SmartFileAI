package models

import (
	"RAG/utils"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"userID" gorm:"unique"`
	Password  string    `json:"password"`
	Nickname  string    `json:"nickname"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var DB *gorm.DB

func InitDB() {
	dsn := utils.Config.GetString("database.user") + ":" +
		utils.Config.GetString("database.password") + "@tcp(" +
		utils.Config.GetString("database.host") + ":" +
		utils.Config.GetString("database.port") + ")/" +
		utils.Config.GetString("database.name") + "?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.Logger.Fatal("Failed to connect to mysql:", err)
	}

	utils.Logger.Println("Connected to mysql")
	// Auto migrate
	if err := DB.AutoMigrate(&User{}); err != nil {
		utils.Logger.Fatal("Failed to migrate mysql:", err)
	}
}

// CloseDB 关闭数据库连接
func CloseDB() {
	// Retrieve the underlying *sql.DB object to close the connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		utils.Logger.Fatalf("Failed to get mysql instance: %v", err)
		return
	}

	// Close the connection pool gracefully
	err = sqlDB.Close()
	if err != nil {
		utils.Logger.Fatalf("Failed to close mysql connection: %v", err)
		return
	}

	utils.Logger.Println("mysql connection closed.")
}
