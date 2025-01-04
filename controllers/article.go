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

// SQL queries as constants
const (
	// Article queries
	createArticleQuery = `
		INSERT INTO articles (title, link, description, course_id) 
		VALUES ($1, $2, $3, $4) RETURNING id`
	getArticleQuery = `
		SELECT id, title, link, description, course_id 
		FROM articles WHERE id = $1`
	getAllArticlesQuery = `
		SELECT id, title, link, description, course_id 
		FROM articles`
	getArticlesByCourseQuery = `
		SELECT id, title, link, description, course_id 
		FROM articles WHERE course_id = $1`
	updateArticleQuery = `
		UPDATE articles 
		SET title = $1, link = $2, description = $3, course_id = $4 
		WHERE id = $5`
	deleteArticleQuery = `
		DELETE FROM articles WHERE id = $1`
)

// ArticleController handles operations on articles
// @title Article API
// @description CRUD operations for managing articles
type ArticleController struct {
	db *sql.DB
}

func NewArticleController(db *sql.DB) *ArticleController {
	return &ArticleController{db: db}
}

func (h *ArticleController) validateArticle(article *models.Article) error {
	if article.Title == "" {
		return errors.New("article title is required")
	}
	if article.Link == "" {
		return errors.New("article link is required")
	}
	if article.CourseID <= 0 {
		return errors.New("valid course ID is required")
	}
	return nil
}

// CreateArticle godoc
// @Summary Create a new article
// @Description Create a new article for a specific course
// @Tags articles
// @Accept json
// @Produce json
// @Param article body models.Article true "Article object to be created"
// @Success 201 {object} models.Article
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /articles/createArticle [post]
func (h *ArticleController) CreateArticle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateArticle(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify course exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", article.CourseID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify course"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	err = h.db.QueryRowContext(ctx, createArticleQuery,
		article.Title, article.Link, article.Description, article.CourseID).Scan(&article.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	c.JSON(http.StatusCreated, article)
}

// GetArticle godoc
// @Summary Get a specific article
// @Description Get an article by its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} models.Article
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /articles/get [post]
func (h *ArticleController) GetArticle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var article models.Article
	err = h.db.QueryRowContext(ctx, getArticleQuery, id).Scan(
		&article.ID, &article.Title, &article.Link,
		&article.Description, &article.CourseID,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve article"})
		return
	}

	c.JSON(http.StatusOK, article)
}

// GetAllArticles godoc
// @Summary Get all articles
// @Description Retrieve all articles from the database
// @Tags articles
// @Accept json
// @Produce json
// @Success 200 {array} models.Article
// @Failure 500 {object} map[string]interface{}
// @Router /articles/all [get]
func (h *ArticleController) GetAllArticles(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllArticlesQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles"})
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(
			&article.ID, &article.Title, &article.Link,
			&article.Description, &article.CourseID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process articles"})
			return
		}
		articles = append(articles, article)
	}

	c.JSON(http.StatusOK, articles)
}

// GetArticlesByCourse godoc
// @Summary Get articles by course
// @Description Get all articles for a specific course
// @Tags articles
// @Accept json
// @Produce json
// @Param courseId path int true "Course ID"
// @Success 200 {array} models.Article
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /articles/GetArticlesByCourse [post]
func (h *ArticleController) GetArticlesByCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	rows, err := h.db.QueryContext(ctx, getArticlesByCourseQuery, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles"})
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(
			&article.ID, &article.Title, &article.Link,
			&article.Description, &article.CourseID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process articles"})
			return
		}
		articles = append(articles, article)
	}

	c.JSON(http.StatusOK, articles)
}

// UpdateArticle godoc
// @Summary Update an article
// @Description Update an existing article by its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param article body models.Article true "Updated article object"
// @Success 200 {object} models.Article
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /articles/updateArticle[put]
func (h *ArticleController) UpdateArticle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateArticle(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.db.ExecContext(ctx, updateArticleQuery,
		article.Title, article.Link, article.Description,
		article.CourseID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	article.ID = uint(id)
	c.JSON(http.StatusOK, article)
}

// DeleteArticle godoc
// @Summary Delete an article
// @Description Delete an article by its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /articles/DeleteArticle [delete]
func (h *ArticleController) DeleteArticle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteArticleQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}
