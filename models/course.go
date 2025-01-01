package models

import "time"

type Course struct {
	ID          uint          `gorm:"primaryKey" json:"ID"`
	Name        string        `json:"Name"`
	Description string        `json:"Description"`
	Pricing     string        `json:"Pricing"`
	Duration    time.Duration `json:"Duration"`
	Image       string        `json:"Image"`
	Language    string        `json:"Language"`
	Level       string        `json:"Level"`
	CourseQuizz []CourseQuizz `gorm:"foreignKey:CourseID"`
	TeacherID   uint          `json:"teacher_id"`
	Teacher     Teacher       `gorm:"foreignKey:TeacherID"`
	Student     []Student     `gorm:"many2many:student_courses;"`
	CategoryID  uint          `json:"category_id"`
	Category    Category      `gorm:"foreignKey:CategoryID"`
	Articles    []Article     `gorm:"foreignKey:CourseID"`
	Videos      []Video       `gorm:"foreignKey:CourseID"`
	Questions   []Question    `gorm:"foreignKey:CourseID"`
}
