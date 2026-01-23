package repositories

import (
	"context"

	"github.com/iyanuloluwa-Miracle/ecommerce-api/config"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCart(userID primitive.ObjectID) (models.Cart, error) {
	var cart models.Cart
	err := config.DB.Collection("carts").
		FindOne(context.TODO(), bson.M{"user_id": userID}).
		Decode(&cart)

	if err != nil {
		// Return empty cart if not found
		return models.Cart{
			UserID: userID,
			Items:  []models.CartItem{},
		}, err
	}

	if cart.Items == nil {
		cart.Items = []models.CartItem{}
	}

	return cart, nil
}

func SaveCart(cart models.Cart) error {
	opts := options.Update().SetUpsert(true)
	_, err := config.DB.Collection("carts").UpdateOne(
		context.TODO(),
		bson.M{"user_id": cart.UserID},
		bson.M{"$set": bson.M{
			"user_id": cart.UserID,
			"items":   cart.Items,
		}},
		opts,
	)
	return err
}

func ClearCart(userID primitive.ObjectID) error {
	_, err := config.DB.Collection("carts").UpdateOne(
		context.TODO(),
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"items": []models.CartItem{}}},
	)
	return err
}
