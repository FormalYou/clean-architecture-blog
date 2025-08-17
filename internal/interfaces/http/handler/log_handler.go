package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"

	"github.com/FormalYou/clean-architecture-blog/internal/infrastructure/log"
)

// LogHandler handles log level changes.
type LogHandler struct{}

// NewLogHandler creates a new LogHandler.
func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

type changeLogLevelRequest struct {
	Level string `json:"level" binding:"required"`
}

// ChangeLogLevel handles the HTTP request to change the log level.
func (h *LogHandler) ChangeLogLevel(c *gin.Context) {
	var req changeLogLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(req.Level)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log level"})
		return
	}

	log.AtomicLevel.SetLevel(level)
	c.JSON(http.StatusOK, gin.H{"message": "log level changed to " + req.Level})
}
