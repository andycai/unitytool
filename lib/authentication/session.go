package authentication

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"
	"github.com/gofiber/storage/mysql/v2"
)

var (
	store *session.Store
)

const (
	sessionName = "unitool_session"
	userIDKey   = "user_id"
)

// SessionSetup 初始化session存储
func SessionSetup(dbDriver string, db *sql.DB, dsn, tableName string) {
	var storage fiber.Storage

	// 默认使用内存存储，适用于所有平台
	storageConfig := memory.Config{
		GCInterval: 10 * time.Second,
	}
	storage = memory.New(storageConfig)

	// 仅在使用MySQL且提供了数据库连接时使用MySQL存储
	if dbDriver == "mysql" && db != nil {
		storage = mysql.New(mysql.Config{
			Db:         db,
			Table:      tableName,
			Reset:      false,
			GCInterval: 10 * time.Second,
		})
	}

	// 创建session存储
	store = session.New(session.Config{
		Storage:        storage,
		Expiration:     24 * time.Hour, // 默认session过期时间
		KeyLookup:      "cookie:" + sessionName,
		CookieName:     sessionName,
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})
}

// AuthStore 存储用户认证信息
func AuthStore(c *fiber.Ctx, userID uint) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	sess.Set(userIDKey, userID)
	return sess.Save()
}

// AuthGet 获取用户认证信息
func AuthGet(c *fiber.Ctx) (bool, uint) {
	sess, err := store.Get(c)
	if err != nil {
		return false, 0
	}

	userID := sess.Get(userIDKey)
	if userID == nil {
		return false, 0
	}

	// 转换为uint类型
	if id, ok := userID.(uint); ok {
		return true, id
	}

	// 处理可能存储为float64的情况
	if id, ok := userID.(float64); ok {
		return true, uint(id)
	}

	return false, 0
}

// AuthDestroy 销毁用户认证信息
func AuthDestroy(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	return sess.Destroy()
}

// SetSessionExpiration 设置session过期时间
func SetSessionExpiration(c *fiber.Ctx, duration time.Duration) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	sess.SetExpiry(duration)
	return sess.Save()
}
