package models

type Category struct {
	ID      uint     `gorm:"primaryKey" json:"ID"`
	Name    string   `json:"Name"`
	SubCats []SubCat `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	Courses []Course `gorm:"foreignKey:CategoryID"`
}
