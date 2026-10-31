package fnos

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"memo/models"
	"memo/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrNotBound = errors.New("此飞牛 NAS 账号尚未绑定应用账号")

type Identity struct {
	UserID   int64
	Username string
	IsAdmin  bool
}

type Service struct{ DB *sql.DB }

func IdentityFromRequest(c *gin.Context) (Identity, bool) {
	if !IsGatewayRequest(c.Request.Context()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请从飞牛桌面中的应用入口使用 NAS 登录"})
		return Identity{}, false
	}
	uid, err := strconv.ParseInt(c.GetHeader("X-Trim-Userid"), 10, 64)
	username := c.GetHeader("X-Trim-Username")
	if err != nil || uid <= 0 || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到飞牛 NAS 登录信息"})
		return Identity{}, false
	}
	return Identity{UserID: uid, Username: username, IsAdmin: c.GetHeader("X-Trim-Isadmin") == "true"}, true
}

func scanUser(row interface{ Scan(...interface{}) error }) (*models.User, string, error) {
	var user models.User
	var password string
	var active, admin int
	err := row.Scan(&user.ID, &user.Username, &user.Nickname, &user.Email, &user.Avatar, &active, &admin, &user.Token, &user.APIKey, &user.CreatedAt, &user.UpdatedAt, &password)
	if err != nil {
		return nil, "", err
	}
	user.IsActive = active == 1
	user.IsAdmin = admin == 1
	return &user, password, nil
}

func (s Service) Login(identity Identity) (*models.User, error) {
	user, _, err := scanUser(s.DB.QueryRow(`
		SELECT u.id, u.username, u.nickname, u.email, u.avatar, u.is_active, u.is_admin,
		       u.token, u.api_key, u.created_at, u.updated_at, u.password
		FROM users u JOIN fnos_bindings f ON f.user_id = u.id
		WHERE f.fnos_user_id = ?`, identity.UserID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotBound
	}
	if err != nil {
		return nil, err
	}
	if !user.IsActive {
		return nil, fmt.Errorf("用户已被禁用")
	}
	token := uuid.New().String()
	if _, err := s.DB.Exec("UPDATE users SET token = ?, token_expires_at = ? WHERE id = ?", token, time.Now().Add(24*time.Hour).Format("2006-01-02 15:04:05"), user.ID); err != nil {
		return nil, err
	}
	if _, err := s.DB.Exec("UPDATE fnos_bindings SET fnos_username = ?, updated_at = CURRENT_TIMESTAMP WHERE fnos_user_id = ?", identity.Username, identity.UserID); err != nil {
		return nil, err
	}
	user.Token = token
	return user, nil
}

func (s Service) Bind(identity Identity, mode, username, password string) (*models.User, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var exists int
	if err := tx.QueryRow("SELECT COUNT(*) FROM fnos_bindings WHERE fnos_user_id = ?", identity.UserID).Scan(&exists); err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, fmt.Errorf("此飞牛 NAS 账号已绑定其他应用账号")
	}

	var user *models.User
	switch mode {
	case "register":
		var userCount int
		if err := tx.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
			return nil, err
		}
		if userCount > 0 {
			var allow string
			if err := tx.QueryRow("SELECT value FROM sys_config WHERE user_id = 0 AND key = 'allow_register'").Scan(&allow); err != nil || (allow != "1" && allow != "true") {
				return nil, fmt.Errorf("系统暂不允许注册")
			}
		}
		if err := tx.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&exists); err != nil {
			return nil, err
		}
		if exists > 0 {
			return nil, fmt.Errorf("用户名已存在")
		}
		token, apiKey := uuid.New().String(), uuid.New().String()
		result, err := tx.Exec(`INSERT INTO users (username, password, nickname, email, is_admin, token, api_key, is_active, created_at, updated_at)
			VALUES (?, ?, ?, '', ?, ?, ?, 1, ?, ?)`, username, services.MD5Hash(password+"memo"), username, userCount == 0, token, apiKey, time.Now(), time.Now())
		if err != nil {
			return nil, err
		}
		userID, _ := result.LastInsertId()
		user = &models.User{ID: userID, Username: username, Nickname: username, IsAdmin: userCount == 0, IsActive: true, Token: token, APIKey: apiKey}
	case "bind":
		var stored string
		user, stored, err = scanUser(tx.QueryRow(`SELECT id, username, nickname, email, avatar, is_active, is_admin, token, api_key, created_at, updated_at, password FROM users WHERE username = ?`, username))
		if err != nil || !user.IsActive || stored != services.MD5Hash(password+"memo") {
			return nil, fmt.Errorf("用户名或密码错误")
		}
	default:
		return nil, fmt.Errorf("无效的绑定方式")
	}

	if err := tx.QueryRow("SELECT COUNT(*) FROM fnos_bindings WHERE user_id = ?", user.ID).Scan(&exists); err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, fmt.Errorf("该应用账号已绑定其他飞牛 NAS 账号")
	}
	if _, err := tx.Exec("INSERT INTO fnos_bindings (user_id, fnos_user_id, fnos_username) VALUES (?, ?, ?)", user.ID, identity.UserID, identity.Username); err != nil {
		return nil, err
	}
	token := uuid.New().String()
	if _, err := tx.Exec("UPDATE users SET token = ?, token_expires_at = ? WHERE id = ?", token, time.Now().Add(24*time.Hour).Format("2006-01-02 15:04:05"), user.ID); err != nil {
		return nil, err
	}
	user.Token = token
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}

func RegisterRoutes(group *gin.RouterGroup, db *sql.DB) {
	service := Service{DB: db}
	group.POST("/login", func(c *gin.Context) {
		identity, ok := IdentityFromRequest(c)
		if !ok {
			return
		}
		user, err := service.Login(identity)
		if errors.Is(err, ErrNotBound) {
			c.JSON(http.StatusOK, gin.H{"binding_required": true, "fnos_username": identity.Username})
			return
		}
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "飞牛 NAS 登录失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": user.Token, "user": user})
	})
	group.POST("/bind", func(c *gin.Context) {
		identity, ok := IdentityFromRequest(c)
		if !ok {
			return
		}
		var req struct {
			Mode     string `json:"mode"`
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Mode == "" || req.Username == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请输入应用账号和密码"})
			return
		}
		user, err := service.Bind(identity, req.Mode, req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": user.Token, "user": user})
	})
}
