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

// @Summary Create new course
// @Description Create a new course
// @Tags courses
// @Accept json
// @Produce json
// @Param course body models.Course true "Course object"
// @Success 201 {object} models.Course
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /Courses/createCourse [post]
// CreateCourse creates a new course
func (h *CourseController) CreateCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var course models.Course
	if err := c.ShouldBind(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateCourse(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id uint
	err := h.db.QueryRowContext(ctx, createCourseQuery,
		course.Name, course.Description, course.Pricing,
		course.Duration, course.Image, course.Language,
		course.Level, course.TeacherID, course.CategoryID,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	course.ID = id
	c.JSON(http.StatusCreated, course)
}

// @Summary Get course by ID
// @Description Retrieve a course by its ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} models.Course
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /Courses/get [post]
// GetCourse retrieves a course by ID
package controllers

import (
    "net/http"
    "path/filepath"
    
    "github.com/gin-gonic/gin"
    "your-project/models"
)

type CourseController struct {
    db *gorm.DB
}

func (h *CourseController) CreateCourse(c *gin.Context) {
    var course models.Course
    
    // Bind form data and file
    if err := c.ShouldBind(&course); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Save uploaded file
    if course.Image != nil {
        filename := filepath.Join("uploads", course.Image.Filename)
        if err := c.SaveUploadedFile(course.Image, filename); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed"})
            return
        }
    }

    // Create course in database
    result := h.db.Create(&course)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusCreated, course)
}

// Course Controller Swagger Documentation

// @Summary Get all courses
// @Description Retrieve a list of all courses
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {array} models.Course
// @Failure 500 {object} map[string]interface{}
// @Router /Courses/all [get]
// GetAllCourses retrieves all courses
func (h *CourseController) GetAllCourses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllCoursesQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve courses"})
		return
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(
			&course.ID, &course.Name, &course.Description,
			&course.Pricing, &course.Duration, &course.Image,
			&course.Language, &course.Level, &course.TeacherID,
			&course.CategoryID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process courses"})
			return
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// @Summary Update course
// @Description Update an existing course
// @Tags courses
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param course body models.Course true "Course object"
// @Success 200 {object} models.Course
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /Courses/updateCourse [put]
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

// @Summary Delete course
// @Description Delete a course by ID
// @Tags courses
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /Courses/DeleteCourse [delete]
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
