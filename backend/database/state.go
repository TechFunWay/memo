package database

import (
	"log"
	"os"
	"path/filepath"
)

// 数据状态标记（启动时评估一次，经 /api/auth/has-users 暴露给前端）
var (
	// FreshDB 本次启动时数据库文件不存在、由应用新建（正常仅出现在首次安装）
	FreshDB bool
	// SuspectDataLoss 历史上曾有账号但当前数据库为空，疑似数据目录未挂载或被更换
	SuspectDataLoss bool
)

// sentinelName 哨兵文件名，放在可执行文件同目录而不是数据目录：
// 数据目录未挂载/被换成空目录时它仍然存在，用于识别"这里曾经有过数据"。
const sentinelName = ".had-users"

func sentinelPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exe), sentinelName), nil
}

// WriteHadUsersSentinel 落数哨兵，标记"系统曾有账号"。尽力而为，失败只记日志。
func WriteHadUsersSentinel() {
	p, err := sentinelPath()
	if err != nil {
		log.Printf("计算哨兵文件路径失败（不影响功能）: %v", err)
		return
	}
	if _, err := os.Stat(p); err == nil {
		return
	}
	if err := os.WriteFile(p, []byte("1"), 0644); err != nil {
		log.Printf("写入哨兵文件失败（不影响功能）: %v", err)
	}
}

// DetectDataState 在建表完成后调用，hasUsers 为当前用户表是否非空。
// 有用户时落哨兵；无用户但哨兵在，说明数据库被换成了空库，判定疑似数据丢失，
// 前端据此告警并阻止用户直接创建新管理员。
func DetectDataState(hasUsers bool) {
	SuspectDataLoss = false
	if hasUsers {
		WriteHadUsersSentinel()
		return
	}
	p, err := sentinelPath()
	if err != nil {
		return
	}
	if _, err := os.Stat(p); err == nil {
		SuspectDataLoss = true
		log.Println("⚠️  数据状态异常：历史上已有账号，但当前数据库为空，数据目录可能未挂载或被更换")
	}
}
