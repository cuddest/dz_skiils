package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/cuddest/dz-skills/models"
	"github.com/gin-gonic/gin"
)

type CourseController struct {
	db *sql.DB
}

func NewCourseController(db *sql.DB) *CourseController {
	return &CourseController{db: db}
}

// Create a new course
func (h *CourseController) CreateCourse(c *gin.Context) {
	var course models.Course

	// Bind JSON input to the Course struct
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO courses (name, description, pricing, duration, image, language, level, teacher_id, category_id) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	var id uint
	err := h.db.QueryRow(query, course.Name, course.Description, course.Pricing, course.Duration, course.Image, course.Language, course.Level, course.TeacherID, course.CategoryID).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the generated ID and return the created course
	course.ID = id
	c.JSON(http.StatusCreated, course)
}

// Get a course by ID
func (h *CourseController) GetCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var course models.Course
	query := `SELECT id, name, description, pricing, duration, image, language, level, teacher_id, category_id 
			  FROM courses WHERE id = $1`

	err = h.db.QueryRow(query, id).Scan(
		&course.ID,
		&course.Name,
		&course.Description,
		&course.Pricing,
		&course.Duration,
		&course.Image,
		&course.Language,
		&course.Level,
		&course.TeacherID,
		&course.CategoryID,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

// Get all courses
func (h *CourseController) GetAllCourses(c *gin.Context) {
	rows, err := h.db.Query(`SELECT id, name, description, pricing, duration, image, language, level, teacher_id, category_id 
							  FROM courses`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.Pricing, &course.Duration, &course.Image, &course.Language, &course.Level, &course.TeacherID, &course.CategoryID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// Update a course by ID
func (h *CourseController) UpdateCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE courses SET name = $1, description = $2, pricing = $3, duration = $4, image = $5, language = $6, level = $7, teacher_id = $8, category_id = $9 WHERE id = $10`
	_, err = h.db.Exec(query, course.Name, course.Description, course.Pricing, course.Duration, course.Image, course.Language, course.Level, course.TeacherID, course.CategoryID, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	course.ID = uint(id)
	c.JSON(http.StatusOK, course)
}

// Delete a course by ID
func (h *CourseController) DeleteCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	query := `DELETE FROM courses WHERE id = $1`
	_, err = h.db.Exec(query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}
