package services

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"memo/database"
	"memo/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AuthService 认证服务
type AuthService struct{}

// MD5Hash 计算 MD5 哈希
func MD5Hash(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// hashAnswer 计算安全答案哈希 (空答案返回空串)；答案不区分大小写与首尾空格
func hashAnswer(answer string) string {
	a := strings.TrimSpace(strings.ToLower(answer))
	if a == "" {
		return ""
	}
	return MD5Hash(a + "memo")
}

// normSlice3 将切片规整为长度 3 (不足补空串)
func normSlice3(s []string) [3]string {
	var out [3]string
	for i := 0; i < 3 && i < len(s); i++ {
		out[i] = strings.TrimSpace(s[i])
	}
	return out
}

// Register 用户注册
func (s *AuthService) Register(username, password, nickname, email string, questions, answers []string) (*models.User, error) {
	// 检查是否是第一个用户 (设为管理员)
	var userCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	isAdmin := userCount == 0

	// 首个用户(管理员)不受"允许注册"开关限制，其余用户需要开关开启
	if !isAdmin {
		configService := &ConfigService{}
		if !configService.IsRegistrationAllowed() {
			return nil, errors.New("系统暂不允许注册")
		}
	}

	// 检查用户名是否存在
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 密码加密 (前端已做 MD5，后端加盐)
	secondHash := MD5Hash(password + "memo")

	// 生成 token 和 api_key
	token := uuid.New().String()
	apiKey := uuid.New().String()

	// 安全问题 (可选，用于找回密码，最多三个)
	q := normSlice3(questions)
	a := normSlice3(answers)

	// 插入用户
	result, err := database.DB.Exec(`
		INSERT INTO users (username, password, nickname, email, is_admin, token, api_key,
			security_question, security_answer, security_question2, security_answer2, security_question3, security_answer3,
			is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?)
	`, username, secondHash, nickname, email, isAdmin, token, apiKey,
		q[0], hashAnswer(a[0]), q[1], hashAnswer(a[1]), q[2], hashAnswer(a[2]),
		time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	userID, _ := result.LastInsertId()

	return &models.User{
		ID:       userID,
		Username: username,
		Nickname: nickname,
		Email:    email,
		IsAdmin:  isAdmin,
		Token:    token,
		IsActive: true,
	}, nil
}

// CreateUser 管理员创建用户
func (s *AuthService) CreateUser(username, password, nickname, email string, isAdmin bool) (*models.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("用户名和密码不能为空")
	}

	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	secondHash := MD5Hash(password + "memo")
	token := uuid.New().String()
	apiKey := uuid.New().String()

	if nickname == "" {
		nickname = username
	}

	result, err := database.DB.Exec(`
		INSERT INTO users (username, password, nickname, email, is_admin, token, api_key, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, 1, ?, ?)
	`, username, secondHash, nickname, email, isAdmin, token, apiKey, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	userID, _ := result.LastInsertId()
	return &models.User{
		ID:       userID,
		Username: username,
		Nickname: nickname,
		Email:    email,
		IsAdmin:  isAdmin,
		IsActive: true,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(username, password string, keepLogin ...bool) (*models.User, error) {
	keep := len(keepLogin) > 0 && keepLogin[0]

	// 查询用户
	var user models.User
	var dbPassword string
	var isActive int
	var isAdmin int

	err := database.DB.QueryRow(`
		SELECT id, username, nickname, email, avatar, is_active, is_admin, token, api_key, created_at, updated_at, password
		FROM users WHERE username = ?
	`, username).Scan(
		&user.ID, &user.Username, &user.Nickname, &user.Email, &user.Avatar,
		&isActive, &isAdmin, &user.Token, &user.APIKey, &user.CreatedAt, &user.UpdatedAt, &dbPassword,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	if isActive == 0 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码 (前端已做 MD5，后端加盐)
	secondHash := MD5Hash(password + "memo")

	if dbPassword != secondHash {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成新 token
	token := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 默认7天
	if !keep {
		expiresAt = time.Now().Add(24 * time.Hour) // 普通登录24小时
	}
	_, err = database.DB.Exec(
		"UPDATE users SET token = ?, token_expires_at = ? WHERE id = ?",
		token, expiresAt.Format("2006-01-02 15:04:05"), user.ID,
	)
	if err != nil {
		return nil, err
	}
	user.Token = token

	user.IsActive = isActive == 1
	user.IsAdmin = isAdmin == 1

	return &user, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	// 查询当前密码
	var dbPassword string
	err := database.DB.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&dbPassword)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码 (前端已做 MD5，后端加盐)
	secondHash := MD5Hash(oldPassword + "memo")

	if dbPassword != secondHash {
		return errors.New("旧密码错误")
	}

	// 更新密码 (前端已做 MD5，后端加盐)
	newSecondHash := MD5Hash(newPassword + "memo")

	_, err = database.DB.Exec("UPDATE users SET password = ?, updated_at = ? WHERE id = ?", newSecondHash, time.Now(), userID)
	return err
}

// GetCurrentUser 获取当前用户
func (s *AuthService) GetCurrentUser(userID int64) (*models.User, error) {
	var user models.User
	var isActive int
	var isAdmin int

	err := database.DB.QueryRow(`
		SELECT id, username, nickname, email, avatar, is_active, is_admin, token, api_key, created_at, updated_at
		FROM users WHERE id = ?
	`, userID).Scan(
		&user.ID, &user.Username, &user.Nickname, &user.Email, &user.Avatar,
		&isActive, &isAdmin, &user.Token, &user.APIKey, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.IsActive = isActive == 1
	user.IsAdmin = isAdmin == 1

	return &user, nil
}

// GetAllUsers 获取所有用户 (管理员)
func (s *AuthService) GetAllUsers() ([]models.User, error) {
	rows, err := database.DB.Query(`
		SELECT id, username, nickname, email, avatar, is_active, is_admin, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var isActive int
		var isAdmin int
		err := rows.Scan(
			&user.ID, &user.Username, &user.Nickname, &user.Email, &user.Avatar,
			&isActive, &isAdmin, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		user.IsActive = isActive == 1
		user.IsAdmin = isAdmin == 1
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser 更新用户
func (s *AuthService) UpdateUser(userID, targetUserID int64, nickname, email, avatar string) error {
	_, err := database.DB.Exec(`
		UPDATE users SET nickname = ?, email = ?, avatar = ?, updated_at = ?
		WHERE id = ?
	`, nickname, email, avatar, time.Now(), targetUserID)
	return err
}

// DeleteUser 删除用户
func (s *AuthService) DeleteUser(currentUserID, targetUserID int64) error {
	// 不能删除自己
	if currentUserID == targetUserID {
		return errors.New("不能删除自己")
	}

	// 检查目标用户是否为管理员
	var isAdmin bool
	err := database.DB.QueryRow("SELECT is_admin FROM users WHERE id = ?", targetUserID).Scan(&isAdmin)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 不能删除管理员
	if isAdmin {
		return errors.New("不能删除管理员")
	}

	_, err = database.DB.Exec("DELETE FROM users WHERE id = ?", targetUserID)
	return err
}

// AdminResetUserPassword 管理员重置指定用户的密码
func (s *AuthService) AdminResetUserPassword(userID int64, newPassword string) error {
	// 检查用户是否存在
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&exists)
	if err != nil || !exists {
		return errors.New("用户不存在")
	}

	// 更新密码 (前端已做 MD5，后端加盐)
	newHash := MD5Hash(newPassword + "memo")
	_, err = database.DB.Exec("UPDATE users SET password = ?, updated_at = ? WHERE id = ?", newHash, time.Now(), userID)
	return err
}

// AdminResetPassword 管理员直接重置密码 (通过命令行)
func (s *AuthService) AdminResetPassword(newPassword string) error {
	// 查找管理员
	var userID int64
	err := database.DB.QueryRow("SELECT id FROM users WHERE is_admin = 1 LIMIT 1").Scan(&userID)
	if err != nil {
		return errors.New("未找到管理员用户，请先注册")
	}

	newHash := MD5Hash(newPassword + "memo")
	_, err = database.DB.Exec("UPDATE users SET password = ?, updated_at = ? WHERE id = ?", newHash, time.Now(), userID)
	return err
}

// HasUsers 检查系统中是否有用户
func (s *AuthService) HasUsers() bool {
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count > 0
}

// SetSecurityQuestions 设置三个安全问题 (必须全部填写且问题互不相同)
func (s *AuthService) SetSecurityQuestions(userID int64, questions, answers []string) error {
	q := normSlice3(questions)
	a := normSlice3(answers)

	for i := 0; i < 3; i++ {
		if q[i] == "" || a[i] == "" {
			return errors.New("请完整设置三个安全问题及答案")
		}
	}
	if q[0] == q[1] || q[0] == q[2] || q[1] == q[2] {
		return errors.New("三个安全问题不能重复")
	}

	_, err := database.DB.Exec(
		`UPDATE users SET security_question = ?, security_answer = ?,
			security_question2 = ?, security_answer2 = ?,
			security_question3 = ?, security_answer3 = ?, updated_at = ? WHERE id = ?`,
		q[0], hashAnswer(a[0]), q[1], hashAnswer(a[1]), q[2], hashAnswer(a[2]), time.Now(), userID,
	)
	return err
}

// GetSecurityQuestions 获取用户的三个安全问题
func (s *AuthService) GetSecurityQuestions(username string) ([]string, error) {
	var q1, q2, q3 string
	err := database.DB.QueryRow(
		`SELECT security_question, security_question2, security_question3
		 FROM users WHERE username = ?`,
		username,
	).Scan(&q1, &q2, &q3)
	if err != nil {
		return nil, errors.New("用户未设置安全问题")
	}
	if q1 == "" || q2 == "" || q3 == "" {
		return nil, errors.New("用户未设置安全问题")
	}
	return []string{q1, q2, q3}, nil
}

// VerifySecurityAnswers 验证三个安全问题答案并重置密码 (必须全部正确)
func (s *AuthService) VerifySecurityAnswers(username string, answers []string, newPassword string) error {
	var a1, a2, a3 string
	err := database.DB.QueryRow(
		`SELECT security_answer, security_answer2, security_answer3
		 FROM users WHERE username = ?`,
		username,
	).Scan(&a1, &a2, &a3)
	if err != nil {
		return errors.New("用户未设置安全问题")
	}
	if a1 == "" || a2 == "" || a3 == "" {
		return errors.New("用户未设置安全问题")
	}

	ans := normSlice3(answers)
	if a1 != hashAnswer(ans[0]) || a2 != hashAnswer(ans[1]) || a3 != hashAnswer(ans[2]) {
		return errors.New("安全问题答案错误")
	}

	newHash := MD5Hash(newPassword + "memo")
	_, err = database.DB.Exec(
		"UPDATE users SET password = ?, updated_at = ? WHERE username = ?",
		newHash, time.Now(), username,
	)
	return err
}

// HasSecurityQuestion 检查用户是否设置了全部三个安全问题
func (s *AuthService) HasSecurityQuestion(userID int64) (bool, error) {
	var count int
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM users WHERE id = ?
		 AND security_question != '' AND security_answer != ''
		 AND security_question2 != '' AND security_answer2 != ''
		 AND security_question3 != '' AND security_answer3 != ''`,
		userID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckToken 检查 token 是否有效
func (s *AuthService) CheckToken(token string) (int64, error) {
	if token == "" {
		return 0, errors.New("未提供认证 token")
	}

	var userID int64
	var isActive bool
	err := database.DB.QueryRow("SELECT id, is_active FROM users WHERE token = ?", token).Scan(&userID, &isActive)
	if err != nil {
		return 0, errors.New("无效的认证 token")
	}

	if !isActive {
		return 0, errors.New("账号已被禁用")
	}

	return userID, nil
}
