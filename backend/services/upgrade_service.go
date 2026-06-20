package services

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"memo/logger"
)

// UpgradeStep 升级步骤
type UpgradeStep struct {
	Version string
	Migrate func(tx *sql.Tx) error
}

// UpgradeService 升级服务
type UpgradeService struct {
	db             *sql.DB
	currentVersion string
	dataDir        string
	steps          []UpgradeStep
}

// NewUpgradeService 创建升级服务
func NewUpgradeService(db *sql.DB, currentVersion string, dataDir string) *UpgradeService {
	return &UpgradeService{
		db:             db,
		currentVersion: currentVersion,
		dataDir:        dataDir,
		steps:          registry,
	}
}

// RunUpgrades 执行升级
func (s *UpgradeService) RunUpgrades() error {
	if err := s.createUpgradeTable(); err != nil {
		return err
	}

	applied, err := s.getAppliedVersions()
	if err != nil {
		return err
	}

	sugar := logger.GetSugarLogger()

	for _, step := range s.steps {
		if applied[step.Version] {
			continue
		}

		// 跳过比当前版本还新的步骤（未来的版本）
		if isNewerVersion(step.Version, s.currentVersion) {
			continue
		}

		sugar.Infof("正在执行升级 %s ...", step.Version)

		tx, err := s.db.Begin()
		if err != nil {
			s.recordUpgrade(step.Version, "", "失败：启动事务失败")
			return fmt.Errorf("升级 %s 启动事务失败: %w", step.Version, err)
		}

		if err := step.Migrate(tx); err != nil {
			tx.Rollback()
			s.recordUpgrade(step.Version, "", fmt.Sprintf("失败：%v", err))
			return fmt.Errorf("升级 %s 失败: %w", step.Version, err)
		}

		if err := tx.Commit(); err != nil {
			s.recordUpgrade(step.Version, "", fmt.Sprintf("失败：提交事务失败: %v", err))
			return fmt.Errorf("升级 %s 提交事务失败: %w", step.Version, err)
		}

		s.recordUpgrade(step.Version, s.currentVersion, "成功")
		sugar.Infof("升级 %s 完成", step.Version)
	}

	return nil
}

// parseVersion 解析版本号 v1.0.0 -> [1, 0, 0]
func parseVersion(v string) ([]int, error) {
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("版本号格式错误: %s", v)
	}
	nums := make([]int, 3)
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("版本号格式错误: %s", v)
		}
		nums[i] = n
	}
	return nums, nil
}

// isNewerVersion 检查 version 是否比 currentVersion 更新
func isNewerVersion(version, currentVersion string) bool {
	v, err := parseVersion(version)
	if err != nil {
		return false
	}
	c, err := parseVersion(currentVersion)
	if err != nil {
		return true
	}
	for i := 0; i < 3; i++ {
		if v[i] > c[i] {
			return true
		}
		if v[i] < c[i] {
			return false
		}
	}
	return false
}

// createUpgradeTable 创建升级记录表
func (s *UpgradeService) createUpgradeTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS sys_update (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version TEXT NOT NULL,
		content TEXT NOT NULL,
		upgrade_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		result TEXT NOT NULL,
		status INTEGER DEFAULT 0
	);
	`
	_, err := s.db.Exec(query)
	return err
}

// getAppliedVersions 获取已应用的升级版本
func (s *UpgradeService) getAppliedVersions() (map[string]bool, error) {
	rows, err := s.db.Query("SELECT DISTINCT version FROM sys_update WHERE status = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, nil
}

// recordUpgrade 记录升级结果
func (s *UpgradeService) recordUpgrade(version, content, result string) {
	sugar := logger.GetSugarLogger()
	_, err := s.db.Exec(
		"INSERT INTO sys_update (version, content, result, status) VALUES (?, ?, ?, ?)",
		version, content, result, 1,
	)
	if err != nil {
		sugar.Errorf("记录升级日志失败: %v", err)
	}
}

// 升级步骤注册表
var registry = []UpgradeStep{}
