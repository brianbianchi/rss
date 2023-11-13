package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/araddon/dateparse"
	"github.com/brianbianchi/rss/data"
	"github.com/brianbianchi/rss/mail"
	"github.com/mmcdole/gofeed"
)

func main() {
	db := data.InitDb()
	defer db.Close()

	rows, err := data.GetSubsJoinUsers(db)
	if err != nil {
		log.Fatalf("Error fetching subscriptions: %v", err)
	}
	defer rows.Close()

	processSubscriptions(rows)
}

func processSubscriptions(rows *sql.Rows) {
	fp := gofeed.NewParser()
	prevReqs := make(map[string]*gofeed.Feed)
	var previousCode string
	var previousEmail string
	var feeds []*gofeed.Feed

	for rows.Next() {
		var code, email, link string
		if err := rows.Scan(&code, &email, &link); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Send email when the next row has a different address
		if previousEmail != email && previousEmail != "" {
			mail.SendSubEmail(previousEmail, feeds, previousCode)
			feeds = nil
		}
		previousCode = code
		previousEmail = email

		if val, ok := prevReqs[link]; ok {
			feeds = append(feeds, val)
		} else {
			feed, err := fp.ParseURL(link)
			if err != nil {
				log.Printf("Error parsing feed URL %s: %v", link, err)
				continue
			}
			prevReqs[link] = feed
			filtered := filterItemsByDate(feed)
			feeds = append(feeds, filtered)
		}
	}

	if feeds != nil {
		mail.SendSubEmail(previousEmail, feeds, previousCode)
	}
}

func filterItemsByDate(rss *gofeed.Feed) *gofeed.Feed {
	var filteredItems []*gofeed.Item
	for _, item := range rss.Items {
		if item.Published == "" {
			filteredItems = append(filteredItems, item)
			continue
		}
		publishedDate, _ := dateparse.ParseAny(item.Published)
		sevenDays := 7 * 24 * time.Hour
		sevenDaysAgo := time.Now().Add(-sevenDays)

		if publishedDate.After(sevenDaysAgo) {
			filteredItems = append(filteredItems, item)
		}
	}
	// include first 10
	if len(filteredItems) > 10 {
		filteredItems = filteredItems[:10]
	}
	rss.Items = filteredItems
	return rss
}
