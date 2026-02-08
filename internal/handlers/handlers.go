package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// uploadAvatar handles file uploads for player avatars
func uploadAvatar(c *gin.Context) {
	userID := c.PostForm("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_INPUT",
			"message": "user_id is required",
		})
		return
	}

	// Get file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_INPUT",
			"message": "avatar file is required",
		})
		return
	}

	// Validate file size (max 5MB)
	maxFileSize := int64(5 * 1024 * 1024)
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "FILE_TOO_LARGE",
			"message": "file size exceeds 5MB limit",
		})
		return
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_FILE_TYPE",
			"message": "only image files are allowed (jpeg, png, gif, webp)",
		})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "./uploads"
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	filepath := uploadDir + "/" + filename

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "UPLOAD_FAILED",
			"message": "Failed to save avatar file",
		})
		return
	}

	// Return response with avatar URL
	avatarURL := "/uploads/" + filename
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"user_id":    userID,
		"avatar_url": avatarURL,
		"message":    "Avatar uploaded successfully",
	})
}
