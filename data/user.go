package data

import (
	"database/sql"
	"math/rand"
	"time"
)

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
	code := generateRandomID(10)
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

func generateRandomID(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
