package handlers

import (
	"net/http"

	"github.com/Yw332/campus-moments-go/pkg/database"
	"github.com/gin-gonic/gin"
)

// Home 根路径
func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"message": "Campus Moments Go API 运行中",
			"version": "1.0.0",
			"author":  "Yw332",
			"status":  "healthy",
		},
	})
}

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	status := "ok"
	dbStatus := "disconnected"
	if database.DB != nil {
		if err := database.DB.Ping(); err == nil {
			dbStatus = "connected"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"status":   status,
			"service":  "Campus Moments Go API",
			"database": dbStatus,
		},
	})
}
