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

// SQL queries as constants to improve maintainability
const (
	createCourseQuery = `
		INSERT INTO courses (name, description, pricing, duration, image, language, level, teacher_id, category_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id`

	getCourseQuery = `
		SELECT id, name, description, pricing, duration, image, language, level, teacher_id, category_id 
		FROM courses 
		WHERE id = $1`

	getAllCoursesQuery = `
		SELECT id, name, description, pricing, duration, image, language, level, teacher_id, category_id 
		FROM courses`

	updateCourseQuery = `
		UPDATE courses 
		SET name = $1, description = $2, pricing = $3, duration = $4, 
			image = $5, language = $6, level = $7, teacher_id = $8, category_id = $9 
		WHERE id = $10`

	deleteCourseQuery = `DELETE FROM courses WHERE id = $1`
)

type CourseController struct {
	db *sql.DB
}

func NewCourseController(db *sql.DB) *CourseController {
	return &CourseController{db: db}
}

// validateCourse performs basic validation on course data
func (h *CourseController) validateCourse(course *models.Course) error {
	if course.Name == "" {
		return errors.New("course name is required")
	}
	if course.TeacherID <= 0 {
		return errors.New("valid teacher ID is required")
	}
	if course.CategoryID <= 0 {
		return errors.New("valid category ID is required")
	}
	return nil
}

// CreateCourse creates a new course
func (h *CourseController) CreateCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
 
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
 
	if course.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course name is required"})
		return
	}
 
	var id uint
	err := h.db.QueryRowContext(ctx, `
		INSERT INTO courses 
		(name, description, pricing, duration, image, language, level, teacher_id, category_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id`,
		course.Name, course.Description, course.Pricing,
		course.Duration,course.Image, course.Language,  course.Level, 
		course.TeacherID, course.CategoryID,
	).Scan(&id)
 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}
 
	course.ID = id
	c.JSON(http.StatusCreated, course)
 }
// GetAllCourses retrieves all courses
func (h *CourseController) GetAllCourses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Update the query to include category details
	const getAllCoursesQuery = `
		SELECT 
			c.id, c.name, c.description, c.pricing, c.duration, c.image, 
			c.language, c.level, c.teacher_id, c.category_id, 
			cat.id AS category_id, cat.name AS category_name, cat.description AS category_description
		FROM 
			courses c
		JOIN 
			categories cat ON c.category_id = cat.id
	`

	rows, err := h.db.QueryContext(ctx, getAllCoursesQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve courses"})
		return
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		var category models.Category // Assuming there's a Category struct in models

		// Scan course and category details
		if err := rows.Scan(
			&course.ID, &course.Name, &course.Description,
			&course.Pricing, &course.Duration, &course.Image,
			&course.Language, &course.Level, &course.TeacherID,
			&category.ID, &category.Name, &category.Description,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process courses"})
			return
		}

		// Attach the category object to the course
		course.Category = category
		courses = append(courses, course)
	}

	if err := rows

// UpdateCourse updates a course
func (h *CourseController) UpdateCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Check if course exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check course existence"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateCourse(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.db.ExecContext(ctx, updateCourseQuery,
		course.Name, course.Description, course.Pricing,
		course.Duration, course.Image, course.Language,
		course.Level, course.TeacherID, course.CategoryID, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	course.ID = uint(id)
	c.JSON(http.StatusOK, course)
}

// DeleteCourse deletes a course
func (h *CourseController) DeleteCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteCourseQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}