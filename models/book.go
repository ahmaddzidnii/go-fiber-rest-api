package models

type Book struct {
	Id          int64  `json:"id" gorm:"primaryKey"`
	Title       string `gorm:"type:varchar(300)" json:"title"`
	Description string `gorm:"text" json:"description"`
	Author      string `gorm:"type:varchar(300)" json:"author"`
	PublishDate string `gorm:"type:date" json:"publish_date"`
}