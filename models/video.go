package models

type Video struct {
	ID         uint   `gorm:"primaryKey" json:"ID"`
	Title      string `json:"Title"`
	Link       string `json:"Link"`
	CourseID   uint   `json:"course_id"`
	Course     Course `gorm:"foreignKey:CourseID"`
}