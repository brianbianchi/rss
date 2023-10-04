package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/brianbianchi/rss/data"
	"github.com/brianbianchi/rss/util"
)

func main() {
	db := data.InitDb()
	defer db.Close()

	CreateTables(db)
	CreateFeeds(db)
}

func CreateTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		code TEXT PRIMARY KEY NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS feeds  (
		link TEXT PRIMARY KEY NOT NULL,
		blog TEXT NOT NULL,
		title TEXT
	);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS subs (
		link TEXT NOT NULL,
		code TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}
	return nil
}

func CreateFeeds(db *sql.DB) {
	path := util.GetRootPath()
	file, err := os.Open(fmt.Sprint(path, "db/feeds.csv"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	sqlTemplate := "INSERT INTO feeds (link, blog, title) VALUES (?, ?, ?);"
	for {
		record, err := reader.Read()
		if err != nil {
			panic(err)
		}

		link := record[0]
		blog := record[1]
		title := record[2]

		stmt, err := db.Prepare(sqlTemplate)
		if err != nil {
			fmt.Println("Error preparing SQL statement:", err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(link, blog, title)
		if err != nil {
			fmt.Println("Error executing SQL statement:", err)
			return
		}
		fmt.Println("Record inserted successfully.")
	}
}
