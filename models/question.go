package models

type Question struct {
	CourseID  uint     `json:"course_id"`
	StudentID uint     `json:"student_id"`
	Question  string   `json:"rating"`
	Answer    []Answer `gorm:"foreignKey:QuestionID"`
}
