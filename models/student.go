package models

type Student struct {
	ID       uint       `gorm:"primaryKey" json:"ID"`
	FullName string     `json:"FullName"`
	Email    string     `json:"Email"`
	Password string     `json:"Password"`
	Picture  string     `json:"Picture"`
	Courses  []Course   `gorm:"many2many:student_courses;"`
	Feedback []Feedback `gorm:"foreignKey:StudentID"`
}
