package models

// User 用户模型
type User struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Nickname  string  `json:"nickname"`
	Email     string  `json:"email"`
	Avatar    string  `json:"avatar"`
	IsActive  bool    `json:"is_active"`
	IsAdmin   bool    `json:"is_admin"`
	Token     string  `json:"-"`
	APIKey    string  `json:"-"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
