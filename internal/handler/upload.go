package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
}

func NewUploadHandler(r *gin.Engine) {
	h := &UploadHandler{}
	r.POST("/upload", h.Upload)
	r.Static("/static", "./uploads")
}

// Upload godoc
// @Summary Upload image file
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image File"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), sanitizeFilename(file.Filename))
	path := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	url := "/static/" + filename
	c.JSON(http.StatusOK, gin.H{"image_url": url})
}

func sanitizeFilename(name string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' || r == '-' || r == '_' || r == '.' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, name)
}
