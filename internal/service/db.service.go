package service

import (
	"fmt"
	"log/slog"
	"os/exec"
	"time"
)

type DBService struct {
}

func NewDBService() *DBService {
	return &DBService{}
}

/*func (d *DBService) startSave() {
	go d.SaveDatabase()
	slog.Info("Go rountine initialized")
}*/

func (d *DBService) SaveDatabase() {
	start := time.Now()

	cmdAdd := exec.Command("git", "add", "./mental-log.db")
	if err := cmdAdd.Run(); err != nil {
		slog.Error("git add failed", "error", err)
	}
	slog.Info("git add ./mental-log.db success")

	commitMsg := fmt.Sprintf("Database saved %s", time.Now().Format("2006-01-02 15:04:05"))
	cmdCommit := exec.Command("git", "commit", "-m", commitMsg)
	if err := cmdCommit.Run(); err != nil {
		slog.Error("git commit failed", "error", err)
	}
	slog.Info("git commit success", "msg", commitMsg)

	cmdPush := exec.Command("git", "push", "origin", "main")
	if err := cmdPush.Run(); err != nil {
		slog.Error("git push failed", "error", err)
	}
	slog.Info("git push success")

	slog.Info("Database saved", "time", time.Since(start))
}
