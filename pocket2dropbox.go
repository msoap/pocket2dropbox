package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// "github.com/k0kubun/pp"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Items   Items    `xml:"channel"`
}

type Items struct {
	XMLName  xml.Name `xml:"channel"`
	ItemList []Item   `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Date  string `xml:"pubDate"`
}

// ----------------------------------------------------------------------------
func http_auth_get(url, user, pass string) ([]byte, error) {

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(user, pass)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

// ----------------------------------------------------------------------------
func main() {
	rssText, err := http_auth_get("http://getpocket.com/users/msoap/feed/unread", os.Getenv("POCKET_USER"), os.Getenv("POCKET_PASS"))
	if err != nil {
		log.Fatal(err)
	}

	rss := RSS{}
	if xml.Unmarshal(rssText, &rss) != nil {
		log.Fatal(err)
	}

	for i, item := range rss.Items.ItemList {
		fmt.Printf("%d: %s %s (%s)\n", i, item.Title, item.Link, item.Date)
	}
}
