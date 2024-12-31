package models

type Question struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	CourseID  uint     `json:"course_id"`
	StudentID uint     `json:"student_id"`
	Question  string   `json:"rating"`
	Answer    []Answer `gorm:"foreignKey:QuestionID" json:"answers"`
}
