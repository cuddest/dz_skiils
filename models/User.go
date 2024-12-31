package models

import "golang.org/x/crypto/bcrypt"

type Student struct {
	ID       uint       `gorm:"primaryKey" json:"ID"`
	FullName string     `json:"FullName"`
	Email    string     `json:"Email"`
	Password string     `json:"Password"`
	Picture  string     `json:"Picture"`
	Courses  []Course   `gorm:"many2many:student_courses;"`
	Feedback []Feedback `gorm:"foreignKey:StudentID"`
}

func (s *Student) GetPassword() string {
	return s.Password
}

func (s *Student) SetPassword(password string) {
	s.Password = password
}

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

func (t *Teacher) GetPassword() string {
	return t.Password
}

func (t *Teacher) SetPassword(password string) {
	t.Password = password
}

type User interface {
	GetPassword() string
	SetPassword(password string)
}

func HashPassword(user User, password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.SetPassword(string(bytes))
	return nil
}

func CheckPassword(user User, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
