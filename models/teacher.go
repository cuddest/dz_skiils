package models

type Teacher struct {
	ID         uint   `gorm:"primaryKey" json:"ID"`
	FullName   string `json:"FullName"`
	Email      string `json:"Email"`
	Password   string `json:"Password"`
	Picture    string `json:"Picture"`
	Skills     string `json:"Skills"`
	Degrees    string `json:"Degree"`
	Experience string `json:"Experience"`
}
