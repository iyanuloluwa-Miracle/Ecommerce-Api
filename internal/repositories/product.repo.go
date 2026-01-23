package repositories

import (
	"context"

	"github.com/iyanuloluwa-Miracle/ecommerce-api/config"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateProduct(p models.Product) error {
	_, err := config.DB.Collection("products").InsertOne(context.TODO(), p)
	return err
}

func GetAllProducts() ([]models.Product, error) {
	cur, err := config.DB.Collection("products").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var products []models.Product
	if err := cur.All(context.TODO(), &products); err != nil {
		return nil, err
	}

	if products == nil {
		products = []models.Product{}
	}

	return products, nil
}

func GetProductByID(id primitive.ObjectID) (models.Product, error) {
	var p models.Product
	err := config.DB.Collection("products").
		FindOne(context.TODO(), bson.M{"_id": id}).Decode(&p)
	return p, err
}

func UpdateProduct(id primitive.ObjectID, name, description string, price float64, stock int) error {
	update := bson.M{}

	if name != "" {
		update["name"] = name
	}
	if description != "" {
		update["description"] = description
	}
	if price > 0 {
		update["price"] = price
	}
	if stock >= 0 {
		update["stock"] = stock
	}

	if len(update) == 0 {
		return nil
	}

	_, err := config.DB.Collection("products").UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": update},
	)
	return err
}

func DeleteProduct(id primitive.ObjectID) error {
	_, err := config.DB.Collection("products").DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func DecrementProductStock(productID primitive.ObjectID, quantity int) error {
	_, err := config.DB.Collection("products").UpdateOne(
		context.TODO(),
		bson.M{"_id": productID},
		bson.M{"$inc": bson.M{"stock": -quantity}},
		options.Update().SetUpsert(false),
	)
	return err
}
