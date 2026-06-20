package models

// Memo 备忘录模型
type Memo struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"user_id"`
	FolderID    *int64  `json:"folder_id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	ContentHTML string  `json:"content_html"`
	IsPinned    bool    `json:"is_pinned"`
	IsFavorite  bool    `json:"is_favorite"`
	Tags        string  `json:"tags"`
	Position    int     `json:"position"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
