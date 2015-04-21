package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	// POKET_API_URL - Pocket API URL for get articles
	POKET_API_URL = "https://getpocket.com/v3/get?consumer_key=%s&access_token=%s&state=unread"

	// POKET_RSS_URL - Pocket RSS for get articles
	POKET_RSS_URL = "https://getpocket.com/users/%s/feed/unread"
)

// ----------------------------------------------------------------------------

// Article - one article from pocket
type Article struct {
	Title          string `json:"resolved_title" xml:"title"`
	URL            string `json:"resolved_url" xml:"link"`
	Date           string `json:"time_added" xml:"pubDate"`
	IsFavorite     bool   `json:"favorite"`
	Timestamp      int64  `json:"timestamp"`
	FileName       string `json:"filename"`          // local filename in cache
	IsDownloaded   bool   `json:"is_downloaded"`     // download to cache
	IsUploadedInDB bool   `json:"is_uploaded_in_db"` // uploaded to dropbox
}

// ArticleAPI - one article from pocket (for parse API out)
type ArticleAPI struct {
	Title    string `json:"resolved_title" xml:"title"`
	URL      string `json:"resolved_url" xml:"link"`
	Date     string `json:"time_added" xml:"pubDate"`
	Favorite string `json:"favorite"`
}

// Articles - list
type Articles []Article

// InfoJSON - save local info about all articles
type InfoJSON struct {
	Timestamp int64    `json:"timestamp"` // last exec time
	Items     Articles `json:"articles"`  // list of articles
}

// PocketRSS - RSS XML struct
type PocketRSS struct {
	XMLName xml.Name `xml:"rss"`
	Items   struct {
		XMLName  xml.Name `xml:"channel"`
		ItemList Articles `xml:"item"`
	} `xml:"channel"`
}

// PocketAPI - JSON struct
type PocketAPI struct {
	Since float32               `json:"since"`
	Items map[string]ArticleAPI `json:"list"`
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

	raw_data := PocketAPI{}
	if err := json.Unmarshal(jsonText, &raw_data); err != nil {
		return nil, err
	}

	result := Articles{}
	for _, item := range raw_data.Items {
		article := Article{
			Title:      item.Title,
			URL:        item.URL,
			IsFavorite: item.Favorite == "1",
		}
		timestamp, err := strconv.ParseInt(item.Date, 10, 64)
		if err == nil {
			article.Date = time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
			article.Timestamp = timestamp
		}

		result = append(result, article)
	}
	sort.Sort(result)

	return result, nil
}

// ----------------------------------------------------------------------------
func load_articles_info(cfg Config) (Articles, error) {
	info := InfoJSON{}
	local_info_path := os.Getenv("HOME") + "/" + LOCAL_INFO_PATH
	json_info, err := ioutil.ReadFile(local_info_path)
	if err == nil {
		if err := json.Unmarshal(json_info, &info); err != nil {
			return nil, err
		}
	}

	return info.Items, nil
}

// ----------------------------------------------------------------------------
func save_articles_info(articles Articles, cfg Config) error {
	info := InfoJSON{
		Timestamp: time.Now().Unix(),
		Items:     articles,
	}
	json_info, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}

	local_info_path := os.Getenv("HOME") + "/" + LOCAL_INFO_PATH
	err = ioutil.WriteFile(local_info_path, json_info, 0644)
	if err != nil {
		return err
	}

	err = upload_to_dropbox(local_info_path, LOCAL_INFO_FILENAME, cfg)
	if err != nil {
		return err
	}

	return nil
}

// ----------------------------------------------------------------------------
func merge_local_and_remote_info(local_articles, remote_articles Articles) (Articles, bool) {
	local_as_map := make(map[string]Article, len(local_articles))
	for _, item := range local_articles {
		local_as_map[item.URL] = item
	}

	result := make(Articles, 0, len(remote_articles))
	hasChanges := false

	for _, item := range remote_articles {
		if local_item, ok := local_as_map[item.URL]; ok {
			item.FileName = local_item.FileName
			item.IsDownloaded = local_item.IsDownloaded
			item.IsUploadedInDB = local_item.IsUploadedInDB
			if item.IsFavorite != local_item.IsFavorite {
				hasChanges = true
			}
		}
		result = append(result, item)
	}

	return result, hasChanges
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
