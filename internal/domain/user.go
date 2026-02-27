package domain

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	StatusID  int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`

	Status Status `gorm:"foreignKey:StatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}
