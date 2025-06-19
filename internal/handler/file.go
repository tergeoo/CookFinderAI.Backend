package handler

import (
	"CookFinder.Backend/internal/model"
	"CookFinder.Backend/internal/service"
	"CookFinder.Backend/internal/storage"
	"CookFinder.Backend/pkg/uuid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type FileHandler struct {
	fileService *service.FileService
	storage     *storage.YandexStorage
}

func NewFileHandler(
	r *gin.Engine,
	fileService *service.FileService,
	storage *storage.YandexStorage,
) {
	h := &FileHandler{
		fileService: fileService,
		storage:     storage,
	}
	r.POST("/upload", h.Upload)
	r.GET("/files", h.GetAll)
	r.DELETE("/files/:id", h.Delete)
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
	fileHeader, err := c.FormFile("image")
	if err != nil {
		slog.Error("failed to get uploaded file", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		slog.Error("failed to open uploaded file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer file.Close()

	url, err := it.storage.UploadFile(c, file, fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	f := &model.File{
		ID:   uuid.V7().String(),
		Name: fileHeader.Filename,
		Path: url,
	}

	err = it.fileService.CreateFile(c, f)
	if err != nil {
		slog.Error("failed to save file metadata", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"path": url})
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

	f, err := it.fileService.GetFileByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	if err := it.storage.DeleteFile(c.Request.Context(), f.Name); err != nil {
		slog.Error("failed to delete file from storage", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete from storage"})
		return
	}

	if err := it.fileService.DeleteFile(c, id); err != nil {
		slog.Error("failed to delete file metadata", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete metadata"})
		return
	}

	c.Status(http.StatusNoContent)
}
