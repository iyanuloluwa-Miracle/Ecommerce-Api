package repositories

import (
	"context"

	"github.com/iyanuloluwa-Miracle/ecommerce-api/config"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateUser(user models.User) error {
	_, err := config.DB.Collection("users").InsertOne(context.TODO(), user)
	return err
}

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Collection("users").
		FindOne(context.TODO(), bson.M{"email": email}).
		Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
