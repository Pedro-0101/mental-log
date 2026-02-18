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

func (d *DBService) SaveDatabase() error {
	start := time.Now()

	// 1. git add
	cmdAdd := exec.Command("git", "add", "./mental-log.db")
	if err := cmdAdd.Run(); err != nil {
		slog.Error("git add failed", "error", err)
		return err
	}
	slog.Info("git add ./mental-log.db success")

	// 2. git commit
	commitMsg := fmt.Sprintf("Database saved %s", time.Now().Format("2006-01-02 15:04:05"))
	cmdCommit := exec.Command("git", "commit", "-m", commitMsg)
	if err := cmdCommit.Run(); err != nil {
		slog.Error("git commit failed", "error", err)
		return err
	}
	slog.Info("git commit success", "msg", commitMsg)

	// 3. git push
	cmdPush := exec.Command("git", "push", "origin", "main")
	if err := cmdPush.Run(); err != nil {
		slog.Error("git push failed", "error", err)
		return err
	}
	slog.Info("git push success")

	slog.Info("Database saved", "time", time.Since(start))
	return nil
}
