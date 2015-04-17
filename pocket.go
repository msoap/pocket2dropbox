package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	// "github.com/k0kubun/pp"
)

const (
	// Pocket API URL for get articles
	POKET_API_URL = "https://getpocket.com/v3/get?consumer_key=%s&access_token=%s&state=unread"

	// Pocket RSS for get articles
	POKET_RSS_URL = "https://getpocket.com/users/%s/feed/unread"
)

// ----------------------------------------------------------------------------

// Article - one article from pocket
type Article struct {
	Title     string `xml:"title" json:"resolved_title"`
	Link      string `xml:"link" json:"resolved_url"`
	Date      string `xml:"pubDate" json:"time_added"`
	Timestamp int64
}

// Articles - list
type Articles []Article

// PocketRSS - RSS XML struct
type PocketRSS struct {
	XMLName xml.Name `xml:"rss"`
	Items   struct {
		XMLName  xml.Name `xml:"channel"`
		ItemList Articles `xml:"item"`
	} `xml:"channel"`
}

// PocketJSON - JSON struct
type PocketJSON struct {
	Since float32            `json:"since"`
	Items map[string]Article `json:"list"`
}

// ----------------------------------------------------------------------------
func get_pocket_by_rss() (Articles, error) {
	rssText, err := http_get(
		fmt.Sprintf(POKET_RSS_URL, os.Getenv("POCKET_USER")),
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

	result := Articles{}
	for _, item := range rss.Items.ItemList {
		date, err := time.Parse(time.RFC1123Z, item.Date)
		if err == nil {
			item.Date = date.Format("2006-01-02 15:04:05")
			item.Timestamp = date.Unix()
		}
		result = append(result, item)
	}
	sort.Sort(result)

	return result, nil
}

// ----------------------------------------------------------------------------
func get_pocket_by_api(cfg Config) (Articles, error) {
	jsonText, err := http_get(
		fmt.Sprintf(POKET_API_URL, cfg.PocketKey, cfg.PocketToken),
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

	result := Articles{}
	for _, item := range raw_data.Items {
		timestamp, err := strconv.ParseInt(item.Date, 10, 64)
		if err == nil {
			item.Date = time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
			item.Timestamp = timestamp
		}

		result = append(result, item)
	}
	sort.Sort(result)

	return result, nil
}

// ----------------------------------------------------------------------------
// sorting
func (list Articles) Len() int {
	return len(list)
}

func (list Articles) Less(i, j int) bool {
	return list[i].Timestamp < list[j].Timestamp
}

func (list Articles) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
