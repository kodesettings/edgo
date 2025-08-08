package logger

import (
	"log/slog"
	"os"
	"io"
)

type FileHandler struct {
	file    *os.File
	handler *slog.TextHandler
}

func (h *FileHandler) newFileHandler(filePath string) (*FileHandler, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FileHandler{file: file}, nil
}

func (h *FileHandler) textHandler() *slog.TextHandler {
	return slog.NewTextHandler(h.file, &slog.HandlerOptions{Level: slog.LevelInfo})
}

func (h *FileHandler) nullHandler() *slog.TextHandler {
	return slog.NewTextHandler(io.Discard, nil)
}

func (h *FileHandler) SetLogger() {
	logFileName, exists := os.LookupEnv("EDGO_LOG")
	if !exists { slog.SetDefault(slog.New(h.nullHandler())); return; }

	// Initialize the custom file handler
	fh, err := h.newFileHandler(logFileName)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	} else {
		h.file = fh.file
	}

	// Set the custom handler for the logger
	slog.SetDefault(slog.New(h.textHandler()))
}
