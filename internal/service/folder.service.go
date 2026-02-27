package service

type FolderService struct {
}

func NewFolderService() *FolderService {
	return &FolderService{}
}

func (f *FolderService) ListFolders() []string {
	return []string{"Folder 1", "Folder 2", "Folder 3"}
}
