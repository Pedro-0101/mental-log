package domain

import "time"

type Entry struct {
	ID         int64  `gorm:"primaryKey"`
	Content    string `gorm:"type:text;not null"`
	UserID     int64  `gorm:"not null"`
	EntryType  string
	Title      string
	IsEvent    bool `gorm:"default:false"`
	ScheduleAt *time.Time
	StatusID   int64 `gorm:"not null"`
	NoteID     *int64
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
	Note   Note   `gorm:"foreignKey:NoteID;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL;"`
	Status Status `gorm:"foreignKey:StatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}
