package handlers

import (
	"memo/database"
	"memo/middleware"
	"memo/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username          string   `json:"username"`
	Password          string   `json:"password"`
	Nickname          string   `json:"nickname"`
	Email             string   `json:"email"`
	SecurityQuestions []string `json:"security_questions"`
	SecurityAnswers   []string `json:"security_answers"`
}

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	KeepLogin bool   `json:"keep_login"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// SetSecurityQuestion 设置三个安全问题 (需登录)
func SetSecurityQuestion(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		Questions []string `json:"questions"`
		Answers   []string `json:"answers"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	authService := &services.AuthService{}
	if err := authService.SetSecurityQuestions(userID, req.Questions, req.Answers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "安全问题设置成功"})
}

// GetSecurityQuestion 获取用户三个安全问题 (公开)
func GetSecurityQuestion(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	authService := &services.AuthService{}
	questions, err := authService.GetSecurityQuestions(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

// VerifySecurityAnswer 验证三个安全问题答案并重置密码 (公开)
func VerifySecurityAnswer(c *gin.Context) {
	var req struct {
		Username    string   `json:"username"`
		Answers     []string `json:"answers"`
		NewPassword string   `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if req.Username == "" || len(req.Answers) < 3 || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数不能为空"})
		return
	}

	if len(req.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码至少6位"})
		return
	}

	authService := &services.AuthService{}
	if err := authService.VerifySecurityAnswers(req.Username, req.Answers, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}

// AdminResetUserPassword 管理员重置用户密码 (需要管理员权限)
func AdminResetUserPassword(c *gin.Context) {
	targetUserID := middleware.GetUserIDFromParam(c)

	var req struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "新密码不能为空"})
		return
	}

	authService := &services.AuthService{}
	if err := authService.AdminResetUserPassword(targetUserID, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}

// CheckHasUsers 检查系统是否有用户
func CheckHasUsers(c *gin.Context) {
	authService := &services.AuthService{}
	hasUsers := authService.HasUsers()
	c.JSON(http.StatusOK, gin.H{
		"has_users":         hasUsers,
		"fresh_db":          database.FreshDB,
		"suspect_data_loss": database.SuspectDataLoss && !hasUsers,
	})
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
		return
	}

	authService := &services.AuthService{}
	user, err := authService.Register(req.Username, req.Password, req.Nickname, req.Email, req.SecurityQuestions, req.SecurityAnswers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": user.Token,
		"user":  user,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
		return
	}

	authService := &services.AuthService{}
	user, err := authService.Login(req.Username, req.Password, req.KeepLogin)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": user.Token,
		"user":  user,
	})
}

// CheckAuth 检查认证状态
func CheckAuth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}

	authService := &services.AuthService{}
	userID, err := authService.CheckToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user_id":       userID,
	})
}

// GetCurrentUser 获取当前用户
func GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)

	authService := &services.AuthService{}
	user, err := authService.GetCurrentUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
		return
	}

	authService := &services.AuthService{}
	if err := authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// CheckSecurityQuestion 检查当前用户是否设置了安全问题
func CheckSecurityQuestion(c *gin.Context) {
	userID := middleware.GetUserID(c)

	authService := &services.AuthService{}
	hasQuestion, err := authService.HasSecurityQuestion(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"has_security_question": hasQuestion})
}

// GetAllUsers 获取所有用户 (管理员)
func GetAllUsers(c *gin.Context) {
	authService := &services.AuthService{}
	users, err := authService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// CreateUser 创建用户 (管理员)
func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
		return
	}

	// 系统只允许一个管理员（首个注册用户），管理员后台创建的一律为普通用户
	authService := &services.AuthService{}
	user, err := authService.CreateUser(req.Username, req.Password, req.Nickname, req.Email, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "创建成功", "user": user})
}

// UpdateUser 更新用户 (管理员)
func UpdateUser(c *gin.Context) {
	currentUserID := middleware.GetUserID(c)
	targetUserID := middleware.GetUserIDFromParam(c)

	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	authService := &services.AuthService{}
	if err := authService.UpdateUser(currentUserID, targetUserID, req.Nickname, req.Email, req.Avatar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteUser 删除用户 (管理员)
func DeleteUser(c *gin.Context) {
	currentUserID := middleware.GetUserID(c)
	targetUserID := middleware.GetUserIDFromParam(c)

	authService := &services.AuthService{}
	if err := authService.DeleteUser(currentUserID, targetUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ToggleUserStatus 切换用户状态 (管理员)
func ToggleUserStatus(c *gin.Context) {
	targetUserID := middleware.GetUserIDFromParam(c)

	var req struct {
		IsActive int `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	configService := &services.ConfigService{}
	if err := configService.ToggleUserStatus(targetUserID, req.IsActive); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "操作成功"})
}
