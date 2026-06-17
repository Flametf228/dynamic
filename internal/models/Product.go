package models

type Product struct {
	ID         int     `gorm:"primaryKey"`
	Name       string  `gorm:"not null"`
	BrandID    int     `gorm:"not null"`
	CategoryID int     `gorm:"not null"`
	Price      float64 `gorm:"type:numeric(12,2);not null"`
	Stock      int     `gorm:"not null"`
}
