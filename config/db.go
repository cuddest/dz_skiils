package config

import (
	"fmt"
	"log"
	"os"

	"github.com/cuddest/dz-skills/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Fetch the database URL from environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set in the environment variables")
	}

	// Connect to the database using GORM
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}
	log.Println("Connected to Database!")

	// Run database migrations
	err = db.AutoMigrate(
		models.Answer{},
		models.Course{},
		models.CourseQuizz{},
		models.Category{},
		models.SubCat{},
		models.Student{},
		models.Teacher{},
		models.Article{},
		models.Video{},
		models.StudentCourse{},
		models.Crating{},
		models.Exam{},
		models.Feedback{},
		models.Question{},
		models.ExamQuizz{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	log.Println("Database Migration Completed!")

	// Assign the connected DB to the global variable
	DB = db
	return nil
}
