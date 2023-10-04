package data

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/brianbianchi/rss/util"
	_ "github.com/mattn/go-sqlite3"
)

type Feed struct {
	Link  string
	Blog  string
	Title string
}

func InitDb() *sql.DB {
	path := util.GetRootPath()
	db, err := sql.Open("sqlite3", fmt.Sprint(path, "rss.db"))
	if err != nil {
		panic(err)
	}
	return db
}

func GetSubs(db *sql.DB, code string) ([]string, error) {
	rows, err := db.Query("SELECT link FROM subs WHERE code=?", code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []string
	for rows.Next() {
		var link string
		err = rows.Scan(&link)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, link)
	}
	return feeds, nil
}

func GetSubsJoinUsers(db *sql.DB) (*sql.Rows, error) {
	query := `
		SELECT users.code, email, link 
		FROM subs JOIN users 
		ON subs.code = users.code 
		ORDER BY email ASC;
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func GetFeeds(db *sql.DB) ([]Feed, error) {
	rows, err := db.Query("SELECT link, blog, title FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []Feed
	for rows.Next() {
		var link string
		var blog string
		var title string
		err = rows.Scan(&link, &blog, &title)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, Feed{Link: link, Blog: blog, Title: title})
	}
	return feeds, nil
}

func CreateSubs(db *sql.DB, links []string, code string) error {
	for _, link := range links {
		_, err := db.Exec(`INSERT INTO subs (link, code) 
			VALUES (?, ?);`, link, code)

		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteSubs(db *sql.DB, code string) error {
	_, err := db.Exec("DELETE FROM subs WHERE code = ?", code)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByCode(db *sql.DB, code string) (string, error) {
	rows, err := db.Query("SELECT email FROM users WHERE code=?", code)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var email string
	for rows.Next() {
		err = rows.Scan(&email)
		if err != nil {
			return "", err
		}
	}
	return email, nil
}

func GetUserByEmail(db *sql.DB, email string) (string, error) {
	rows, err := db.Query("SELECT email FROM users WHERE email=?", email)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var code string
	for rows.Next() {
		err = rows.Scan(&code)
		if err != nil {
			return "", err
		}
	}
	return code, nil
}

func CreateUser(db *sql.DB, email string) (string, error) {
	code := generateRandomID()
	_, err := db.Exec(`INSERT INTO users (code, email, created) 
		VALUES (?, ?, ?);`, code, email, time.Now())
	if err != nil {
		return "", err
	}
	return code, nil
}

func DeleteUser(db *sql.DB, code string) error {
	_, err := db.Exec("DELETE FROM users WHERE code = ?", code)
	if err != nil {
		return err
	}
	return nil
}

func generateRandomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const idLength = 10

	b := make([]byte, idLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
