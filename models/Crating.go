package models

type Crating struct {
	CourseID uint    `json:"course_id"`
	StudentID uint   `json:"student_id"`
	Rating    float64 `json:"rating"`
}