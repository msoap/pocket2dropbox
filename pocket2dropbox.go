package main

import (
	"fmt"

	// "github.com/k0kubun/pp"
)

// ----------------------------------------------------------------------------
func main() {
	articles := get_pocket_by_rss()
	for i, article := range articles {
		fmt.Printf("%d: %s %s (%s)\n", i, article.Title, article.Link, article.Date)
	}
}
