package domain

type Status struct {
	ID     int64  `gorm:"primaryKey"`
	Title  string `gorm:"not null"`
	Active bool   `gorm:"not null"`
}
