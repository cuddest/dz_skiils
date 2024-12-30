package models

type Exam struct {
	ID          uint        `gorm:"primaryKey" json:"ID"`
	Description string      `json:"Description"`
	ExamQuizzes []ExamQuizz `gorm:"foreignKey:ExamID"`
	CourseID    uint        `gorm:"unique" json:"course_id"`                          // Ensure that each exam is linked to one course
	Course      Course      `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE;"` // One-to-one relationship with Course
}
