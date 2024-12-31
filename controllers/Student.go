package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/cuddest/dz-skills/models"
	"github.com/gin-gonic/gin"
)

type StudentController struct {
	db *sql.DB
}

func NewStudentController(db *sql.DB) *StudentController {
	return &StudentController{db: db}
}

// Create a new student
func (h *StudentController) CreateStudent(c *gin.Context) {
	var student models.Student

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO students (full_name, email, password, picture) 
				VALUES ($1, $2, $3, $4) RETURNING id`

	var id uint
	err := h.db.QueryRow(query, student.FullName, student.Email, student.Password, student.Picture).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	student.ID = id
	c.JSON(http.StatusCreated, student)
}

func (h *StudentController) GetStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var student models.Student
	query := `SELECT id, full_name, email, password, picture FROM students WHERE id = $1`
	err = h.db.QueryRow(query, id).Scan(
		&student.ID,
		&student.FullName,
		&student.Email,
		&student.Password,
		&student.Picture,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *StudentController) GetAllStudents(c *gin.Context) {
	rows, err := h.db.Query(`SELECT id, full_name, email, password, picture FROM students`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.FullName, &student.Email, &student.Password, &student.Picture); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

// Update a student by ID
func (h *StudentController) UpdateStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE students SET full_name = $1, email = $2, password = $3, picture = $4 WHERE id = $5`
	_, err = h.db.Exec(query, student.FullName, student.Email, student.Password, student.Picture, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	student.ID = uint(id)
	c.JSON(http.StatusOK, student)
}

// Delete a student by ID
func (h *StudentController) DeleteStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	query := `DELETE FROM students WHERE id = $1`
	_, err = h.db.Exec(query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}
