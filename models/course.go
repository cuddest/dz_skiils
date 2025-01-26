package models

type Course struct {
    ID          uint   `gorm:"primaryKey" json:"ID"`
    Name        string `json:"Name"`
    Description string `json:"Description"`
    Pricing     string `json:"Pricing"`
    Duration    string `json:"Duration"`
    Image       string `json:"Image" form:"image"`
    Language    string `json:"Language"`
    Level       string `json:"Level"`
    TeacherID   uint   `json:"teacher_id"`
    CategoryID  uint   `json:"category_id"`
    Category    Category    `gorm:"foreignKey:CategoryID"`
    Articles    []Article   `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
    Videos      []Video     `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
    Questions   []Question  `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
    Crating     []Crating   `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
}