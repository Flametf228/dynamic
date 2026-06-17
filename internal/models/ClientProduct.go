package models

type ClientProduct struct {
	ClientID  int  `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
}
