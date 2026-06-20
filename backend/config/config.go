package config

import (
	"os"
	"path/filepath"
)

// AppVersion 应用版本（单一定义源）
const AppVersion = "v1.0.0"

var (
	DataDir   string
	UploadDir string
	DBDir     string
)

// InitConfig 初始化配置
func InitConfig(dataDir string) {
	DataDir = dataDir
	UploadDir = filepath.Join(dataDir, "upload")
	DBDir = filepath.Join(dataDir, "db")

	// 确保目录存在
	os.MkdirAll(UploadDir, 0755)
	os.MkdirAll(DBDir, 0755)
}
