package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/stacktic/dropbox"
)

// ----------------------------------------------------------------------------
func http_get(url, user, pass string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if user != "" && pass != "" {
		request.SetBasicAuth(user, pass)
	}

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
func upload_to_dropbox(src, dst string, cfg Config) error {
	db := dropbox.NewDropbox()
	db.SetAppInfo(cfg.DBClientId, cfg.DBClientSecret)
	db.SetAccessToken(cfg.DBToken)

	_, err := db.UploadFile(src, dst, true, "")
	return err
}

// ----------------------------------------------------------------------------
func get_host(URL string) string {
	url_object, err := url.Parse(URL)
	if err != nil {
		return ""
	}

	host := url_object.Host
	host = regexp.MustCompile(`:\d+$`).ReplaceAllString(host, "")

	return host
}

// ----------------------------------------------------------------------------
func create_dir_if_need(dir string) {
	if _, err := os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal("create dir error:", dir)
		}
	}
}
