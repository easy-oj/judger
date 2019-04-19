package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/common/settings"
)

var (
	DB *sql.DB
)

func InitDatabase() {
	DB = dial(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", settings.MySQL.Username, settings.MySQL.Password,
		settings.MySQL.Host, settings.MySQL.Port, settings.MySQL.Database))
}

func dial(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	logs.Info("[Database] dial database on %s", strings.Split(dsn, "@")[1])
	return db
}
