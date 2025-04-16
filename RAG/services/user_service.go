package services

import (
	"RAG/models"
	"RAG/utils"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Register(userID, password, nickname string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		UserID:   userID,
		Password: string(hashedPassword),
		Nickname: nickname,
	}

	result := models.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	utils.Logger.Printf("User registered: %s", userID)
	return nil
}

func Login(userID, password string) (*models.User, error) {
	var user models.User
	result := models.DB.Where("user_id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	if err := AddOnlineUser(userID); err != nil {
		return nil, err
	}

	utils.Logger.Printf("User logged in: %s", userID)

	users, _ := GetOnlineUsers()
	fmt.Println("Online users:", users)
	return &user, nil
}

func Logout(userID string) error {
	if err := RemoveOnlineUser(userID); err != nil {
		return err
	}
	utils.Logger.Printf("User logged out: %s", userID)
	users, _ := GetOnlineUsers()
	fmt.Println("Online users:", users)
	return nil
}

func DeleteAccount(userID string) error {
	result := models.DB.Where("user_id = ?", userID).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}
	if err := RemoveOnlineUser(userID); err != nil {
		return err
	}
	utils.Logger.Printf("User account deleted: %s", userID)
	return nil
}
