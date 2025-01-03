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

// SQL queries for Crating
const (
	createCratingQuery = `
		INSERT INTO cratings (course_id, student_id, rating) 
		VALUES ($1, $2, $3)`

	getAverageRatingByCourseIDQuery = `
        SELECT COALESCE(AVG(rating), 0) as average_rating, COUNT(*) as total_ratings
        FROM cratings 
        WHERE course_id = $1`

	getCratingByCourseIDQuery = `
		SELECT course_id, student_id, rating 
		FROM cratings WHERE course_id = $1`

	getCratingByStudentIDQuery = `
		SELECT course_id, student_id, rating 
		FROM cratings WHERE student_id = $1`

	getCratingByCourseAndStudentIDQuery = `
		SELECT course_id, student_id, rating 
		FROM cratings WHERE course_id = $1 AND student_id = $2`

	getAllCratingsQuery = `
		SELECT course_id, student_id, rating 
		FROM cratings`

	updateCratingQuery = `
		UPDATE cratings 
		SET rating = $1 
		WHERE course_id = $2 AND student_id = $3`

	deleteCratingQuery = `
		DELETE FROM cratings WHERE course_id = $1 AND student_id = $2`
)

type CratingController struct {
	db *sql.DB
}

func NewCratingController(db *sql.DB) *CratingController {
	return &CratingController{db: db}
}

func (h *CratingController) validateCrating(crating *models.Crating) error {
	if crating.CourseID == 0 || crating.StudentID == 0 {
		return errors.New("course_id and student_id are required")
	}
	if crating.Rating < 0 || crating.Rating > 5 {
		return errors.New("rating must be between 0 and 5")
	}
	return nil
}

// CreateCrating handles the creation of a new rating
func (h *CratingController) CreateCrating(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var crating models.Crating
	if err := c.ShouldBindJSON(&crating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateCrating(&crating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.ExecContext(ctx, createCratingQuery, crating.CourseID, crating.StudentID, crating.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rating"})
		return
	}

	c.JSON(http.StatusCreated, crating)
}

// GetCratingsByCourse retrieves all ratings for a specific course
func (h *CratingController) GetCratingsByCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	rows, err := h.db.QueryContext(ctx, getCratingByCourseIDQuery, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ratings"})
		return
	}
	defer rows.Close()

	var cratings []models.Crating
	for rows.Next() {
		var crating models.Crating
		if err := rows.Scan(&crating.CourseID, &crating.StudentID, &crating.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process ratings"})
			return
		}
		cratings = append(cratings, crating)
	}

	c.JSON(http.StatusOK, cratings)
}

// GetCratingsByStudent retrieves all ratings for a specific student
func (h *CratingController) GetCratingsByStudent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	studentID, err := strconv.Atoi(c.Param("student_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID format"})
		return
	}

	rows, err := h.db.QueryContext(ctx, getCratingByStudentIDQuery, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ratings"})
		return
	}
	defer rows.Close()

	var cratings []models.Crating
	for rows.Next() {
		var crating models.Crating
		if err := rows.Scan(&crating.CourseID, &crating.StudentID, &crating.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process ratings"})
			return
		}
		cratings = append(cratings, crating)
	}

	c.JSON(http.StatusOK, cratings)
}

// GetCratingByCourseAndStudent retrieves a specific rating by course ID and student ID
func (h *CratingController) GetCratingByCourseAndStudent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	studentID, err := strconv.Atoi(c.Param("student_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID format"})
		return
	}

	var crating models.Crating
	err = h.db.QueryRowContext(ctx, getCratingByCourseAndStudentIDQuery, courseID, studentID).Scan(&crating.CourseID, &crating.StudentID, &crating.Rating)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rating"})
		return
	}

	c.JSON(http.StatusOK, crating)
}

// UpdateCrating updates an existing rating
func (h *CratingController) UpdateCrating(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var crating models.Crating
	if err := c.ShouldBindJSON(&crating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateCrating(&crating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.ExecContext(ctx, updateCratingQuery, crating.Rating, crating.CourseID, crating.StudentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating updated successfully"})
}

// DeleteCrating deletes a rating by course ID and student ID
func (h *CratingController) DeleteCrating(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	studentID, err := strconv.Atoi(c.Param("student_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID format"})
		return
	}

	_, err = h.db.ExecContext(ctx, deleteCratingQuery, courseID, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating deleted successfully"})
}
func (h *CratingController) GetCourseAverageRating(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	var averageRating float64
	var totalRatings int

	err = h.db.QueryRowContext(ctx, getAverageRatingByCourseIDQuery, courseID).Scan(&averageRating, &totalRatings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average rating"})
		return
	}

	response := gin.H{
		"course_id":      courseID,
		"average_rating": averageRating,
		"total_ratings":  totalRatings,
	}

	c.JSON(http.StatusOK, response)
}

func (h *CratingController) GetAllCratings(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllCratingsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ratings"})
		return
	}
	defer rows.Close()

	var cratings []models.Crating
	for rows.Next() {
		var crating models.Crating
		if err := rows.Scan(&crating.CourseID, &crating.StudentID, &crating.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process ratings"})
			return
		}
		cratings = append(cratings, crating)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating through ratings"})
		return
	}

	c.JSON(http.StatusOK, cratings)
}
