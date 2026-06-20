package handlers

import (
	"memo/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetConfigs 获取所有配置
func GetConfigs(c *gin.Context) {
	configService := &services.ConfigService{}
	configs, err := configService.GetAllConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

// UpdateConfig 更新配置
func UpdateConfig(c *gin.Context) {
	configService := &services.ConfigService{}

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if req.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "配置键不能为空"})
		return
	}

	if err := configService.SetConfig(req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetConfig 获取单个配置
func GetConfig(c *gin.Context) {
	key := c.Param("key")

	configService := &services.ConfigService{}
	value, err := configService.GetConfig(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

// GetPublicConfig 获取公开配置（无需认证）
func GetPublicConfig(c *gin.Context) {
	configService := &services.ConfigService{}
	configs := make(map[string]string)

	// 只返回允许公开的配置
	if configService.IsRegistrationAllowed() {
		configs["allow_register"] = "1"
	} else {
		configs["allow_register"] = "0"
	}

	c.JSON(http.StatusOK, gin.H{"configs": configs})
}
