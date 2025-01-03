package controllers

import (
	"net/http"

	"github.com/cuddest/dz-skills/auth"
	"github.com/cuddest/dz-skills/config"
	"github.com/cuddest/dz-skills/models"
	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func GenerateToken(context *gin.Context) {
    var input = struct {
        Email    string `json:"email"`
        Username string `json:"username"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }{}

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
        record := config.DB.Where("email = ? OR username = ?", input.Email, input.Username).First(&teacher)
        if record.Error != nil {
            context.JSON(http.StatusUnauthorized, gin.H{"error": "user not found or invalid credentials"})
            context.Abort()
            return
        }
        user = &teacher
    } else {
        var student models.Student
        record := config.DB.Where("email = ? OR username = ?", input.Email, input.Username).First(&student)
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

    // Get the concrete type's email and username
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