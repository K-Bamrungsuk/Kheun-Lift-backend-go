package handler

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"kheun-Lift-backend-go/internal/middleware"
	"kheun-Lift-backend-go/internal/user/model"
	"kheun-Lift-backend-go/internal/user/repository"
	"kheun-Lift-backend-go/internal/user/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxProfileImageSize = 2 << 20 // 2 MB

// NewUploadProfileImageHandler takes a multipart file in the "profileImage" field.
func NewUploadProfileImageHandler(pool *pgxpool.Pool, store *repository.ProfileImageStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, ok := middleware.SessionFromContext(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Please login."})
			return
		}

		// Extra 4 KB is room for the multipart headers around the file.
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxProfileImageSize+4096)
		file, _, err := c.Request.FormFile("profileImage")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Send one image up to 2 MB in the profileImage field."})
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil || len(data) > maxProfileImageSize {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Send one image up to 2 MB in the profileImage field."})
			return
		}

		// Check the bytes, not the file name, so only real images reach the public bucket.
		contentType := http.DetectContentType(data)
		if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Profile image must be JPEG, PNG or WebP."})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()

		url, err := store.Put(ctx, session.User.ID, contentType, data)
		if err != nil {
			log.Printf("Upload profile image failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to upload profile image."})
			return
		}

		user, err := service.UpdateUser(ctx, pool, session.User.ID, model.UpdateUserInput{ProfileImage: &url})
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Session is invalid. Please login."})
			return
		}
		if err != nil {
			log.Print("Save profile image failed")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to update user."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func NewDeleteProfileImageHandler(pool *pgxpool.Pool, store *repository.ProfileImageStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, ok := middleware.SessionFromContext(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Please login."})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()

		// Delete from R2 first: a 200 must mean the photo is really gone.
		// R2 does not error on a missing key, so calling this without an image is fine.
		if err := store.Delete(ctx, session.User.ID); err != nil {
			log.Printf("Delete profile image failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to delete profile image."})
			return
		}

		err := repository.ClearProfileImage(ctx, pool, session.User.ID)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Session is invalid. Please login."})
			return
		}
		if err != nil {
			log.Print("Clear profile image failed")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to delete profile image."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Profile image deleted."})
	}
}
