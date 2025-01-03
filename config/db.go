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

func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set in the environment variables")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}
	log.Println("Connected to Database!")

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}
	log.Println("Database Migration Completed!")
	DB = db
	return db, nil
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Answer{},
		&models.Course{},
		&models.CourseQuizz{},
		&models.Category{},
		&models.SubCat{},
		&models.Student{},
		&models.Teacher{},
		&models.Article{},
		&models.Video{},
		&models.StudentCourse{},
		&models.Crating{},
		&models.Exam{},
		&models.Feedback{},
		&models.Question{},
		&models.ExamQuizz{},
	)
}
