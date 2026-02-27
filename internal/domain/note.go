package domain

import "time"

type Note struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    int64     `gorm:"not null"`
	FolderID  int64     `gorm:"not null"`
	Content   string    `gorm:"type:text"`
	Title     string    `gorm:"not null"`
	Tags      string    `gorm:"not null"`
	StatusID  int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
	Folder Folder `gorm:"foreignKey:FolderID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
	Status Status `gorm:"foreignKey:StatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}
