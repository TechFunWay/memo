package handlers

import (
	"memo/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetVersion 获取版本信息
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": config.AppVersion,
		"name":    "备忘录",
	})
}

// GetSystemConfig 获取系统配置
func GetSystemConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data_dir":    config.DataDir,
		"upload_dir":  config.UploadDir,
		"db_dir":      config.DBDir,
	})
}
