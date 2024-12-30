package models

type SubCat struct {
	ID         uint      `gorm:"primaryKey" json:"ID"`
	Name       string    `json:"Name"`
	CategoryID uint      `json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID"`
}