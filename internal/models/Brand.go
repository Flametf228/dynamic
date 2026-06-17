package models

type Brand struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"unique;not null"`
}
