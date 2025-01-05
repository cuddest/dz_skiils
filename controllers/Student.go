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

// @Summary Create a new student
// @Description Register a new student in the system
// @Tags students
// @Accept json
// @Produce json
// @Param student body models.Student true "Student information"
// @Success 201 {object} models.Student
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /students/CreateStudent [post]
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

// @Summary Get student by ID
// @Description Retrieve a student's information by their ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} models.Student
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /students/GetStudent [post]
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

// @Summary Get all students
// @Description Retrieve a list of all students
// @Tags students
// @Accept json
// @Produce json
// @Success 200 {array} models.Student
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /students/all [get]
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

// @Summary Update student
// @Description Update a student's information
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param student body models.Student true "Updated student information"
// @Success 200 {object} models.Student
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /students/UpdateUser [put]
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

// @Summary Delete student
// @Description Delete a student from the system
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /students/DeleteUser [delete]
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
