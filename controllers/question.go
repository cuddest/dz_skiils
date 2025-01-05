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

// SQL queries for Question
const (
	createQuestionQuery = `
		INSERT INTO questions (course_id, student_id, question) 
		VALUES ($1, $2, $3) RETURNING id`

	getQuestionQuery = `
		SELECT id, course_id, student_id, question 
		FROM questions WHERE id = $1`

	getAllQuestionsQuery = `
		SELECT id, course_id, student_id, question 
		FROM questions`

	updateQuestionQuery = `
		UPDATE questions 
		SET course_id = $1, student_id = $2, question = $3 
		WHERE id = $4`

	deleteQuestionQuery = `
		DELETE FROM questions WHERE id = $1`
)

type QuestionController struct {
	db *sql.DB
}

func NewQuestionController(db *sql.DB) *QuestionController {
	return &QuestionController{db: db}
}

func (h *QuestionController) validateQuestion(question *models.Question) error {
	if question.CourseID <= 0 {
		return errors.New("valid course ID is required")
	}
	if question.StudentID <= 0 {
		return errors.New("valid student ID is required")
	}
	if question.Question == "" {
		return errors.New("question text is required")
	}
	return nil
}

// @Summary Create a new question
// @Description Create a new question in the system
// @Tags questions
// @Accept json
// @Produce json
// @Param question body models.Question true "Question information"
// @Success 201 {object} models.Question
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /questions/createQuestion [post]
func (h *QuestionController) CreateQuestion(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateQuestion(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, createQuestionQuery,
		question.CourseID, question.StudentID, question.Question).Scan(&question.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, question)
}

// @Summary Get question by ID
// @Description Retrieve a question by its ID
// @Tags questions
// @Accept json
// @Produce json
// @Param id path int true "Question ID"
// @Success 200 {object} models.Question
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /questions/get [post]
func (h *QuestionController) GetQuestion(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var question models.Question
	err = h.db.QueryRowContext(ctx, getQuestionQuery, id).Scan(
		&question.ID, &question.CourseID, &question.StudentID, &question.Question,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve question"})
		return
	}

	c.JSON(http.StatusOK, question)
}

// @Summary Get all questions
// @Description Retrieve a list of all questions
// @Tags questions
// @Accept json
// @Produce json
// @Success 200 {array} models.Question
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /questions/all [get]
func (h *QuestionController) GetAllQuestions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllQuestionsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
		return
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var question models.Question
		if err := rows.Scan(
			&question.ID, &question.CourseID, &question.StudentID, &question.Question,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process questions"})
			return
		}
		questions = append(questions, question)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing questions"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// @Summary Update question
// @Description Update an existing question
// @Tags questions
// @Accept json
// @Produce json
// @Param id path int true "Question ID"
// @Param question body models.Question true "Updated question information"
// @Success 200 {object} models.Question
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /questions/updateQuestion [put]

func (h *QuestionController) UpdateQuestion(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateQuestion(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.db.ExecContext(ctx, updateQuestionQuery,
		question.CourseID, question.StudentID, question.Question, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	question.ID = uint(id)
	c.JSON(http.StatusOK, question)
}

// @Summary Delete question
// @Description Delete a question from the system
// @Tags questions
// @Accept json
// @Produce json
// @Param id path int true "Question ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /questions/DeleteQuestion [delete]

func (h *QuestionController) DeleteQuestion(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteQuestionQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}
