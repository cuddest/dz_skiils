package controllers

import (
	"net/http"

	"github.com/cuddest/dz-skills/auth"
	"github.com/cuddest/dz-skills/config"
	"github.com/cuddest/dz-skills/models"
	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Identifier string `json:"email"` // Can be either email or username
	Password   string `json:"password"`
	Role       string `json:"role"`
}

// @Summary User login
// @Description Authenticate a user (teacher or student) and generate an access token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body TokenRequest true "Login credentials (identifier (email or username) and password)"
// @Success 200 {object} map[string]interface{} "Returns JWT token"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Authentication failed"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /teachers/login [post]
// @Router /students/login [post]
func GenerateToken(context *gin.Context) {
	var input TokenRequest

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Validate role
	if input.Role != "teacher" && input.Role != "student" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid role specified"})
		context.Abort()
		return
	}

	var user models.User
	if input.Role == "teacher" {
		var teacher models.Teacher
		record := config.DB.Where("email = ? OR username = ?", input.Identifier, input.Identifier).First(&teacher)
		if record.Error != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "user not found or invalid credentials"})
			context.Abort()
			return
		}
		user = &teacher
	} else {
		var student models.Student
		record := config.DB.Where("email = ? OR username = ?", input.Identifier, input.Identifier).First(&student)
		if record.Error != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "user not found or invalid credentials"})
			context.Abort()
			return
		}
		user = &student
	}

	credentialError := models.CheckPassword(user, input.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}

	var email, username string
	switch v := user.(type) {
	case *models.Teacher:
		email = v.Email
		username = v.Username
	case *models.Student:
		email = v.Email
		username = v.Username
	}

	tokenString, err := auth.GenerateJWT(email, username, input.Role)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"token":    tokenString,
		"username": username,
		"role":     input.Role,
	})
}
