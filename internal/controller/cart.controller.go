package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/services"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddToCartRequest represents the add to cart request body
// @Description Add item to cart request
type AddToCartRequest struct {
	ProductID string `json:"product_id" binding:"required" example:"507f1f77bcf86cd799439011"`
	Quantity  int    `json:"quantity" binding:"required,gt=0" example:"2"`
}

// RemoveFromCartRequest represents the remove from cart request body
// @Description Remove item from cart request
type RemoveFromCartRequest struct {
	ProductID string `json:"product_id" binding:"required" example:"507f1f77bcf86cd799439011"`
}

// AddToCart adds an item to the user's cart
// @Summary Add item to cart
// @Description Add a product to the authenticated user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddToCartRequest true "Product and quantity to add"
// @Success 200 {object} utils.Response "Item added to cart"
// @Failure 400 {object} utils.Response "Bad request - validation error or invalid product ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /cart/add [post]
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
// @Summary Remove item from cart
// @Description Remove a product from the authenticated user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RemoveFromCartRequest true "Product to remove"
// @Success 200 {object} utils.Response "Item removed from cart"
// @Failure 400 {object} utils.Response "Bad request - validation error or invalid product ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /cart/remove [delete]
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
// @Summary View cart
// @Description Get the authenticated user's shopping cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=object{cart=models.Cart}} "Cart retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /cart [get]
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
