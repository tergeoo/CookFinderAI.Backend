package handler

import (
	"CookFinder.Backend/internal/model"
	"CookFinder.Backend/internal/service"
	"CookFinder.Backend/pkg/uuid"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(
	r *gin.Engine,
	fileService *service.FileService,
) {
	h := &FileHandler{
		fileService: fileService,
	}
	r.POST("/upload", h.Upload)
	r.GET("/files", h.GetAll)
	r.Static("/static", "./uploads")
}

// Upload godoc
// @Summary Upload image file
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image File"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /upload [post]
func (it *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		slog.Error("failed to get uploaded file", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), sanitizeFilename(file.Filename))
	path := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		slog.Error("failed to save uploaded file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	url := "/static/" + filename

	f := &model.File{
		ID:   uuid.V7().String(),
		Name: file.Filename,
		Path: url,
	}

	err = it.fileService.CreateFile(c, f)
	if err != nil {
		slog.Error("failed to save file metadata", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAll godoc
// @Summary GetAll files
// @Tags Files
// @Produce json
// @Success 204 {string} string "No Content"
// @Failure 500 {object} map[string]string
// @Router /files [get]
func (it *FileHandler) GetAll(c *gin.Context) {
	files, err := it.fileService.GetAllFiles(c)
	if err != nil {
		slog.Error("failed to get files", "error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// Delete godoc
// @Summary Delete by id
// @Tags Files
// @Param id path string true "File id"
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /files/{id} [delete]
func (it *FileHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Получаем метаинформацию о файле из БД
	f, err := it.fileService.GetFileByID(c, id)
	if err != nil {
		slog.Error("failed to get file", "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// Удаляем сам файл с диска
	filePath := filepath.Join("uploads", filepath.Base(f.Path)) // безопасное соединение пути
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		slog.Error("failed to delete file from disk", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file"})
		return
	}

	// Удаляем метаинформацию из БД
	if err := it.fileService.DeleteFile(c, id); err != nil {
		slog.Error("failed to delete file metadata", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file metadata"})
		return
	}

	c.Status(http.StatusNoContent)
}

func sanitizeFilename(name string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' || r == '-' || r == '_' || r == '.' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, name)
}
