package handlers

import (
	"memo/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// processStartedAt 进程启动时间，前端用 started_at 区分「本次运行」与「重启后」：
// 刷新页面不重置，应用重启后变化，赞赏弹窗据此判断是否再次提示
var processStartedAt = time.Now()

// GetVersion 获取版本信息
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    config.AppVersion,
		"name":       "备忘录",
		"started_at": processStartedAt.Format(time.RFC3339),
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
