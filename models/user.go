package models

import "errors"

type Role string

const (
	Admin   Role = "admin"
	Student Role = "student"
)

func (r Role) IsValid() error {
	switch r {
	case Admin, Student:
		return nil
	}
	return errors.New("invalid role")
}

type User struct {
	Id           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName     string `gorm:"type:varchar(300)" json:"full_name"`
	Username     string `gorm:"type:varchar(300)" json:"username"`
	Email        string `gorm:"type:varchar(300)" json:"email"`
	Password     string `gorm:"type:varchar(300)" json:"password"`
	// Role         Role   `gorm:"type:enum('student','admin');default:'student'" json:"role"`
	Role         Role   `gorm:"type:varchar(50);default:'student'" json:"role"`
	RefreshToken string `gorm:"type:text" json:"refresh_token"`
}

type Register struct {
	FullName        string `json:"full_name"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ResponseRegister struct {
	Id       uint   `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Logout struct {
	UserID uint `json:"user_id"`
}