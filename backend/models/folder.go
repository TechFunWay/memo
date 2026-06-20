package models

// Folder 文件夹模型
type Folder struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	ParentID  *int64  `json:"parent_id"`
	Name      string  `json:"name"`
	Position  int     `json:"position"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
