package main

import (
	"fmt"
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
		panic(err)
	}
	defer rows.Close()

	fp := gofeed.NewParser()
	prevReqs := make(map[string]*gofeed.Feed)
	var previousCode string
	var previousEmail string
	var feeds []*gofeed.Feed

	for rows.Next() {
		var code string
		var email string
		var link string
		err := rows.Scan(&code, &email, &link)
		if err != nil {
			panic(err)
		}

		// Send email when the next row has a different address
		if previousEmail != email && previousEmail != "" {
			fmt.Println("Email")
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
				fmt.Println(err)
			}
			prevReqs[link] = feed
			filtered := filterItemsByDate(feed)
			feeds = append(feeds, filtered)
		}
	}

	if feeds != nil {
		fmt.Println("Email")
		mail.SendSubEmail(previousEmail, feeds, previousCode)
	}
}

func filterItemsByDate(rss *gofeed.Feed) *gofeed.Feed {
	var filtered []*gofeed.Item
	for _, item := range rss.Items {
		pd, _ := dateparse.ParseAny(item.Published)
		duration := 7 * 24 * time.Hour // 7 days for now
		diff := time.Since(pd)

		if diff < duration {
			filtered = append(filtered, item)
		}
	}
	rss.Items = filtered
	return rss
}
