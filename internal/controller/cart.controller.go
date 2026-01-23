package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/services"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddToCartRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type RemoveFromCartRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}

// AddToCart adds an item to the user's cart
func AddToCart(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid user ID")
		return
	}

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid product ID")
		return
	}

	item := models.CartItem{
		ProductID: productID,
		Quantity:  req.Quantity,
	}

	if err := services.AddToCart(userID, item); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "item added to cart", nil)
}

// RemoveFromCart removes an item from the user's cart
func RemoveFromCart(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid user ID")
		return
	}

	var req RemoveFromCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid product ID")
		return
	}

	if err := services.RemoveFromCart(userID, productID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "item removed from cart", nil)
}

// ViewCart returns the user's cart
func ViewCart(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid user ID")
		return
	}

	cart, err := services.GetCart(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch cart")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "cart retrieved successfully", gin.H{"cart": cart})
}
