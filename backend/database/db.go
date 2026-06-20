package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB 初始化数据库
func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// 设置连接参数
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(2)

	// 测试连接
	if err := DB.Ping(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}
