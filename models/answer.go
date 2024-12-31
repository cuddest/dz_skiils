package models

type Answer struct {
	ID         uint     `gorm:"primaryKey" json:"ID"`
	Answer     string   `json:"Answer"`
	QuestionID uint     `json:"question_id"`
	Question   Question `gorm:"foreignKey:QuestionID" json:"question"`
}
