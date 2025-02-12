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

// SQL queries for ExamQuizz
const (
	createExamQuizzQuery = `
		INSERT INTO exam_quizzes (question, option1, option2, option3, option4, answer, exam_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	getExamQuizzQuery = `
		SELECT eq.id, eq.question, eq.option1, eq.option2, eq.option3, eq.option4, eq.answer, eq.exam_id,
			json_build_object(
				'ID', e.id,
				'Description', e.description
			) as exam
		FROM exam_quizzes eq
		LEFT JOIN exams e ON eq.exam_id = e.id 
		WHERE eq.id = $1`

	getAllExamQuizzesQuery = `
		SELECT eq.id, eq.question, eq.option1, eq.option2, eq.option3, eq.option4, eq.answer, eq.exam_id,
			json_build_object(
				'ID', e.id,
				'Description', e.description
			) as exam
		FROM exam_quizzes eq
		LEFT JOIN exams e ON eq.exam_id = e.id`

	getExamQuizzesByExamQuery = `
		SELECT eq.id, eq.question, eq.option1, eq.option2, eq.option3, eq.option4, eq.answer, eq.exam_id,
			json_build_object(
				'ID', e.id,
				'Description', e.description
			) as exam
		FROM exam_quizzes eq
		LEFT JOIN exams e ON eq.exam_id = e.id
		WHERE eq.exam_id = $1`

	updateExamQuizzQuery = `
		UPDATE exam_quizzes 
		SET question = $1, option1 = $2, option2 = $3, option3 = $4, option4 = $5, answer = $6, exam_id = $7 
		WHERE id = $8`

	deleteExamQuizzQuery = `
		DELETE FROM exam_quizzes WHERE id = $1`
)

type ExamQuizzController struct {
	db *sql.DB
}

func NewExamQuizzController(db *sql.DB) *ExamQuizzController {
	return &ExamQuizzController{db: db}
}

func (h *ExamQuizzController) validateExamQuizz(quizz *models.ExamQuizz) error {
	if quizz.Question == "" {
		return errors.New("question is required")
	}
	if quizz.Option1 == "" || quizz.Option2 == "" {
		return errors.New("at least two options are required")
	}
	if quizz.Answer == 0 {
		return errors.New("answer is required")
	}
	if quizz.ExamID <= 0 {
		return errors.New("valid exam ID is required")
	}
	return nil
}

// @Summary Create new exam quiz
// @Description Create a new exam quiz in the system
// @Tags examquizzes
// @Accept json
// @Produce json
// @Param quiz body models.ExamQuizz true "Exam Quiz object"
// @Success 201 {object} models.ExamQuizz
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /examquizzes/createExamQuiz [post]
// CreateExamQuizz creates a new exam quiz
func (h *ExamQuizzController) CreateExamQuizz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var quizz models.ExamQuizz
	if err := c.ShouldBindJSON(&quizz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateExamQuizz(&quizz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify exam exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exams WHERE id = $1)", quizz.ExamID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify exam"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	err = h.db.QueryRowContext(ctx, createExamQuizzQuery,
		quizz.Question, quizz.Option1, quizz.Option2, quizz.Option3,
		quizz.Option4, quizz.Answer, quizz.ExamID).Scan(&quizz.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exam quiz"})
		return
	}

	c.JSON(http.StatusCreated, quizz)
}

// @Summary Get exam quiz by ID
// @Description Get a specific exam quiz by its ID
// @Tags examquizzes
// @Accept json
// @Produce json
// @Param id body int true "Exam Quiz ID"
// @Success 200 {object} models.ExamQuizz
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /examquizzes/get [post]
// GetExamQuizz retrieves a specific exam quiz
func (h *ExamQuizzController) GetExamQuizz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var quizz models.ExamQuizz
	var examJSON []byte
	err = h.db.QueryRowContext(ctx, getExamQuizzQuery, id).Scan(
		&quizz.ID, &quizz.Question, &quizz.Option1, &quizz.Option2,
		&quizz.Option3, &quizz.Option4, &quizz.Answer, &quizz.ExamID, &examJSON,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam quiz not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exam quiz"})
		return
	}

	c.JSON(http.StatusOK, quizz)
}

// @Summary Get all exam quizzes
// @Description Retrieve all exam quizzes from the database
// @Tags examquizzes
// @Accept json
// @Produce json
// @Success 200 {array} models.ExamQuizz
// @Failure 500 {object} map[string]interface{}
// @Router /examquizzes/all [get]
// GetAllExamQuizzes retrieves all exam quizzes

func (h *ExamQuizzController) GetAllExamQuizzes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllExamQuizzesQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exam quizzes"})
		return
	}
	defer rows.Close()

	var quizzes []models.ExamQuizz
	for rows.Next() {
		var quizz models.ExamQuizz
		var examJSON []byte
		if err := rows.Scan(
			&quizz.ID, &quizz.Question, &quizz.Option1, &quizz.Option2,
			&quizz.Option3, &quizz.Option4, &quizz.Answer, &quizz.ExamID, &examJSON,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process exam quizzes"})
			return
		}
		quizzes = append(quizzes, quizz)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing exam quizzes"})
		return
	}

	c.JSON(http.StatusOK, quizzes)
}

// @Summary Get exam quizzes by exam
// @Description Get all quizzes for a specific exam
// @Tags examquizzes
// @Accept json
// @Produce json
// @Param examId body int true "Exam ID"
// @Success 200 {array} models.ExamQuizz
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /examquizzes/GetExamQuizzesByExam [post]
// GetExamQuizzesByExam retrieves quizzes by exam ID
func (h *ExamQuizzController) GetExamQuizzesByExam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	examID, err := strconv.Atoi(c.Param("examId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID format"})
		return
	}

	// Verify exam exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exams WHERE id = $1)", examID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify exam"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	rows, err := h.db.QueryContext(ctx, getExamQuizzesByExamQuery, examID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exam quizzes"})
		return
	}
	defer rows.Close()

	var quizzes []models.ExamQuizz
	for rows.Next() {
		var quizz models.ExamQuizz
		var examJSON []byte
		if err := rows.Scan(
			&quizz.ID, &quizz.Question, &quizz.Option1, &quizz.Option2,
			&quizz.Option3, &quizz.Option4, &quizz.Answer, &quizz.ExamID, &examJSON,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process exam quizzes"})
			return
		}
		quizzes = append(quizzes, quizz)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing exam quizzes"})
		return
	}

	c.JSON(http.StatusOK, quizzes)
}

// @Summary Update exam quiz
// @Description Update an existing exam quiz
// @Tags examquizzes
// @Accept json
// @Produce json
// @Param quiz body models.ExamQuizz true "Updated exam quiz object"
// @Success 200 {object} models.ExamQuizz
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /examquizzes/updateExamQuiz [put]
// UpdateExamQuizz updates an exam quiz
func (h *ExamQuizzController) UpdateExamQuizz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var quizz models.ExamQuizz
	if err := c.ShouldBindJSON(&quizz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateExamQuizz(&quizz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify exam exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exams WHERE id = $1)", quizz.ExamID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify exam"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	result, err := h.db.ExecContext(ctx, updateExamQuizzQuery,
		quizz.Question, quizz.Option1, quizz.Option2, quizz.Option3,
		quizz.Option4, quizz.Answer, quizz.ExamID, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exam quiz"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam quiz not found"})
		return
	}

	quizz.ID = uint(id)
	c.JSON(http.StatusOK, quizz)
}

// @Summary Delete exam quiz
// @Description Delete an exam quiz by its ID
// @Tags examquizzes
// @Accept json
// @Produce json
// @Param id body int true "Exam Quiz ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /examquizzes/DeleteExamQuiz [delete]
// DeleteExamQuizz deletes an exam quiz

func (h *ExamQuizzController) DeleteExamQuizz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteExamQuizzQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exam quiz"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam quiz not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exam quiz deleted successfully"})
}
