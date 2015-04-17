package main

import (
	"fmt"
	"log"

	// "github.com/k0kubun/pp"
)

// ----------------------------------------------------------------------------
func main() {
	cfg, err := get_config()
	if err != nil {
		log.Fatal(err)
	}

	articles, err := get_pocket_by_api(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = save_articles_info(articles, cfg)
	if err != nil {
		log.Println(err)
	}

	for i, article := range articles {
		fmt.Printf("%d: %s %s (%s / %d)\n", i, article.Title, article.Link, article.Date, article.Timestamp)
	}

	err = save_config(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
