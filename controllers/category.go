package controllers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/cuddest/dz-skills/models"
	"github.com/gin-gonic/gin"
)

// SQL queries for Category
const (
	createCategoryQuery = `
		INSERT INTO categories (name) 
		VALUES ($1) RETURNING id`

	getCategoryQuery = `
		SELECT c.id, c.name,
			COALESCE(json_agg(
				json_build_object(
					'ID', s.id,
					'Name', s.name,
					'category_id', s.category_id
				)
			) FILTER (WHERE s.id IS NOT NULL), '[]') as subcats
		FROM categories c
		LEFT JOIN sub_cats s ON c.id = s.category_id
		WHERE c.id = $1
		GROUP BY c.id, c.name`

	getAllCategoriesQuery = `
		SELECT c.id, c.name,
			COALESCE(json_agg(
				json_build_object(
					'ID', s.id,
					'Name', s.name,
					'category_id', s.category_id
				)
			) FILTER (WHERE s.id IS NOT NULL), '[]') as subcats
		FROM categories c
		LEFT JOIN sub_cats s ON c.id = s.category_id
		GROUP BY c.id, c.name`

	updateCategoryQuery = `
		UPDATE categories 
		SET name = $1 
		WHERE id = $2`

	deleteCategoryQuery = `
		DELETE FROM categories WHERE id = $1`
)

// CategoryController handles operations on categories
// @title Category API
// @description CRUD operations for managing categories
type CategoryController struct {
	db *sql.DB
}

func NewCategoryController(db *sql.DB) *CategoryController {
	return &CategoryController{db: db}
}

// Category Controller Methods

func (h *CategoryController) validateCategory(category *models.Category) error {
	if category.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category object to be created"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/createCategory [post]
func (h *CategoryController) CreateCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateCategory(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.db.QueryRowContext(ctx, createCategoryQuery, category.Name).Scan(&category.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategory godoc
// @Summary Get a specific category
// @Description Get a category by its ID, including its subcategories
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/get [post]
func (h *CategoryController) GetCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var category models.Category
	var subcatsJSON []byte
	err = h.db.QueryRowContext(ctx, getCategoryQuery, id).Scan(
		&category.ID, &category.Name, &subcatsJSON,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Retrieve all categories with their subcategories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]interface{}
// @Router /categories/all [get]
func (h *CategoryController) GetAllCategories(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllCategoriesQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		var subcatsJSON []byte
		if err := rows.Scan(&category.ID, &category.Name, &subcatsJSON); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process categories"})
			return
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.Category true "Updated category object"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/updateCategory [put]
func (h *CategoryController) UpdateCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateCategory(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.db.ExecContext(ctx, updateCategoryQuery, category.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	category.ID = uint(id)
	c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/DeleteCategory [delete]
func (h *CategoryController) DeleteCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteCategoryQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
