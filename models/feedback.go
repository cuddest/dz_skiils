package models

type Feedback struct {
	ID        uint   `gorm:"primaryKey" json:"ID"`
	Description string `json:"Description"`
	Review     uint   `json:"Review"`
	StudentID  uint   `json:"student_id"`
	Student    Student `gorm:"foreignKey:StudentID"`
}