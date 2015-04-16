package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"time"

	// "github.com/k0kubun/pp"
)

// ----------------------------------------------------------------------------
type Article struct {
	Title string `xml:"title" json:"resolved_title"`
	Link  string `xml:"link" json:"resolved_url"`
	Date  string `xml:"pubDate" json:"time_added"`
}

type PocketRSS struct {
	XMLName xml.Name `xml:"rss"`
	Items   struct {
		XMLName  xml.Name  `xml:"channel"`
		ItemList []Article `xml:"item"`
	} `xml:"channel"`
}

type PocketJSON struct {
	Since float32            `json:"since"`
	Items map[string]Article `json:"list"`
}

// ----------------------------------------------------------------------------
func get_pocket_by_rss() ([]Article, error) {
	rssText, err := http_get(
		fmt.Sprintf("https://getpocket.com/users/%s/feed/unread", os.Getenv("POCKET_USER")),
		os.Getenv("POCKET_USER"),
		os.Getenv("POCKET_PASS"),
	)
	if err != nil {
		return nil, err
	}

	rss := PocketRSS{}
	if xml.Unmarshal(rssText, &rss) != nil {
		return nil, err
	}

	result := []Article{}
	for _, item := range rss.Items.ItemList {
		date, err := time.Parse(time.RFC1123Z, item.Date)
		if err == nil {
			item.Date = date.Format("2006-01-02 15:04:05")
		}
		result = append(result, item)
	}

	return result, nil
}

// ----------------------------------------------------------------------------
func get_pocket_by_api() ([]Article, error) {
	jsonText, err := http_get(
		fmt.Sprintf("https://getpocket.com/v3/get?consumer_key=%s&access_token=%s&state=unread", os.Getenv("POCKET_KEY"), os.Getenv("POCKET_TOKEN")),
		"",
		"",
	)
	if err != nil {
		return nil, err
	}

	raw_data := PocketJSON{}
	if err := json.Unmarshal(jsonText, &raw_data); err != nil {
		return nil, err
	}

	result := []Article{}
	for _, item := range raw_data.Items {
		timestamp, err := strconv.ParseInt(item.Date, 10, 64)
		if err == nil {
			item.Date = time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
		}

		result = append(result, item)
	}

	return result, nil
}
