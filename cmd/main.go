package main

import (
	"log/slog"

	"github.com/Pedro-0101/mental-dump/internal/config"
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
