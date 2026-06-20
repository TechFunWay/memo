package database

import (
	"log"
	"strings"
)

// InitTables 初始化数据库表
func InitTables() error {
	// 用户表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		nickname TEXT,
		email TEXT,
		avatar TEXT DEFAULT '',
		is_active INTEGER DEFAULT 1,
		is_admin INTEGER DEFAULT 0,
		token TEXT,
		api_key TEXT,
		security_question TEXT DEFAULT '',
		security_answer TEXT DEFAULT '',
		security_question2 TEXT DEFAULT '',
		security_answer2 TEXT DEFAULT '',
		security_question3 TEXT DEFAULT '',
		security_answer3 TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE INDEX IF NOT EXISTS idx_users_token ON users(token);
	`

	// 备忘录表
	createMemosTable := `
	CREATE TABLE IF NOT EXISTS memos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		folder_id INTEGER,
		title TEXT DEFAULT '',
		content TEXT,
		content_html TEXT,
		is_pinned INTEGER DEFAULT 0,
		is_favorite INTEGER DEFAULT 0,
		tags TEXT DEFAULT '',
		position INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE SET NULL
	);
	CREATE INDEX IF NOT EXISTS idx_memos_user_id ON memos(user_id);
	CREATE INDEX IF NOT EXISTS idx_memos_folder_id ON memos(folder_id);
	`

	// 文件夹表
	createFoldersTable := `
	CREATE TABLE IF NOT EXISTS folders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		parent_id INTEGER,
		name TEXT NOT NULL,
		position INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (parent_id) REFERENCES folders(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_folders_user_id ON folders(user_id);
	CREATE INDEX IF NOT EXISTS idx_folders_parent_id ON folders(parent_id);
	`

	// 附件表
	createAttachmentsTable := `
	CREATE TABLE IF NOT EXISTS memo_attachments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		memo_id INTEGER NOT NULL,
		file_name TEXT NOT NULL,
		file_path TEXT NOT NULL,
		file_size INTEGER DEFAULT 0,
		mime_type TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (memo_id) REFERENCES memos(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_attachments_memo_id ON memo_attachments(memo_id);
	`

	// 标签表
	createTagsTable := `
	CREATE TABLE IF NOT EXISTS memo_tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		color TEXT DEFAULT '#18a058',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, name),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`

	// 备忘录 - 标签关联表
	createMemoTagsTable := `
	CREATE TABLE IF NOT EXISTS memo_memo_tags (
		memo_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (memo_id, tag_id),
		FOREIGN KEY (memo_id) REFERENCES memos(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES memo_tags(id) ON DELETE CASCADE
	);
	`

	// 系统配置表
	createSysConfigTable := `
	CREATE TABLE IF NOT EXISTS sys_config (
		user_id INTEGER DEFAULT 0,
		key TEXT NOT NULL,
		value TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, key)
	);
	`

	// 执行所有表创建语句
	tables := []string{
		createUsersTable,
		createMemosTable,
		createFoldersTable,
		createAttachmentsTable,
		createTagsTable,
		createMemoTagsTable,
		createSysConfigTable,
	}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			log.Printf("Error creating table: %v", err)
			return err
		}
	}

	// 兼容迁移：为旧数据库补充可能缺失的列
	migrateColumns := []string{
		"ALTER TABLE users ADD COLUMN security_question TEXT DEFAULT ''",
		"ALTER TABLE users ADD COLUMN security_answer TEXT DEFAULT ''",
		"ALTER TABLE users ADD COLUMN security_question2 TEXT DEFAULT ''",
		"ALTER TABLE users ADD COLUMN security_answer2 TEXT DEFAULT ''",
		"ALTER TABLE users ADD COLUMN security_question3 TEXT DEFAULT ''",
		"ALTER TABLE users ADD COLUMN security_answer3 TEXT DEFAULT ''",
		"ALTER TABLE users ADD COLUMN token_expires_at DATETIME",
	}
	for _, m := range migrateColumns {
		if _, err := DB.Exec(m); err != nil {
			if !strings.Contains(err.Error(), "duplicate column") {
				log.Printf("Migration warning: %v", err)
			}
		}
	}

	log.Println("All tables created successfully")
	return nil
}
