package domain

import "time"

type Folder struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64 `gorm:"not null"`
	ParentID  *int64
	Title     string `gorm:"not null"`
	Tags      string
	StatusID  int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`

	User       User     `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
	Parent     *Folder  `gorm:"foreignKey:ParentID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
	Subfolders []Folder `gorm:"foreignKey:ParentID"`
	Status     Status   `gorm:"foreignKey:StatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}
