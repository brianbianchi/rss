package data

import (
	"database/sql"
)

type Feed struct {
	Link  string
	Blog  string
	Title string
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
