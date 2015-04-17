package main

import (
	// "fmt"
	"log"

	"github.com/k0kubun/pp"
)

// ----------------------------------------------------------------------------
func main() {
	// articles, err := get_pocket_by_api()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// for i, article := range articles {
	// 	fmt.Printf("%d: %s %s (%s / %d)\n", i, article.Title, article.Link, article.Date, article.Timestamp)
	// }
	cfg, err := get_config()
	err = save_config(cfg)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(cfg)
}
