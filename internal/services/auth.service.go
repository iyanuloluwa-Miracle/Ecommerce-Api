package services

import (
	"errors"
	"time"

	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/repositories"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/utils"
)

func Register(email, password string) error {
	// Check if user already exists
	existingUser, _ := repositories.FindUserByEmail(email)
	if existingUser != nil && existingUser.Email != "" {
		return errors.New("user already exists")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := models.User{
		Email:     email,
		Password:  hash,
		Role:      "USER",
		CreatedAt: time.Now(),
	}
	return repositories.CreateUser(user)
}

func Login(email, password string) (string, error) {
	user, err := repositories.FindUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateJWT(user.ID.Hex(), user.Role)
}
