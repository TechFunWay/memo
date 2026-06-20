package handlers

import (
	"memo/database"
	"memo/middleware"
	"memo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetFolders 获取文件夹列表
func GetFolders(c *gin.Context) {
	userID := middleware.GetUserID(c)

	rows, err := database.DB.Query(
		"SELECT id, user_id, parent_id, name, position, created_at, updated_at FROM folders WHERE user_id = ? ORDER BY position ASC, created_at ASC",
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败: " + err.Error()})
		return
	}
	defer rows.Close()

	var folders []models.Folder
	for rows.Next() {
		var folder models.Folder
		var parentID sqlNullInt64
		err := rows.Scan(&folder.ID, &folder.UserID, &parentID, &folder.Name, &folder.Position, &folder.CreatedAt, &folder.UpdatedAt)
		if err != nil {
			continue
		}
		if parentID.Valid {
			folder.ParentID = &parentID.Int64
		}
		folders = append(folders, folder)
	}

	c.JSON(http.StatusOK, gin.H{"folders": folders})
}

// CreateFolder 创建文件夹
func CreateFolder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		Name     string `json:"name"`
		ParentID *int64 `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	result, err := database.DB.Exec(
		"INSERT INTO folders (user_id, parent_id, name, position, created_at, updated_at) VALUES (?, ?, ?, 0, ?, ?)",
		userID, req.ParentID, req.Name, now, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{
		"id":         id,
		"user_id":    userID,
		"parent_id":  req.ParentID,
		"name":       req.Name,
		"position":   0,
		"created_at": now,
		"updated_at": now,
	})
}

// UpdateFolder 更新文件夹
func UpdateFolder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req struct {
		Name     string `json:"name"`
		ParentID *int64 `json:"parent_id"`
		Position *int   `json:"position"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	var existingUserID int64
	err := database.DB.QueryRow("SELECT user_id FROM folders WHERE id = ?", id).Scan(&existingUserID)
	if err != nil || existingUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	setClauses := []string{"updated_at = ?"}
	args := []any{now}

	if req.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, req.Name)
	}
	if req.ParentID != nil {
		setClauses = append(setClauses, "parent_id = ?")
		args = append(args, req.ParentID)
	}
	if req.Position != nil {
		setClauses = append(setClauses, "position = ?")
		args = append(args, req.Position)
	}

	args = append(args, id)
	query := "UPDATE folders SET " + joinStrings(setClauses, ", ") + " WHERE id = ?"
	_, err = database.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteFolder 删除文件夹
func DeleteFolder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var existingUserID int64
	err := database.DB.QueryRow("SELECT user_id FROM folders WHERE id = ?", id).Scan(&existingUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件夹不存在"})
		return
	}
	if existingUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	database.DB.Exec("UPDATE memos SET folder_id = NULL WHERE folder_id = ? AND user_id = ?", id, userID)
	_, err = database.DB.Exec("DELETE FROM folders WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
