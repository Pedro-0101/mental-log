package service

import (
	"log/slog"
	"strings"

	"github.com/Pedro-0101/mental-dump/internal/domain"
	"github.com/Pedro-0101/mental-dump/internal/repo"
)

type FolderService struct {
	folderRepo *repo.FolderRepo
}

func NewFolderService(folderRepo *repo.FolderRepo) *FolderService {
	return &FolderService{folderRepo: folderRepo}
}

func (f *FolderService) FindAll() ([]domain.Folder, error) {
	return f.folderRepo.FindAll()
}

// FindRootFolders returns top-level folders (no parent).
func (f *FolderService) FindRootFolders() ([]domain.Folder, error) {
	return f.folderRepo.FindRootFolders()
}

// FindByParentID returns subfolders of a given parent folder.
func (f *FolderService) FindByParentID(parentID int64) ([]domain.Folder, error) {
	return f.folderRepo.FindByParentID(parentID)
}

func (f *FolderService) FindByID(id int64) (*domain.Folder, error) {
	return f.folderRepo.FindByID(id)
}

func (f *FolderService) CreateFolder(title string, parentID *int64) (*domain.Folder, error) {
	mockedUserId := int64(1)
	tags := ""
	return f.folderRepo.Create(title, mockedUserId, tags, parentID)
}

// PropagateTags walks up the parent chain from folderID, merging noteTags
// into each ancestor folder's Tags (comma-separated, skipping duplicates).
func (f *FolderService) PropagateTags(folderID *int64, noteTags string) {
	if folderID == nil || noteTags == "" {
		return
	}

	newTags := splitTags(noteTags)
	currentID := folderID

	for currentID != nil {
		folder, err := f.folderRepo.FindByID(*currentID)
		if err != nil {
			slog.Error("Error finding folder for tag propagation", "id", *currentID, "error", err)
			return
		}

		existingTags := splitTags(folder.Tags)
		tagSet := make(map[string]bool)
		for _, t := range existingTags {
			tagSet[t] = true
		}

		changed := false
		for _, t := range newTags {
			if !tagSet[t] {
				existingTags = append(existingTags, t)
				tagSet[t] = true
				changed = true
			}
		}

		if changed {
			merged := strings.Join(existingTags, ",")
			if err := f.folderRepo.UpdateTags(folder.ID, merged); err != nil {
				slog.Error("Error updating folder tags", "id", folder.ID, "error", err)
				return
			}
			slog.Info("Propagated tags to folder", "id", folder.ID, "tags", merged)
		}

		currentID = folder.ParentID
	}
}

func splitTags(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
