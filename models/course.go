package models

type Course struct {
	ID          uint          `gorm:"primaryKey" json:"ID"`
	Name        string        `json:"Name"`
	Description string        `json:"Description"`
	Pricing     string        `json:"Pricing"`
	Duration    string        `json:"Duration"`
	Image       string        `json:"Image"`
	Language    string        `json:"Language"`
	Level       string        `json:"Level"`
	CourseQuizz []CourseQuizz `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
	TeacherID   uint          `json:"teacher_id"`
	Teacher     Teacher       `gorm:"foreignKey:TeacherID"`
	Student     []Student     `gorm:"many2many:student_courses;"`
	CategoryID  uint          `json:"category_id"`
	Category    Category      `gorm:"foreignKey:CategoryID"`
	Articles    []Article     `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
	Videos      []Video       `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
	Questions   []Question    `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
	Crating     []Crating     `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
}
