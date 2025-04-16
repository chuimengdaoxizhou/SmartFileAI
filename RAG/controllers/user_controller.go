package controllers

import (
	"RAG/services"
	"RAG/utils"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	UserID   string `json:"userID" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type LoginRequest struct {
	UserID   string `json:"userID" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request")
		return
	}

	if err := services.Register(req.UserID, req.Password, req.Nickname); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, "User registered successfully")
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request")
		return
	}

	user, err := services.Login(req.UserID, req.Password)
	if err != nil {
		utils.ErrorResponse(c, 401, err.Error())
		return
	}

	token, err := services.GenerateJWT(user.UserID)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to generate token")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func Logout(c *gin.Context) {
	userID := c.GetString("userID")
	if err := services.Logout(userID); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, "Logged out successfully")
}

func DeleteAccount(c *gin.Context) {
	userID := c.GetString("userID")
	if err := services.DeleteAccount(userID); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, "Account deleted successfully")
}

func GetOnlineUsers(c *gin.Context) {
	users, err := services.GetOnlineUsers()
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, users)
}
