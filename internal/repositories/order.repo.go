package repositories

import (
	"context"

	"github.com/iyanuloluwa-Miracle/ecommerce-api/config"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateOrder(order models.Order) error {
	_, err := config.DB.Collection("orders").InsertOne(context.TODO(), order)
	return err
}

func GetOrdersByUserID(userID primitive.ObjectID) ([]models.Order, error) {
	cur, err := config.DB.Collection("orders").Find(context.TODO(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var orders []models.Order
	if err := cur.All(context.TODO(), &orders); err != nil {
		return nil, err
	}

	if orders == nil {
		orders = []models.Order{}
	}

	return orders, nil
}

func GetOrderBySessionID(sessionID string) (*models.Order, error) {
	var order models.Order
	err := config.DB.Collection("orders").
		FindOne(context.TODO(), bson.M{"stripe_session": sessionID}).
		Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func UpdateOrderPaymentStatus(sessionID string, status string) error {
	_, err := config.DB.Collection("orders").UpdateOne(
		context.TODO(),
		bson.M{"stripe_session": sessionID},
		bson.M{"$set": bson.M{"payment_status": status}},
	)
	return err
}
