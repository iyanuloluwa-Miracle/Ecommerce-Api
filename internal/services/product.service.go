package services

import (
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProduct(p models.Product) error {
	return repositories.CreateProduct(p)
}

func ListProducts() ([]models.Product, error) {
	return repositories.GetAllProducts()
}

func GetProductByID(id primitive.ObjectID) (models.Product, error) {
	return repositories.GetProductByID(id)
}

func UpdateProduct(id primitive.ObjectID, name, description string, price float64, stock int) error {
	return repositories.UpdateProduct(id, name, description, price, stock)
}

func DeleteProduct(id primitive.ObjectID) error {
	return repositories.DeleteProduct(id)
}
