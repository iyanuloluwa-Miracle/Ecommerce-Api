package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/services"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// GetProducts returns all products
func GetProducts(c *gin.Context) {
	products, err := services.ListProducts()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch products")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "products retrieved successfully", gin.H{"products": products})
}

// GetProduct returns a single product by ID
func GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid product ID")
		return
	}

	product, err := services.GetProductByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "product not found")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "product retrieved successfully", gin.H{"product": product})
}

// CreateProduct creates a new product (admin only)
func CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   time.Now(),
	}

	if err := services.CreateProduct(product); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create product")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "product created", gin.H{"product": product})
}

// UpdateProduct updates an existing product (admin only)
func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid product ID")
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := services.UpdateProduct(id, req.Name, req.Description, req.Price, req.Stock); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update product")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "product updated", nil)
}

// DeleteProduct deletes a product (admin only)
func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid product ID")
		return
	}

	if err := services.DeleteProduct(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to delete product")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "product deleted", nil)
}
