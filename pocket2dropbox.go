package main

import (
	"fmt"
	"log"

	// "github.com/k0kubun/pp"
)

// ----------------------------------------------------------------------------
func main() {
	// articles, err := get_pocket_by_rss()
	articles, err := get_pocket_by_api()
	if err != nil {
		log.Fatal(err)
	}

	for i, article := range articles {
		fmt.Printf("%d: %s %s (%s)\n", i, article.Title, article.Link, article.Date)
	}
}
