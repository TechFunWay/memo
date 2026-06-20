package handlers

import (
	"memo/database"
	"memo/middleware"
	"memo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 北京时间（东八区）。容器镜像 FROM scratch 无 tzdata，使用固定时区偏移。
var beijingLoc = time.FixedZone("CST", 8*3600)

func nowBeijing() string {
	return time.Now().In(beijingLoc).Format("2006-01-02 15:04:05")
}

// GetMemos 获取备忘录列表
func GetMemos(c *gin.Context) {
	userID := middleware.GetUserID(c)
	search := c.Query("search")
	filter := c.Query("filter")
	folderID := c.Query("folder_id")

	query := "SELECT id, user_id, folder_id, title, content, content_html, is_pinned, is_favorite, tags, position, created_at, updated_at FROM memos WHERE user_id = ?"
	args := []any{userID}

	if search != "" {
		query += " AND (title LIKE ? OR content LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	if filter == "favorite" {
		query += " AND is_favorite = 1"
	}
	if filter == "pinned" {
		query += " AND is_pinned = 1"
	}
	if folderID != "" {
		query += " AND folder_id = ?"
		args = append(args, folderID)
	}

	query += " ORDER BY is_pinned DESC, position ASC, created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败: " + err.Error()})
		return
	}
	defer rows.Close()

	var memos []models.Memo
	for rows.Next() {
		var memo models.Memo
		var folderID sqlNullInt64
		err := rows.Scan(&memo.ID, &memo.UserID, &folderID, &memo.Title, &memo.Content, &memo.ContentHTML, &memo.IsPinned, &memo.IsFavorite, &memo.Tags, &memo.Position, &memo.CreatedAt, &memo.UpdatedAt)
		if err != nil {
			continue
		}
		if folderID.Valid {
			memo.FolderID = &folderID.Int64
		}
		memos = append(memos, memo)
	}

	c.JSON(http.StatusOK, gin.H{"memos": memos})
}

// CreateMemo 创建备忘录
func CreateMemo(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		FolderID *int64 `json:"folder_id"`
		Tags     string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	now := nowBeijing()

	result, err := database.DB.Exec(
		"INSERT INTO memos (user_id, folder_id, title, content, content_html, is_pinned, is_favorite, tags, position, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 0, 0, ?, 0, ?, ?)",
		userID, req.FolderID, req.Title, req.Content, "", req.Tags, now, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{
		"id":           id,
		"user_id":      userID,
		"folder_id":    req.FolderID,
		"title":        req.Title,
		"content":      req.Content,
		"content_html": "",
		"is_pinned":    false,
		"is_favorite":  false,
		"tags":         req.Tags,
		"position":     0,
		"created_at":   now,
		"updated_at":   now,
	})
}

// GetMemo 获取备忘录详情
func GetMemo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var memo models.Memo
	var folderID sqlNullInt64
	err := database.DB.QueryRow(
		"SELECT id, user_id, folder_id, title, content, content_html, is_pinned, is_favorite, tags, position, created_at, updated_at FROM memos WHERE id = ? AND user_id = ?",
		id, userID,
	).Scan(&memo.ID, &memo.UserID, &folderID, &memo.Title, &memo.Content, &memo.ContentHTML, &memo.IsPinned, &memo.IsFavorite, &memo.Tags, &memo.Position, &memo.CreatedAt, &memo.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "备忘录不存在"})
		return
	}

	if folderID.Valid {
		memo.FolderID = &folderID.Int64
	}

	c.JSON(http.StatusOK, memo)
}

// UpdateMemo 更新备忘录
func UpdateMemo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		FolderID   *int64 `json:"folder_id"`
		Tags       string `json:"tags"`
		IsPinned   *bool  `json:"is_pinned"`
		IsFavorite *bool  `json:"is_favorite"`
		Position   *int   `json:"position"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	var existingUserID int64
	err := database.DB.QueryRow("SELECT user_id FROM memos WHERE id = ?", id).Scan(&existingUserID)
	if err != nil || existingUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	now := nowBeijing()
	setClauses := []string{"updated_at = ?"}
	args := []any{now}

	if req.Title != "" {
		setClauses = append(setClauses, "title = ?")
		args = append(args, req.Title)
	}
	if req.Content != "" {
		setClauses = append(setClauses, "content = ?", "content_html = ?")
		args = append(args, req.Content, "")
	}
	if req.FolderID != nil {
		setClauses = append(setClauses, "folder_id = ?")
		args = append(args, req.FolderID)
	}
	if req.Tags != "" {
		setClauses = append(setClauses, "tags = ?")
		args = append(args, req.Tags)
	}
	if req.IsPinned != nil {
		setClauses = append(setClauses, "is_pinned = ?")
		args = append(args, req.IsPinned)
	}
	if req.IsFavorite != nil {
		setClauses = append(setClauses, "is_favorite = ?")
		args = append(args, req.IsFavorite)
	}
	if req.Position != nil {
		setClauses = append(setClauses, "position = ?")
		args = append(args, req.Position)
	}

	args = append(args, id)
	query := "UPDATE memos SET " + joinStrings(setClauses, ", ") + " WHERE id = ?"
	_, err = database.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteMemo 删除备忘录
func DeleteMemo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var existingUserID int64
	err := database.DB.QueryRow("SELECT user_id FROM memos WHERE id = ?", id).Scan(&existingUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "备忘录不存在"})
		return
	}
	if existingUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM memos WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// sqlNullInt64 wrapper for sqlite null int64
type sqlNullInt64 struct {
	Int64 int64
	Valid bool
}

func (n *sqlNullInt64) Scan(value any) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	switch v := value.(type) {
	case int64:
		n.Int64 = v
	case float64:
		n.Int64 = int64(v)
	default:
		n.Valid = false
	}
	return nil
}
