package models

type Article struct {
	ID          uint   `gorm:"primaryKey" json:"ID"`
	Title       string `json:"Title"`
	Link        string `json:"Link"`
	Description string `json:"Description"`
	CourseID    uint   `json:"course_id"`
	Course      Course `gorm:"foreignKey:CourseID"`
}
