package services

import (
	"errors"

	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddToCart adds an item to the user's cart
func AddToCart(userID primitive.ObjectID, item models.CartItem) error {
	// Verify product exists
	product, err := repositories.GetProductByID(item.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	// Check stock availability
	if product.Stock < item.Quantity {
		return errors.New("insufficient stock")
	}

	cart, err := repositories.GetCart(userID)
	if err != nil {
		// Create new cart if doesn't exist
		cart = models.Cart{
			UserID: userID,
			Items:  []models.CartItem{},
		}
	}

	// Check if product already in cart, update quantity if so
	found := false
	for i, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			cart.Items[i].Quantity += item.Quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, item)
	}

	return repositories.SaveCart(cart)
}

// RemoveFromCart removes an item from the user's cart
func RemoveFromCart(userID primitive.ObjectID, productID primitive.ObjectID) error {
	cart, err := repositories.GetCart(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	// Find and remove the item
	newItems := []models.CartItem{}
	found := false
	for _, item := range cart.Items {
		if item.ProductID != productID {
			newItems = append(newItems, item)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("item not in cart")
	}

	cart.Items = newItems
	return repositories.SaveCart(cart)
}

// GetCart returns the user's cart
func GetCart(userID primitive.ObjectID) (models.Cart, error) {
	return repositories.GetCart(userID)
}
