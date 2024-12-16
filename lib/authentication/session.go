package authentication

import (
	"database/sql"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var sessionStore = sync.Map{}

func SessionSetup(dbDriver string, db *sql.DB, dsn, tableName string) {
	if dbDriver == "mysql" {
		sessionMySQLStart(db, tableName)
	} else {
		sessionStart(dsn, tableName)
	}
}

func sessionMySQLStart(db *sql.DB, tableName string) {
}

func sessionStart(dsn, tableName string) {
}

func AuthStore(c *fiber.Ctx, userID uint) {
	sessionStore.Store(c.IP(), userID)
}

func AuthGet(c *fiber.Ctx) (bool, uint) {
	userID, ok := sessionStore.Load(c.IP())

	if !ok {
		return false, 0
	}

	return true, userID.(uint)
}

func AuthDestroy(c *fiber.Ctx) {
	sessionStore.Delete(c.IP())
}
