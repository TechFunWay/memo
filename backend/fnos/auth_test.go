package fnos

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

func setupFnOSDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	for _, statement := range []string{
		`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE, password TEXT, nickname TEXT, email TEXT DEFAULT '', avatar TEXT DEFAULT '', is_active INTEGER DEFAULT 1, is_admin INTEGER DEFAULT 0, token TEXT, api_key TEXT, token_expires_at DATETIME, created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE sys_config (user_id INTEGER DEFAULT 0, key TEXT, value TEXT, PRIMARY KEY(user_id, key))`,
		`CREATE TABLE fnos_bindings (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER UNIQUE, fnos_user_id INTEGER UNIQUE, fnos_username TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME DEFAULT CURRENT_TIMESTAMP)`,
		`INSERT INTO sys_config (user_id, key, value) VALUES (0, 'allow_register', '1')`,
	} {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("schema: %v", err)
		}
	}
	return db
}

func TestBindThenLoginWithFnOS(t *testing.T) {
	service := Service{DB: setupFnOSDB(t)}
	user, err := service.Bind(Identity{UserID: 1000, Username: "nas-admin"}, "register", "admin", "client-md5")
	if err != nil || user.Token == "" || !user.IsAdmin {
		t.Fatalf("bind: user=%#v err=%v", user, err)
	}
	user, err = service.Login(Identity{UserID: 1000, Username: "renamed-admin"})
	if err != nil || user.Username != "admin" || user.Token == "" {
		t.Fatalf("login: user=%#v err=%v", user, err)
	}
	if _, err := service.Bind(Identity{UserID: 2000, Username: "other"}, "bind", "admin", "client-md5"); err == nil {
		t.Fatal("same application account was bound twice")
	}
}

func TestGatewayHeadersRequireTrustedConnection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/fnos/login", nil)
	req.Header.Set("X-Trim-Userid", "1000")
	req.Header.Set("X-Trim-Username", "nas-admin")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req
	if _, ok := IdentityFromRequest(c); ok {
		t.Fatal("trusted headers received over normal connection")
	}
	c.Request = req.WithContext(context.WithValue(req.Context(), gatewayContextKey{}, true))
	identity, ok := IdentityFromRequest(c)
	if !ok || identity.UserID != 1000 {
		t.Fatalf("identity=%#v ok=%v", identity, ok)
	}
}
