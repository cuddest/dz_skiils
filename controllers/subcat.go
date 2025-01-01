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

// SQL queries for SubCat
const (
	createSubCatQuery = `
		INSERT INTO sub_cats (name, category_id) 
		VALUES ($1, $2) RETURNING id`

	getSubCatQuery = `
		SELECT id, name, category_id 
		FROM sub_cats WHERE id = $1`

	getAllSubCatsQuery = `
		SELECT id, name, category_id 
		FROM sub_cats`

	getSubCatsByCategoryQuery = `
		SELECT id, name, category_id 
		FROM sub_cats WHERE category_id = $1`

	updateSubCatQuery = `
		UPDATE sub_cats 
		SET name = $1, category_id = $2 
		WHERE id = $3`

	deleteSubCatQuery = `
		DELETE FROM sub_cats WHERE id = $1`
)

type SubCatController struct {
	db *sql.DB
}

func NewSubCatController(db *sql.DB) *SubCatController {
	return &SubCatController{db: db}
}

func (h *SubCatController) validateSubCat(subcat *models.SubCat) error {
	if subcat.Name == "" {
		return errors.New("name is required")
	}
	if subcat.CategoryID <= 0 {
		return errors.New("valid category ID is required")
	}
	return nil
}

// CreateSubCat handles the creation of a new subcategory
func (h *SubCatController) CreateSubCat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var subcat models.SubCat
	if err := c.ShouldBindJSON(&subcat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateSubCat(&subcat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify category exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", subcat.CategoryID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify category"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Create subcategory
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Insert subcategory
	err = tx.QueryRowContext(ctx, createSubCatQuery,
		subcat.Name, subcat.CategoryID).Scan(&subcat.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subcategory"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, subcat)
}

// GetSubCat retrieves a specific subcategory by ID
func (h *SubCatController) GetSubCat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var subcat models.SubCat
	err = h.db.QueryRowContext(ctx, getSubCatQuery, id).Scan(
		&subcat.ID, &subcat.Name, &subcat.CategoryID,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subcategory not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subcategory"})
		return
	}

	c.JSON(http.StatusOK, subcat)
}

// GetAllSubCats retrieves all subcategories
func (h *SubCatController) GetAllSubCats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllSubCatsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subcategories"})
		return
	}
	defer rows.Close()

	var subcats []models.SubCat
	for rows.Next() {
		var subcat models.SubCat
		if err := rows.Scan(
			&subcat.ID, &subcat.Name, &subcat.CategoryID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process subcategories"})
			return
		}
		subcats = append(subcats, subcat)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing subcategories"})
		return
	}

	c.JSON(http.StatusOK, subcats)
}

// GetSubCatsByCategory retrieves all subcategories for a specific category
func (h *SubCatController) GetSubCatsByCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	categoryID, err := strconv.Atoi(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format"})
		return
	}

	// Verify category exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", categoryID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify category"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	rows, err := h.db.QueryContext(ctx, getSubCatsByCategoryQuery, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subcategories"})
		return
	}
	defer rows.Close()

	var subcats []models.SubCat
	for rows.Next() {
		var subcat models.SubCat
		if err := rows.Scan(
			&subcat.ID, &subcat.Name, &subcat.CategoryID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process subcategories"})
			return
		}
		subcats = append(subcats, subcat)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing subcategories"})
		return
	}

	c.JSON(http.StatusOK, subcats)
}

// UpdateSubCat updates an existing subcategory
func (h *SubCatController) UpdateSubCat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var subcat models.SubCat
	if err := c.ShouldBindJSON(&subcat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateSubCat(&subcat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify category exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", subcat.CategoryID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify category"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	result, err := h.db.ExecContext(ctx, updateSubCatQuery,
		subcat.Name, subcat.CategoryID, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subcategory"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subcategory not found"})
		return
	}

	subcat.ID = uint(id)
	c.JSON(http.StatusOK, subcat)
}

// DeleteSubCat deletes a subcategory by ID
func (h *SubCatController) DeleteSubCat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteSubCatQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subcategory"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subcategory not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subcategory deleted successfully"})
}