package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	// CACHE_DIR - for save local cache of articles in html
	CACHE_DIR = ".cache/pocket2dropbox"

	// CONFIG_DIR - dir in HOME for config file
	CONFIG_DIR = ".config"
	// CONFIG_PATH - config file
	CONFIG_PATH = CONFIG_DIR + "/pocket2dropbox.cfg"

	// LOCAL_INFO_FILENAME - filename for with info about downloaded articles
	LOCAL_INFO_FILENAME = "pocket2dropbox_info.json"
	// LOCAL_INFO_PATH - file with info about downloaded articles
	LOCAL_INFO_PATH = CACHE_DIR + "/" + LOCAL_INFO_FILENAME
)

// ----------------------------------------------------------------------------
func main() {
	cfg, err := get_config()
	if err != nil {
		log.Fatal(err)
	}

	local_articles, err := load_articles_info(cfg)
	if err != nil {
		log.Fatal(err)
	}

	articles, err := get_pocket_by_api(cfg)
	if err != nil {
		log.Fatal(err)
	}
	articles = merge_local_and_remote_info(local_articles, articles)
	hasChanges := false
	year_str := time.Now().Format("2006")
	local_html_path := fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), CACHE_DIR, year_str)
	create_dir_if_need(local_html_path)

	for i, item := range articles {

		if !item.IsDownloaded {
			file_name := "article." + time.Now().Format("2006-01-02_15-04-05")

			host := get_host(item.URL)
			if host != "" {
				file_name += "." + host
			}
			file_name += ".html"

			local_html_name := local_html_path + "/" + file_name

			log.Println("download:", item.URL)
			err = exec.Command("wgethtml.pl", "-j", "-a", local_html_name, item.URL).Run()
			if err == nil {
				if _, err := os.Stat(local_html_name); err == nil {
					item.FileName = file_name
					item.IsDownloaded = true
					hasChanges = true
					log.Println("downloaded:", item.URL)
				}
			} else {
				log.Println("download error:", err)
			}
		}

		if item.IsDownloaded && item.FileName != "" && !item.IsUploadedInDB {
			local_html_name := local_html_path + "/" + item.FileName
			err = upload_to_dropbox(local_html_name, "/"+year_str+"/"+item.FileName, cfg)
			if err != nil {
				log.Println("upload to dropbox error:", err)
			} else {
				item.IsUploadedInDB = true
				hasChanges = true
				log.Println("uploaded:", item.URL)
			}
		}
		articles[i] = item
	}

	if hasChanges {
		err = save_articles_info(articles, cfg)
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println("changes not found")
	}
}
