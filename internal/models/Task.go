package models

import (
	"time"
)

type Task struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	StartTime time.Time `gorm:"autoCreateTime"`
	EndTime   *time.Time
	Status    string `gorm:"type:varchar(20);default:'in_progress'"`
}
