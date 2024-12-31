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
	createQuizzQuery = `
		INSERT INTO course_quizzes (question, option1, option2, option3, option4, answer, course_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	getQuizzQuery = `
		SELECT id, question, option1, option2, option3, option4, answer, course_id 
		FROM course_quizzes WHERE id = $1`

	getAllQuizzesQuery = `
		SELECT id, question, option1, option2, option3, option4, answer, course_id 
		FROM course_quizzes`

	getQuizzesByCourseQuery = `
		SELECT id, question, option1, option2, option3, option4, answer, course_id 
		FROM course_quizzes WHERE course_id = $1`

	updateQuizzQuery = `
		UPDATE course_quizzes 
		SET question = $1, option1 = $2, option2 = $3, option3 = $4, option4 = $5, answer = $6, course_id = $7 
		WHERE id = $8`

	deleteQuizzQuery = `
		DELETE FROM course_quizzes WHERE id = $1`
)

type CourseQuizzController struct {
	db *sql.DB
}

func NewCourseQuizzController(db *sql.DB) *CourseQuizzController {
	return &CourseQuizzController{db: db}
}

func (h *CourseQuizzController) validateQuizz(quizz *models.CourseQuizz) error {
	if quizz.Question == "" {
		return errors.New("question is required")
	}
	if quizz.Option1 == "" || quizz.Option2 == "" {
		return errors.New("at least two options are required")
	}
	if quizz.Answer == "" {
		return errors.New("answer is required")
	}
	if quizz.CourseID <= 0 {
		return errors.New("valid course ID is required")
	}
	return nil
}

// CreateQuizz handles the creation of a new quiz
func (h *CourseQuizzController) CreateQuizz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var quizz models.CourseQuizz
	if err := c.ShouldBindJSON(&quizz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateQuizz(&quizz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify course exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", quizz.CourseID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify course"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	err = h.db.QueryRowContext(ctx, createQuizzQuery,
		quizz.Question, quizz.Option1, quizz.Option2,
		quizz.Option3, quizz.Option4, quizz.Answer,
		quizz.CourseID).Scan(&quizz.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz"})
		return
	}

	c.JSON(http.StatusCreated, quizz)
}

// GetQuizz retrieves a specific quiz by ID
func (h *CourseQuizzController) GetQuizz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var quizz models.CourseQuizz
	err = h.db.QueryRowContext(ctx, getQuizzQuery, id).Scan(
		&quizz.ID, &quizz.Question, &quizz.Option1,
		&quizz.Option2, &quizz.Option3, &quizz.Option4,
		&quizz.Answer, &quizz.CourseID,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve quiz"})
		return
	}

	c.JSON(http.StatusOK, quizz)
}

// GetAllQuizzes retrieves all quizzes
func (h *CourseQuizzController) GetAllQuizzes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllQuizzesQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve quizzes"})
		return
	}
	defer rows.Close()

	var quizzes []models.CourseQuizz
	for rows.Next() {
		var quizz models.CourseQuizz
		if err := rows.Scan(
			&quizz.ID, &quizz.Question, &quizz.Option1,
			&quizz.Option2, &quizz.Option3, &quizz.Option4,
			&quizz.Answer, &quizz.CourseID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process quizzes"})
			return
		}
		quizzes = append(quizzes, quizz)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing quizzes"})
		return
	}

	c.JSON(http.StatusOK, quizzes)
}

// GetQuizzesByCourse retrieves all quizzes for a specific course
func (h *CourseQuizzController) GetQuizzesByCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	// Verify course exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", courseID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify course"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	rows, err := h.db.QueryContext(ctx, getQuizzesByCourseQuery, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve quizzes"})
		return
	}
	defer rows.Close()

	var quizzes []models.CourseQuizz
	for rows.Next() {
		var quizz models.CourseQuizz
		if err := rows.Scan(
			&quizz.ID, &quizz.Question, &quizz.Option1,
			&quizz.Option2, &quizz.Option3, &quizz.Option4,
			&quizz.Answer, &quizz.CourseID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process quizzes"})
			return
		}
		quizzes = append(quizzes, quizz)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing quizzes"})
		return
	}

	c.JSON(http.StatusOK, quizzes)
}

// UpdateQuizz updates an existing quiz
func (h *CourseQuizzController) UpdateQuizz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var quizz models.CourseQuizz
	if err := c.ShouldBindJSON(&quizz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateQuizz(&quizz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify course exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", quizz.CourseID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify course"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	result, err := h.db.ExecContext(ctx, updateQuizzQuery,
		quizz.Question, quizz.Option1, quizz.Option2,
		quizz.Option3, quizz.Option4, quizz.Answer,
		quizz.CourseID, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quiz"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	quizz.ID = uint(id)
	c.JSON(http.StatusOK, quizz)
}

// DeleteQuizz deletes a quiz by ID
func (h *CourseQuizzController) DeleteQuizz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteQuizzQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete quiz"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quiz deleted successfully"})
}
