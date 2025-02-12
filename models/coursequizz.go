package models

type CourseQuizz struct {
	ID       uint   `gorm:"primaryKey" json:"ID"`
	Question string `json:"Question"`
	Option1  string `json:"Option1"`
	Option2  string `json:"Option2"`
	Option3  string `json:"Option3"`
	Option4  string `json:"Option4"`
	Answer   string `json:"Answer"`
	CourseID uint   `json:"exam_id"`
	Course   Course `gorm:"foreignKey:CourseID"`
}
