package data

import (
	"database/sql"
	"fmt"

	"github.com/brianbianchi/rss/util"
	_ "github.com/mattn/go-sqlite3"
)

func InitDb() *sql.DB {
	path := util.GetRootPath()
	db, err := sql.Open("sqlite3", fmt.Sprint(path, "rss.db"))
	if err != nil {
		panic(err)
	}
	return db
}
