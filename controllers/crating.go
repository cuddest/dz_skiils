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

// @Summary Create new rating
// @Description Create a new course rating
// @Tags ratings
// @Accept json
// @Produce json
// @Param rating body models.Crating true "Rating object"
// @Success 201 {object} models.Crating
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cratings/createCrating [post]
// CreateCrating creates a new rating
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

// @Summary Get ratings by course
// @Description Retrieve all ratings for a specific course
// @Tags ratings
// @Accept json
// @Produce json
// @Param course_id path int true "Course ID"
// @Success 200 {array} models.Crating
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cratings/GetCratingsByCourse [post]
// GetCratingsByCourse retrieves ratings by course
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

// @Summary Get ratings by student
// @Description Retrieve all ratings for a specific student
// @Tags ratings
// @Accept json
// @Produce json
// @Param student_id path int true "Student ID"
// @Success 200 {array} models.Crating
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cratings/GetCratingsByStudent [post]
// GetCratingsByStudent retrieves ratings by student
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

// @Summary Get rating by course and student
// @Description Retrieve a specific rating by course and student IDs
// @Tags ratings
// @Accept json
// @Produce json
// @Param course_id path int true "Course ID"
// @Param student_id path int true "Student ID"
// @Success 200 {object} models.Crating
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cratings/GetCratingByCourseAndStudent [post]
// GetCratingByCourseAndStudent retrieves a specific rating
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

// @Summary Update rating
// @Description Update an existing rating
// @Tags ratings
// @Accept json
// @Produce json
// @Param rating body models.Crating true "Rating object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cratings/updateCrating [put]
// UpdateCrating updates a rating
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

// @Summary Delete rating
// @Description Delete a rating by course and student IDs
// @Tags ratings
// @Accept json
// @Produce json
// @Param course_id path int true "Course ID"
// @Param student_id path int true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cratings/DeleteCrating [delete]
// DeleteCrating deletes a rating

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

// @Summary Get course average rating
// @Description Get the average rating for a specific course
// @Tags ratings
// @Accept json
// @Produce json
// @Param course_id path int true "Course ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cratings/GetCourseAverageRating [post]
// GetCourseAverageRating retrieves average rating for a course
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

// @Summary Get all ratings
// @Description Retrieve all course ratings
// @Tags ratings
// @Accept json
// @Produce json
// @Success 200 {array} models.Crating
// @Failure 500 {object} map[string]interface{}
// @Router /cratings/GetAllCratings [get]
// GetAllCratings retrieves all ratings

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
