package models

import (
	"time"
)

type StudentCourse struct {
	StudentID   uint      `gorm:"primaryKey"`
	CourseID    uint      `gorm:"primaryKey"`
	Grade       string    `json:"grade"`
	Enrollment  time.Time `json:"enrollment"`
	Certificate *string   `json:"certificate"`
	Issued      bool      `json:"issued"`
}
