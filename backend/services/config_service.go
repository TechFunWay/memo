package services

import (
	"errors"
	"memo/database"
	"time"
)

// ConfigService 系统配置服务
type ConfigService struct{}

// GetConfig 获取单个配置
func (s *ConfigService) GetConfig(key string) (string, error) {
	var value string
	err := database.DB.QueryRow("SELECT value FROM sys_config WHERE key = ? AND user_id = 0", key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

// SetConfig 设置配置
func (s *ConfigService) SetConfig(key, value string) error {
	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM sys_config WHERE key = ? AND user_id = 0", key).Scan(&exists)
	if err != nil {
		return err
	}

	now := time.Now()
	if exists > 0 {
		_, err = database.DB.Exec("UPDATE sys_config SET value = ?, updated_at = ? WHERE key = ? AND user_id = 0", value, now, key)
	} else {
		_, err = database.DB.Exec("INSERT INTO sys_config (user_id, key, value, created_at, updated_at) VALUES (0, ?, ?, ?, ?)", key, value, now, now)
	}
	return err
}

// GetAllConfigs 获取所有配置
func (s *ConfigService) GetAllConfigs() (map[string]string, error) {
	rows, err := database.DB.Query("SELECT key, value FROM sys_config WHERE user_id = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	configs := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		configs[key] = value
	}
	return configs, nil
}

// IsConfigTrue 判断配置是否为 "1"
func (s *ConfigService) IsConfigTrue(key string) bool {
	val, err := s.GetConfig(key)
	if err != nil {
		return false
	}
	return val == "1"
}

// IsRegistrationAllowed 检查是否允许注册
func (s *ConfigService) IsRegistrationAllowed() bool {
	val, err := s.GetConfig("allow_register")
	if err != nil {
		// 配置不存在时默认允许注册
		return true
	}
	return val == "1"
}

// ToggleUserStatus 切换用户状态
func (s *ConfigService) ToggleUserStatus(userID int64, isActive int) error {
	// 不能禁用管理员
	var isAdmin bool
	err := database.DB.QueryRow("SELECT is_admin FROM users WHERE id = ?", userID).Scan(&isAdmin)
	if err != nil {
		return errors.New("用户不存在")
	}
	if isAdmin {
		return errors.New("不能禁用管理员")
	}

	_, err = database.DB.Exec("UPDATE users SET is_active = ?, updated_at = ? WHERE id = ?", isActive, time.Now(), userID)
	return err
}
