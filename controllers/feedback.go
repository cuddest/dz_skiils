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

// SQL queries for Feedback
const (
	createFeedbackQuery = `
		INSERT INTO feedbacks (description, review, student_id) 
		VALUES ($1, $2, $3) RETURNING id`

	getFeedbackQuery = `
		SELECT f.id, f.description, f.review, f.student_id,
			json_build_object(
				'ID', s.id,
				'Name', s.name,
				'Email', s.email
			) as student
		FROM feedbacks f
		LEFT JOIN students s ON f.student_id = s.id 
		WHERE f.id = $1`

	getAllFeedbacksQuery = `
		SELECT f.id, f.description, f.review, f.student_id,
			json_build_object(
				'ID', s.id,
				'Name', s.name,
				'Email', s.email
			) as student
		FROM feedbacks f
		LEFT JOIN students s ON f.student_id = s.id`

	getFeedbacksByStudentQuery = `
		SELECT f.id, f.description, f.review, f.student_id,
			json_build_object(
				'ID', s.id,
				'Name', s.name,
				'Email', s.email
			) as student
		FROM feedbacks f
		LEFT JOIN students s ON f.student_id = s.id
		WHERE f.student_id = $1`

	updateFeedbackQuery = `
		UPDATE feedbacks 
		SET description = $1, review = $2, student_id = $3 
		WHERE id = $4`

	deleteFeedbackQuery = `
		DELETE FROM feedbacks WHERE id = $1`
)

type FeedbackController struct {
	db *sql.DB
}

func NewFeedbackController(db *sql.DB) *FeedbackController {
	return &FeedbackController{db: db}
}

func (h *FeedbackController) validateFeedback(feedback *models.Feedback) error {
	if feedback.Description == "" {
		return errors.New("description is required")
	}
	if feedback.Review < 1 || feedback.Review > 5 {
		return errors.New("review must be between 1 and 5")
	}
	if feedback.StudentID <= 0 {
		return errors.New("valid student ID is required")
	}
	return nil
}

// CreateFeedback handles the creation of a new feedback
func (h *FeedbackController) CreateFeedback(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var feedback models.Feedback
	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateFeedback(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify student exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM students WHERE id = $1)", feedback.StudentID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify student"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	err = h.db.QueryRowContext(ctx, createFeedbackQuery,
		feedback.Description, feedback.Review,
		feedback.StudentID).Scan(&feedback.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create feedback"})
		return
	}

	c.JSON(http.StatusCreated, feedback)
}

// GetFeedback retrieves a specific feedback by ID
func (h *FeedbackController) GetFeedback(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var feedback models.Feedback
	var studentJSON []byte
	err = h.db.QueryRowContext(ctx, getFeedbackQuery, id).Scan(
		&feedback.ID, &feedback.Description, &feedback.Review,
		&feedback.StudentID, &studentJSON,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback"})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// GetAllFeedbacks retrieves all feedbacks
func (h *FeedbackController) GetAllFeedbacks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllFeedbacksQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedbacks"})
		return
	}
	defer rows.Close()

	var feedbacks []models.Feedback
	for rows.Next() {
		var feedback models.Feedback
		var studentJSON []byte
		if err := rows.Scan(
			&feedback.ID, &feedback.Description, &feedback.Review,
			&feedback.StudentID, &studentJSON,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process feedbacks"})
			return
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing feedbacks"})
		return
	}

	c.JSON(http.StatusOK, feedbacks)
}

// GetFeedbacksByStudent retrieves all feedbacks for a specific student
func (h *FeedbackController) GetFeedbacksByStudent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	studentID, err := strconv.Atoi(c.Param("studentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID format"})
		return
	}

	// Verify student exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM students WHERE id = $1)", studentID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify student"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	rows, err := h.db.QueryContext(ctx, getFeedbacksByStudentQuery, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedbacks"})
		return
	}
	defer rows.Close()

	var feedbacks []models.Feedback
	for rows.Next() {
		var feedback models.Feedback
		var studentJSON []byte
		if err := rows.Scan(
			&feedback.ID, &feedback.Description, &feedback.Review,
			&feedback.StudentID, &studentJSON,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process feedbacks"})
			return
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing feedbacks"})
		return
	}

	c.JSON(http.StatusOK, feedbacks)
}

// UpdateFeedback updates an existing feedback
func (h *FeedbackController) UpdateFeedback(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var feedback models.Feedback
	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateFeedback(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify student exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM students WHERE id = $1)", feedback.StudentID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify student"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	result, err := h.db.ExecContext(ctx, updateFeedbackQuery,
		feedback.Description, feedback.Review,
		feedback.StudentID, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update feedback"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
		return
	}

	feedback.ID = uint(id)
	c.JSON(http.StatusOK, feedback)
}

// DeleteFeedback deletes a feedback by ID
func (h *FeedbackController) DeleteFeedback(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteFeedbackQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete feedback"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feedback deleted successfully"})
}
