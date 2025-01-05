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

// SQL queries for Exam
const (
	createExamQuery = `
		INSERT INTO exams (description, course_id) 
		VALUES ($1, $2) RETURNING id`

	getExamQuery = `
		SELECT e.id, e.description, e.course_id,
			json_build_object(
				'ID', c.id,
				'Name', c.name,
				'Description', c.description
			) as course
		FROM exams e
		LEFT JOIN courses c ON e.course_id = c.id 
		WHERE e.id = $1`

	getAllExamsQuery = `
		SELECT e.id, e.description, e.course_id,
			json_build_object(
				'ID', c.id,
				'Name', c.name,
				'Description', c.description
			) as course
		FROM exams e
		LEFT JOIN courses c ON e.course_id = c.id`

	getExamsByCourseQuery = `
		SELECT e.id, e.description, e.course_id,
			json_build_object(
				'ID', c.id,
				'Name', c.name,
				'Description', c.description
			) as course
		FROM exams e
		LEFT JOIN courses c ON e.course_id = c.id
		WHERE e.course_id = $1`

	updateExamQuery = `
		UPDATE exams 
		SET description = $1, course_id = $2 
		WHERE id = $3`

	deleteExamQuery = `
		DELETE FROM exams WHERE id = $1`
)

type ExamController struct {
	db *sql.DB
}

func NewExamController(db *sql.DB) *ExamController {
	return &ExamController{db: db}
}

func (h *ExamController) validateExam(exam *models.Exam) error {
	if exam.Description == "" {
		return errors.New("description is required")
	}
	if exam.CourseID <= 0 {
		return errors.New("valid course ID is required")
	}
	return nil
}

// @Summary Create new exam
// @Description Create a new exam in the system
// @Tags exams
// @Accept json
// @Produce json
// @Param exam body models.Exam true "Exam object"
// @Success 201 {object} models.Exam
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /exams/createExam [post]
// CreateExam creates a new exam
func (h *ExamController) CreateExam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var exam models.Exam
	if err := c.ShouldBindJSON(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateExam(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify course exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", exam.CourseID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify course"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	err = h.db.QueryRowContext(ctx, createExamQuery,
		exam.Description, exam.CourseID).Scan(&exam.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exam"})
		return
	}

	c.JSON(http.StatusCreated, exam)
}

// @Summary Get exam by ID
// @Description Get a specific exam by its ID
// @Tags exams
// @Accept json
// @Produce json
// @Param id body int true "Exam ID"
// @Success 200 {object} models.Exam
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /exams/get [post]
// GetExam retrieves a specific exam
func (h *ExamController) GetExam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var exam models.Exam
	var courseJSON []byte
	err = h.db.QueryRowContext(ctx, getExamQuery, id).Scan(
		&exam.ID, &exam.Description, &exam.CourseID, &courseJSON,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exam"})
		return
	}

	c.JSON(http.StatusOK, exam)
}

// @Summary Get all exams
// @Description Retrieve all exams from the database
// @Tags exams
// @Accept json
// @Produce json
// @Success 200 {array} models.Exam
// @Failure 500 {object} map[string]interface{}
// @Router /exams/all [get]
// GetAllExams retrieves all exams
func (h *ExamController) GetAllExams(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllExamsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exams"})
		return
	}
	defer rows.Close()

	var exams []models.Exam
	for rows.Next() {
		var exam models.Exam
		var courseJSON []byte
		if err := rows.Scan(
			&exam.ID, &exam.Description, &exam.CourseID, &courseJSON,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process exams"})
			return
		}
		exams = append(exams, exam)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing exams"})
		return
	}

	c.JSON(http.StatusOK, exams)
}

// @Summary Get exams by course
// @Description Get all exams for a specific course
// @Tags exams
// @Accept json
// @Produce json
// @Param courseId body int true "Course ID"
// @Success 200 {array} models.Exam
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /exams/GetExamsByCourse [post]
// GetExamsByCourse retrieves exams by course ID
func (h *ExamController) GetExamsByCourse(c *gin.Context) {
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

	rows, err := h.db.QueryContext(ctx, getExamsByCourseQuery, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exams"})
		return
	}
	defer rows.Close()

	var exams []models.Exam
	for rows.Next() {
		var exam models.Exam
		var courseJSON []byte
		if err := rows.Scan(
			&exam.ID, &exam.Description, &exam.CourseID, &courseJSON,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process exams"})
			return
		}
		exams = append(exams, exam)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing exams"})
		return
	}

	c.JSON(http.StatusOK, exams)
}

// @Summary Update exam
// @Description Update an existing exam
// @Tags exams
// @Accept json
// @Produce json
// @Param exam body models.Exam true "Updated exam object"
// @Success 200 {object} models.Exam
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /exams/updateExam [put]
// UpdateExam updates an exam
func (h *ExamController) UpdateExam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var exam models.Exam
	if err := c.ShouldBindJSON(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateExam(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify course exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", exam.CourseID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify course"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	result, err := h.db.ExecContext(ctx, updateExamQuery,
		exam.Description, exam.CourseID, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exam"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	exam.ID = uint(id)
	c.JSON(http.StatusOK, exam)
}

// @Summary Delete exam
// @Description Delete an exam by its ID
// @Tags exams
// @Accept json
// @Produce json
// @Param id body int true "Exam ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /exams/DeleteExam [delete]
// DeleteExam deletes an exam
func (h *ExamController) DeleteExam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteExamQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exam"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exam deleted successfully"})
}
