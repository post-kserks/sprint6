package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"morse-converter/internal/service"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	workDir, err := os.Getwd()
	if err != nil {
		h.logger.Printf("Error getting working directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	indexPath := filepath.Join(workDir, "index.html")

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		h.logger.Printf("File not found: %s", indexPath)
		http.Error(w, "File not found", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, indexPath)
}

func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Printf("Error getting file: %v", err)
		http.Error(w, "Error getting file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		h.logger.Printf("Error reading file: %v", err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	result, err := service.ConvertText(string(content))
	if err != nil {
		h.logger.Printf("Error converting text: %v", err)
		http.Error(w, "Error converting text", http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(header.Filename)
	outputFilename := fmt.Sprintf("result_%s%s", time.Now().UTC().String(), ext)

	err = os.WriteFile(outputFilename, []byte(result), 0644)
	if err != nil {
		h.logger.Printf("Error writing result file: %v", err)
		http.Error(w, "Error writing result file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(result))
}
