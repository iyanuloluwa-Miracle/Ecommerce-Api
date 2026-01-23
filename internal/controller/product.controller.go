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

// CreateProductRequest represents the create product request body
// @Description Create product request
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required" example:"Wireless Headphones"`
	Description string  `json:"description" binding:"required" example:"High-quality wireless headphones"`
	Price       float64 `json:"price" binding:"required,gt=0" example:"99.99"`
	Stock       int     `json:"stock" binding:"required,gte=0" example:"50"`
}

// UpdateProductRequest represents the update product request body
// @Description Update product request
type UpdateProductRequest struct {
	Name        string  `json:"name" example:"Updated Product Name"`
	Description string  `json:"description" example:"Updated description"`
	Price       float64 `json:"price" example:"149.99"`
	Stock       int     `json:"stock" example:"100"`
}

// GetProducts returns all products
// @Summary Get all products
// @Description Retrieve all products from the catalog
// @Tags Products
// @Produce json
// @Success 200 {object} utils.Response{data=object{products=[]models.Product}} "Products retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products [get]
func GetProducts(c *gin.Context) {
	products, err := services.ListProducts()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch products")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "products retrieved successfully", gin.H{"products": products})
}

// GetProduct returns a single product by ID
// @Summary Get product by ID
// @Description Retrieve a single product by its ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response{data=object{product=models.Product}} "Product retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid product ID"
// @Failure 404 {object} utils.Response "Product not found"
// @Router /products/{id} [get]
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
// @Summary Create a new product
// @Description Create a new product in the catalog (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProductRequest true "Product details"
// @Success 201 {object} utils.Response{data=object{product=models.Product}} "Product created successfully"
// @Failure 400 {object} utils.Response "Bad request - validation error"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/products [post]
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
// @Summary Update a product
// @Description Update an existing product by ID (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body UpdateProductRequest true "Updated product details"
// @Success 200 {object} utils.Response "Product updated successfully"
// @Failure 400 {object} utils.Response "Bad request - invalid ID or validation error"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/products/{id} [put]
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
// @Summary Delete a product
// @Description Delete a product by ID (admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response "Product deleted successfully"
// @Failure 400 {object} utils.Response "Invalid product ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/products/{id} [delete]
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
