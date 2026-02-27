package folder

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Pedro-0101/mental-dump/internal/service"
)

type ListFolder struct {
	service *service.FolderService
	width   float32
	height  float32
	folders []string
}

func NewListFolder(service *service.FolderService, width float32, height float32) *ListFolder {
	folders := service.ListFolders()
	return &ListFolder{service: service, width: width, height: height, folders: folders}
}

func (l *ListFolder) RenderList() fyne.CanvasObject {

	if len(l.folders) == 0 {
		return widget.NewLabel("No folders")
	}

	var buttons []fyne.CanvasObject
	for _, folder := range l.folders {

		folderButton := widget.NewButton(folder, func() {

			slog.Info("Folder selected", "folder", folder)

		})

		folderButton.Resize(fyne.NewSize(l.width, l.height))
		buttons = append(buttons, folderButton)
	}

	return container.NewVBox(buttons...)

}
