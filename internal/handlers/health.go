package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (*HealthHandler) HealthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"date":   time.Now(),
	})
}
