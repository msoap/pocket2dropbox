package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

// ----------------------------------------------------------------------------
type Article struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Date  string `xml:"pubDate"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Items   struct {
		XMLName  xml.Name  `xml:"channel"`
		ItemList []Article `xml:"item"`
	} `xml:"channel"`
}

// ----------------------------------------------------------------------------
func get_pocket_by_rss() []Article {
	rssText, err := http_auth_get(
		fmt.Sprintf("https://getpocket.com/users/%s/feed/unread", os.Getenv("POCKET_USER")),
		os.Getenv("POCKET_USER"),
		os.Getenv("POCKET_PASS"),
	)
	if err != nil {
		log.Fatal(err)
	}

	rss := RSS{}
	if xml.Unmarshal(rssText, &rss) != nil {
		log.Fatal(err)
	}

	result := []Article{}
	for _, item := range rss.Items.ItemList {
		result = append(result, item)
	}

	return result
}
