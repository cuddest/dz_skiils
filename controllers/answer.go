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

// SQL queries for Answer
const (
	createAnswerQuery = `
		INSERT INTO answers (answer, question_id) 
		VALUES ($1, $2) RETURNING id`

	getAnswerQuery = `
		SELECT id, answer, question_id 
		FROM answers WHERE id = $1`

	getAllAnswersQuery = `
		SELECT id, answer, question_id 
		FROM answers`

	getAnswersByQuestionQuery = `
		SELECT id, answer, question_id 
		FROM answers WHERE question_id = $1`

	updateAnswerQuery = `
		UPDATE answers 
		SET answer = $1, question_id = $2 
		WHERE id = $3`

	deleteAnswerQuery = `
		DELETE FROM answers WHERE id = $1`
)

func NewAnswerController(db *sql.DB) *AnswerController {
	return &AnswerController{db: db}
}

func (h *AnswerController) validateAnswer(answer *models.Answer) error {
	if answer.Answer == "" {
		return errors.New("answer text is required")
	}
	if answer.QuestionID <= 0 {
		return errors.New("valid question ID is required")
	}
	return nil
}

// AnswerController handles operations on answers
// @title Answer API
// @description CRUD operations for managing answers
type AnswerController struct {
	db *sql.DB
}

// CreateAnswer godoc
// @Summary Create a new answer
// @Description Create a new answer for a specific question
// @Tags answers
// @Accept json
// @Produce json
// @Param answer body models.Answer true "Answer object to be created"
// @Success 201 {object} models.Answer
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /answers/CreateAnswer [post]

// CreateAnswer handles the creation of a new answer
func (h *AnswerController) CreateAnswer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var answer models.Answer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateAnswer(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify question exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM questions WHERE id = $1)", answer.QuestionID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify question"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	// Create answer
	err = h.db.QueryRowContext(ctx, createAnswerQuery, answer.Answer, answer.QuestionID).Scan(&answer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create answer"})
		return
	}

	c.JSON(http.StatusCreated, answer)
}

// GetAnswer godoc
// @Summary Get a specific answer
// @Description Get an answer by its ID
// @Tags answers
// @Accept json
// @Produce json
// @Param id path int true "Answer ID"
// @Success 200 {object} models.Answer
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /answers/GetAnswer [post]
// GetAnswer retrieves a specific answer by ID
func (h *AnswerController) GetAnswer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var answer models.Answer
	err = h.db.QueryRowContext(ctx, getAnswerQuery, id).Scan(&answer.ID, &answer.Answer, &answer.QuestionID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve answer"})
		return
	}

	c.JSON(http.StatusOK, answer)
}

// GetAllAnswers godoc
// @Summary Get all answers
// @Description Retrieve all answers from the database
// @Tags answers
// @Accept json
// @Produce json
// @Success 200 {array} models.Answer
// @Failure 500 {object} map[string]interface{}
// @Router /answers/GetAllAnswer [get]
// GetAllAnswers retrieves all answers
func (h *AnswerController) GetAllAnswers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllAnswersQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve answers"})
		return
	}
	defer rows.Close()

	var answers []models.Answer
	for rows.Next() {
		var answer models.Answer
		if err := rows.Scan(&answer.ID, &answer.Answer, &answer.QuestionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process answers"})
			return
		}
		answers = append(answers, answer)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing answers"})
		return
	}

	c.JSON(http.StatusOK, answers)
}

// GetAnswersByQuestion godoc
// @Summary Get answers by question
// @Description Get all answers for a specific question
// @Tags answers
// @Accept json
// @Produce json
// @Param questionId path int true "Question ID"
// @Success 200 {array} models.Answer
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /questions/GetAnswersByQuestion [post]
// GetAnswersByQuestion retrieves all answers for a specific question
func (h *AnswerController) GetAnswersByQuestion(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	questionID, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID format"})
		return
	}

	// Verify question exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM questions WHERE id = $1)", questionID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify question"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	rows, err := h.db.QueryContext(ctx, getAnswersByQuestionQuery, questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve answers"})
		return
	}
	defer rows.Close()

	var answers []models.Answer
	for rows.Next() {
		var answer models.Answer
		if err := rows.Scan(&answer.ID, &answer.Answer, &answer.QuestionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process answers"})
			return
		}
		answers = append(answers, answer)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing answers"})
		return
	}

	c.JSON(http.StatusOK, answers)
}

// @Summary Update an answer
// @Description Update an existing answer by its ID
// @Tags answers
// @Accept json
// @Produce json
// @Param id path int true "Answer ID"
// @Param answer body models.Answer true "Updated answer object"
// @Success 200 {object} models.Answer
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /answers/UpdateAnswer [put]
// UpdateAnswer updates an existing answer
func (h *AnswerController) UpdateAnswer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var answer models.Answer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateAnswer(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify question exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM questions WHERE id = $1)", answer.QuestionID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify question"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	_, err = h.db.ExecContext(ctx, updateAnswerQuery, answer.Answer, answer.QuestionID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update answer"})
		return
	}

	answer.ID = uint(id)
	c.JSON(http.StatusOK, answer)
}

// DeleteAnswer godoc
// @Summary Delete an answer
// @Description Delete an answer by its ID
// @Tags answers
// @Accept json
// @Produce json
// @Param id path int true "Answer ID"
// @Success 200 {object} models.Answer
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /answers/DeleteAnswer [delete]
// DeleteAnswer deletes an answer by ID
func (h *AnswerController) DeleteAnswer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	_, err = h.db.ExecContext(ctx, deleteAnswerQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete answer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Answer deleted successfully"})
}
