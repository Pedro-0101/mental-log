package repo

import (
	"time"

	"github.com/Pedro-0101/mental-dump/internal/domain"
	"gorm.io/gorm"
)

type FolderRepo struct {
	DB *gorm.DB
}

func NewFolderRepo(db *gorm.DB) *FolderRepo {
	return &FolderRepo{DB: db}
}

func (r *FolderRepo) FindAll() ([]domain.Folder, error) {
	var folders []domain.Folder
	if err := r.DB.Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

// FindRootFolders returns folders that have no parent (root-level).
func (r *FolderRepo) FindRootFolders() ([]domain.Folder, error) {
	var folders []domain.Folder
	if err := r.DB.Where("parent_id IS NULL").Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

// FindByParentID returns subfolders of a given parent folder.
func (r *FolderRepo) FindByParentID(parentID int64) ([]domain.Folder, error) {
	var folders []domain.Folder
	if err := r.DB.Where("parent_id = ?", parentID).Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

func (r *FolderRepo) Create(title string, userID int64, tags string, parentID *int64) (*domain.Folder, error) {
	newFolder := &domain.Folder{
		UserID:    userID,
		ParentID:  parentID,
		Title:     title,
		Tags:      tags,
		StatusID:  1,
		CreatedAt: time.Now(),
	}

	if err := r.DB.Create(newFolder).Error; err != nil {
		return nil, err
	}

	return newFolder, nil
}

func (r *FolderRepo) FindByID(id int64) (*domain.Folder, error) {
	var folder domain.Folder
	if err := r.DB.First(&folder, id).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *FolderRepo) UpdateTags(id int64, tags string) error {
	return r.DB.Model(&domain.Folder{}).Where("id = ?", id).Update("tags", tags).Error
}
