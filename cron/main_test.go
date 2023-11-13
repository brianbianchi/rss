package main

import (
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

func TestFilterItemsByDate(t *testing.T) {
	rss := &gofeed.Feed{
		Items: []*gofeed.Item{
			{Published: time.Now().Add(-2 * 24 * time.Hour).String()},
			{Published: time.Now().Add(-4 * 24 * time.Hour).String()},
			{Published: time.Now().Add(-6 * 24 * time.Hour).String()},
			{Published: time.Now().Add(-8 * 24 * time.Hour).String()},
		},
	}
	result := len(filterItemsByDate(rss).Items)
	expected := 3
	if result != expected {
		t.Fatalf(`filterItemsByDate() returned %d, expected %d.`, result, expected)
	}

	rss = &gofeed.Feed{
		Items: []*gofeed.Item{
			{Published: time.Now().Add(-2 * 24 * time.Hour).String()},
			{Published: time.Now().Add(-8 * 24 * time.Hour).String()},
			{Published: time.Now().Add(-10 * 24 * time.Hour).String()},
			{Published: time.Now().Add(-12 * 24 * time.Hour).String()},
		},
	}
	result = len(filterItemsByDate(rss).Items)
	expected = 1
	if result != expected {
		t.Fatalf(`filterItemsByDate() returned %d, expected %d.`, result, expected)
	}

	rss = &gofeed.Feed{
		Items: []*gofeed.Item{
			{},
			{},
			{},
			{},
		},
	}
	result = len(filterItemsByDate(rss).Items)
	expected = 4
	if result != expected {
		t.Fatalf(`filterItemsByDate() returned %d, expected %d.`, result, expected)
	}
}
