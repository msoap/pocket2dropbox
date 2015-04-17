package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	CACHE_DIR = ".cache/pocket2dropbox"
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

	for i, item := range articles {

		if !item.IsDownloaded {
			file_name := "article_" + time.Now().Format("2006-01-02_15-04-05.html")
			local_html_name := os.Getenv("HOME") + "/" + CACHE_DIR + "/2015/" + file_name

			err = exec.Command("wgethtml.pl", "-a", local_html_name, item.URL).Run()
			if err == nil {
				item.FileName = file_name
				item.IsDownloaded = true
				log.Println("downloaded:", item.URL)
			} else {
				log.Println("download error:", err)
			}
		}

		if item.IsDownloaded && item.FileName != "" && !item.IsUploadedInDB {
			local_html_name := os.Getenv("HOME") + "/" + CACHE_DIR + "/2015/" + item.FileName
			err = upload_to_dropbox(local_html_name, "/2015/"+item.FileName, cfg)
			if err != nil {
				log.Println("upload to dropbox error:", err)
			} else {
				item.IsUploadedInDB = true
				log.Println("uploaded:", item.URL)
			}
		}
		articles[i] = item

		if i == 2 {
			break
		}
	}

	err = save_articles_info(articles, cfg)
	if err != nil {
		log.Println(err)
	}
}
