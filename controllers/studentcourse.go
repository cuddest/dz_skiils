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

const (
	createStudentCourseQuery = `
		INSERT INTO student_courses (student_id, course_id, grade, enrollment, certificate, issued) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	getStudentCourseQuery = `
		SELECT student_id, course_id, grade, enrollment, certificate, issued 
		FROM student_courses 
		WHERE student_id = $1 AND course_id = $2`

	getAllStudentCoursesQuery = `
		SELECT student_id, course_id, grade, enrollment, certificate, issued 
		FROM student_courses`

	updateStudentCourseQuery = `
		UPDATE student_courses 
		SET grade = $1, enrollment = $2, certificate = $3, issued = $4
		WHERE student_id = $5 AND course_id = $6`

	deleteStudentCourseQuery = `
		DELETE FROM student_courses 
		WHERE student_id = $1 AND course_id = $2`

	getExamQuizzAnswerQuery = `
		SELECT answer 
		FROM exam_quizzes 
		WHERE id = $1`
)

// StudentCourseController handles HTTP requests for StudentCourse operations
type StudentCourseController struct {
	db *sql.DB
}

// ExamAnswer represents a student's answer to an exam question
type ExamAnswer struct {
	QuizzID uint `json:"quizz_id"`
	Answer  uint `json:"answer"`
}

// NewStudentCourseController creates a new StudentCourseController instance
func NewStudentCourseController(db *sql.DB) *StudentCourseController {
	return &StudentCourseController{db: db}
}

// validateStudentCourse performs validation on student course data
func (h *StudentCourseController) validateStudentCourse(sc *models.StudentCourse) error {
	if sc.StudentID == 0 {
		return errors.New("student ID is required")
	}
	if sc.CourseID == 0 {
		return errors.New("course ID is required")
	}
	return nil
}

// @Summary Create student course enrollment
// @Description Create a new student course enrollment
// @Tags student-courses
// @Accept json
// @Produce json
// @Param studentCourse body models.StudentCourse true "Student course enrollment information"
// @Success 201 {object} models.StudentCourse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map/ @Summary Get student course enrollment
// @Description Retrieve a specific student course enrollment
// @Tags student-courses
// @Accept json
// @Produce json
// @Param studentId path int true "Student ID"
// @Param courseId path int true "Course ID"
// @Success 200 {object} models.StudentCourse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth[string]interface{}
// @Security ApiKeyAuth
// @Router /student_courses/createStudentCourse [post]
func (h *StudentCourseController) CreateStudentCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var sc models.StudentCourse
	if err := c.ShouldBindJSON(&sc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateStudentCourse(&sc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values for new enrollment
	sc.Enrollment = time.Now()
	sc.Issued = false

	_, err := h.db.ExecContext(ctx, createStudentCourseQuery,
		sc.StudentID, sc.CourseID, sc.Grade, sc.Enrollment, sc.Certificate, sc.Issued)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student course enrollment"})
		return
	}

	c.JSON(http.StatusCreated, sc)
}

// @Summary Get student course enrollment
// @Description Retrieve a specific student course enrollment
// @Tags student-courses
// @Accept json
// @Produce json
// @Param studentId path int true "Student ID"
// @Param courseId path int true "Course ID"
// @Success 200 {object} models.StudentCourse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /student_courses/get [post]
func (h *StudentCourseController) GetStudentCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	studentID, err := strconv.ParseUint(c.Param("studentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID format"})
		return
	}

	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	var sc models.StudentCourse
	err = h.db.QueryRowContext(ctx, getStudentCourseQuery, studentID, courseID).Scan(
		&sc.StudentID, &sc.CourseID, &sc.Grade, &sc.Enrollment, &sc.Certificate, &sc.Issued,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student course enrollment not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve student course enrollment"})
		return
	}

	c.JSON(http.StatusOK, sc)
}

// @Summary Get all student course enrollments
// @Description Retrieve all student course enrollments
// @Tags student-courses
// @Accept json
// @Produce json
// @Success 200 {array} models.StudentCourse
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /student_courses/all [get]

func (h *StudentCourseController) GetAllStudentCourses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllStudentCoursesQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve student course enrollments"})
		return
	}
	defer rows.Close()

	var studentCourses []models.StudentCourse
	for rows.Next() {
		var sc models.StudentCourse
		if err := rows.Scan(
			&sc.StudentID, &sc.CourseID, &sc.Grade,
			&sc.Enrollment, &sc.Certificate, &sc.Issued,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process student course enrollments"})
			return
		}
		studentCourses = append(studentCourses, sc)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing student course enrollments"})
		return
	}

	c.JSON(http.StatusOK, studentCourses)
}

// @Summary Submit exam answers
// @Description Submit and grade student exam answers
// @Tags student-courses
// @Accept json
// @Produce json
// @Param studentId path int true "Student ID"
// @Param courseId path int true "Course ID"
// @Param answers body []ExamAnswer true "Array of exam answers"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /student_courses/SubmitExamAnswers [post]
func (h *StudentCourseController) SubmitExamAnswers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()

	// Get student and course IDs from parameters
	studentID, err := strconv.ParseUint(c.Param("studentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID format"})
		return
	}

	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	// Get answers from request body
	var answers []ExamAnswer
	if err := c.ShouldBindJSON(&answers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate number of answers
	if len(answers) != 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exactly 20 answers are required"})
		return
	}

	// Start transaction
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Grade the exam
	var correctAnswers uint = 0
	for _, answer := range answers {
		var correctAnswer uint
		err := tx.QueryRowContext(ctx, getExamQuizzAnswerQuery, answer.QuizzID).Scan(&correctAnswer)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found: " + strconv.FormatUint(uint64(answer.QuizzID), 10)})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve correct answer"})
			return
		}

		if answer.Answer == correctAnswer {
			correctAnswers++
		}
	}

	// Calculate grade out of 20
	grade := strconv.FormatUint(uint64(correctAnswers), 10) + "/20"

	// Determine if student passed (>= 10/20)
	passed := correctAnswers >= 10
	var certificate *string
	if passed {
		certText := "System of certificates available soon"
		certificate = &certText
	}

	// Update student course record
	_, err = tx.ExecContext(ctx, updateStudentCourseQuery,
		grade, time.Now(), certificate, passed,
		studentID, courseID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student course record"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Prepare response
	response := gin.H{
		"grade":              grade,
		"passed":             passed,
		"certificate_issued": passed,
	}
	if passed {
		response["certificate"] = *certificate
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update student course enrollment
// @Description Update an existing student course enrollment
// @Tags student-courses
// @Accept json
// @Produce json
// @Param studentId path int true "Student ID"
// @Param courseId path int true "Course ID"
// @Param studentCourse body models.StudentCourse true "Updated student course information"
// @Success 200 {object} models.StudentCourse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /student_courses/updateStudentCourse [put]
func (h *StudentCourseController) UpdateStudentCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	studentID, err := strconv.ParseUint(c.Param("studentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID format"})
		return
	}

	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	var sc models.StudentCourse
	if err := c.ShouldBindJSON(&sc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sc.StudentID = uint(studentID)
	sc.CourseID = uint(courseID)

	result, err := h.db.ExecContext(ctx, updateStudentCourseQuery,
		sc.Grade, sc.Enrollment, sc.Certificate, sc.Issued,
		sc.StudentID, sc.CourseID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student course enrollment"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student course enrollment not found"})
		return
	}

	c.JSON(http.StatusOK, sc)
}

// @Summary Delete student course enrollment
// @Description Delete a student course enrollment
// @Tags student-courses
// @Accept json
// @Produce json
// @Param studentId path int true "Student ID"
// @Param courseId path int true "Course ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /student_courses/DeleteStudentCourse [delete]
func (h *StudentCourseController) DeleteStudentCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	studentID, err := strconv.ParseUint(c.Param("studentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID format"})
		return
	}

	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteStudentCourseQuery, studentID, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student course enrollment"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student course enrollment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student course enrollment deleted successfully"})
}
