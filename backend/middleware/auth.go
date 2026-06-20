package middleware

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"memo/database"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证 token"})
			c.Abort()
			return
		}

		// 查询用户
		var userID int64
		var isActive bool
		var expiresAt sql.NullString
		err := database.DB.QueryRow("SELECT id, is_active, token_expires_at FROM users WHERE token = ?", token).Scan(&userID, &isActive, &expiresAt)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证 token"})
			c.Abort()
			return
		}

		if !isActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "账号已被禁用"})
			c.Abort()
			return
		}

		// 检查 token 是否过期
		if expiresAt.Valid && expiresAt.String != "" {
			expTime, err := time.Parse("2006-01-02 15:04:05", expiresAt.String)
			if err == nil && time.Now().After(expTime) {
				database.DB.Exec("UPDATE users SET token = '', token_expires_at = NULL WHERE id = ?", userID)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "登录已过期，请重新登录"})
				c.Abort()
				return
			}
		}

		// 将用户 ID 存入上下文
		c.Set("user_id", userID)
		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 检查是否为管理员
		var isAdmin bool
		err := database.DB.QueryRow("SELECT is_admin FROM users WHERE id = ?", userID).Scan(&isAdmin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
			c.Abort()
			return
		}

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetUserIDFromParam 从 URL 参数获取用户 ID
func GetUserIDFromParam(c *gin.Context) int64 {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0
	}
	return userID
}
