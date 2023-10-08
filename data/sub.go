package data

import (
	"database/sql"
)

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
