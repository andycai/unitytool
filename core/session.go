package core

import (
	"database/sql"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"
	"github.com/gofiber/storage/mysql/v2"
	"gorm.io/gorm"
)

var (
	store *session.Store
)

const (
	sessionName = "unitool_session"
	userIDKey   = "user_id"
)

// FiberSession Fiber默认的会话表结构
type FiberSession struct {
	K string `gorm:"column:k;primaryKey"` // key
	V string `gorm:"column:v;not null"`   // value
	E int64  `gorm:"column:e;default:0"`  // expiry
}

// TableName 设置表名
func (FiberSession) TableName() string {
	return "sessions"
}

// SQLiteStorage 实现 fiber.Storage 接口
type SQLiteStorage struct {
	db *gorm.DB
}

// NewSQLiteStorage 创建新的SQLite存储
func NewSQLiteStorage(dsn string) (*SQLiteStorage, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 使用GORM的自动迁移来创建或更新表结构
	err = db.AutoMigrate(&FiberSession{})
	if err != nil {
		return nil, err
	}

	// 启动定期清理过期会话的goroutine
	storage := &SQLiteStorage{db: db}
	go storage.gcLoop()

	return storage, nil
}

func (s *SQLiteStorage) gcLoop() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		now := time.Now().Unix()
		s.db.Delete(&FiberSession{}, "e <= ? AND e != 0", now)
	}
}

// Get 获取会话数据
func (s *SQLiteStorage) Get(key string) ([]byte, error) {
	var session FiberSession
	if err := s.db.First(&session, "k = ? AND (e > ? OR e = 0)", key, time.Now().Unix()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return []byte(session.V), nil
}

// Set 设置会话数据
func (s *SQLiteStorage) Set(key string, val []byte, exp time.Duration) error {
	var expiry int64
	if exp != 0 {
		expiry = time.Now().Add(exp).Unix()
	}

	session := FiberSession{
		K: key,
		V: string(val),
		E: expiry,
	}

	return s.db.Save(&session).Error
}

// Delete 删除会话数据
func (s *SQLiteStorage) Delete(key string) error {
	return s.db.Delete(&FiberSession{}, "k = ?", key).Error
}

// Reset 重置存储
func (s *SQLiteStorage) Reset() error {
	return s.db.Where("1 = 1").Delete(&FiberSession{}).Error
}

// Close 关闭存储
func (s *SQLiteStorage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// SessionSetup 初始化session存储
func SessionSetup(dbDriver string, db *sql.DB, dsn, tableName string) {
	var storage fiber.Storage

	// 根据数据库类型选择存储方式
	switch dbDriver {
	case "mysql":
		if db != nil {
			storage = mysql.New(mysql.Config{
				Db:         db,
				Table:      tableName,
				Reset:      false,
				GCInterval: 10 * time.Second,
			})
		}
	case "sqlite":
		sqliteStorage, err := NewSQLiteStorage(dsn)
		if err == nil {
			storage = sqliteStorage
		}
	}

	// 如果没有成功创建数据库存储，则使用内存存储作为后备
	if storage == nil {
		storage = memory.New(memory.Config{
			GCInterval: 10 * time.Second,
		})
	}

	// 创建session存储
	store = session.New(session.Config{
		Storage:        storage,
		Expiration:     24 * time.Hour, // 默认session过期时间
		KeyLookup:      "cookie:" + sessionName,
		CookieSecure:   IsSecureMode(), // 根据安全模式配置决定是否启用secure
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		CookiePath:     "/", // 确保cookie在所有路径下可用
		CookieDomain:   "",  // 自动使用当前域名
	})
}

// StoreSession 存储用户认证信息
func StoreSession(c *fiber.Ctx, userID uint) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	sess.Set(userIDKey, userID)
	return sess.Save()
}

// GetSession 获取用户认证信息
func GetSession(c *fiber.Ctx) (bool, uint) {
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

// DestroySession 销毁用户认证信息
func DestroySession(c *fiber.Ctx) error {
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
