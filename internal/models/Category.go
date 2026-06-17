package models

type Category struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"unique;not null"`
}
