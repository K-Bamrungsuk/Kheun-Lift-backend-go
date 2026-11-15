package middleware

import (
	"encoding/json"
	"mime"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DecodeBody reads a JSON body (16 KB at most, no unknown fields). On failure it answers and returns false.
func DecodeBody(c *gin.Context, input any) bool {
	contentType, _, err := mime.ParseMediaType(c.GetHeader("Content-Type"))
	if err != nil || contentType != "application/json" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"message": "Use application/json."})
		return false
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16*1024)
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON, unsupported fields or request body too large."})
		return false
	}
	return true
}

// ParseID reads a positive 32-bit integer URL parameter. On failure it answers 400 and returns false.
func ParseID(c *gin.Context, name string) (int32, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 32)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": name + " must be a positive 32-bit integer."})
		return 0, false
	}
	return int32(id), true
}
