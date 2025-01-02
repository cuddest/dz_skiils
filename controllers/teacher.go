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

// SQL queries for Teacher operations
const (
	createTeacherQuery = `
		INSERT INTO teachers (full_name, username, email, password, picture, skills, degrees, experience) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id`

	getTeacherQuery = `
		SELECT id, full_name, username, email, password, picture, skills, degrees, experience 
		FROM teachers 
		WHERE id = $1`

	getAllTeachersQuery = `
		SELECT id, full_name, username, email, password, picture, skills, degrees, experience 
		FROM teachers`

	updateTeacherQuery = `
		UPDATE teachers 
		SET full_name = $1, username = $2, email = $3, password = $4, 
		    picture = $5, skills = $6, degrees = $7, experience = $8 
		WHERE id = $9`

	deleteTeacherQuery = `
		DELETE FROM teachers WHERE id = $1`

	checkUsernameQuery = `
		SELECT EXISTS(SELECT 1 FROM teachers WHERE username = $1 AND id != $2)`

	checkEmailQuery = `
		SELECT EXISTS(SELECT 1 FROM teachers WHERE email = $1 AND id != $2)`
)

// TeacherController handles HTTP requests for Teacher operations
type TeacherController struct {
	db *sql.DB
}

// NewTeacherController creates a new TeacherController instance
func NewTeacherController(db *sql.DB) *TeacherController {
	return &TeacherController{db: db}
}

// validateTeacher performs validation on teacher data
func (h *TeacherController) validateTeacher(teacher *models.Teacher) error {
	if teacher.FullName == "" {
		return errors.New("full name is required")
	}
	if teacher.Username == "" {
		return errors.New("username is required")
	}
	if teacher.Email == "" {
		return errors.New("email is required")
	}
	if teacher.Password == "" && teacher.ID == 0 { // Only require password for new teachers
		return errors.New("password is required")
	}
	// Add more validation as needed
	return nil
}

// checkUniqueness verifies username and email uniqueness
func (h *TeacherController) checkUniqueness(ctx context.Context, teacher *models.Teacher) error {
	var exists bool

	// Check username uniqueness
	if err := h.db.QueryRowContext(ctx, checkUsernameQuery,
		teacher.Username, teacher.ID).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return errors.New("username already exists")
	}

	// Check email uniqueness
	if err := h.db.QueryRowContext(ctx, checkEmailQuery,
		teacher.Email, teacher.ID).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	return nil
}

// CreateTeacher handles the creation of a new teacher
func (h *TeacherController) CreateTeacher(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var teacher models.Teacher
	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateTeacher(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.checkUniqueness(ctx, &teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Insert teacher
	err = tx.QueryRowContext(ctx, createTeacherQuery,
		teacher.FullName, teacher.Username, teacher.Email,
		teacher.Password, teacher.Picture, teacher.Skills,
		teacher.Degrees, teacher.Experience).Scan(&teacher.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create teacher"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Clear sensitive data before sending response
	teacher.Password = ""
	c.JSON(http.StatusCreated, teacher)
}

// GetTeacher retrieves a specific teacher by ID
func (h *TeacherController) GetTeacher(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var teacher models.Teacher
	err = h.db.QueryRowContext(ctx, getTeacherQuery, id).Scan(
		&teacher.ID, &teacher.FullName, &teacher.Username,
		&teacher.Email, &teacher.Password, &teacher.Picture,
		&teacher.Skills, &teacher.Degrees, &teacher.Experience,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve teacher"})
		return
	}

	// Clear sensitive data before sending response
	teacher.Password = ""
	c.JSON(http.StatusOK, teacher)
}

// GetAllTeachers retrieves all teachers
func (h *TeacherController) GetAllTeachers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllTeachersQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve teachers"})
		return
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var teacher models.Teacher
		if err := rows.Scan(
			&teacher.ID, &teacher.FullName, &teacher.Username,
			&teacher.Email, &teacher.Password, &teacher.Picture,
			&teacher.Skills, &teacher.Degrees, &teacher.Experience,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process teachers"})
			return
		}
		teacher.Password = "" // Clear sensitive data
		teachers = append(teachers, teacher)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing teachers"})
		return
	}

	c.JSON(http.StatusOK, teachers)
}

// UpdateTeacher updates an existing teacher
func (h *TeacherController) UpdateTeacher(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var teacher models.Teacher
	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher.ID = uint(id)
	if err := h.validateTeacher(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.checkUniqueness(ctx, &teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If password is empty, fetch the existing password
	if teacher.Password == "" {
		var currentPassword string
		err := h.db.QueryRowContext(ctx, "SELECT password FROM teachers WHERE id = $1", id).Scan(&currentPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current password"})
			return
		}
		teacher.Password = currentPassword
	}

	result, err := h.db.ExecContext(ctx, updateTeacherQuery,
		teacher.FullName, teacher.Username, teacher.Email,
		teacher.Password, teacher.Picture, teacher.Skills,
		teacher.Degrees, teacher.Experience, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update teacher"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	teacher.Password = "" // Clear sensitive data
	c.JSON(http.StatusOK, teacher)
}

// DeleteTeacher deletes a teacher by ID
func (h *TeacherController) DeleteTeacher(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteTeacherQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete teacher"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Teacher deleted successfully"})
}