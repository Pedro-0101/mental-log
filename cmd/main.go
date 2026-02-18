package main

import (
	"log/slog"

	"github.com/Pedro-0101/mental-log/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	slog.Info("Starting app...")

	app, err := config.NewApp()

	if err != nil {
		slog.Error("Failed to create app", "error", err)
		return
	}

	defer app.Close()

}
