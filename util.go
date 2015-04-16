package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/stacktic/dropbox"
)

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
func upload_to_dropbox(src, dst string) error {
	db := dropbox.NewDropbox()
	db.SetAppInfo(os.Getenv("DB_CLIENTID"), os.Getenv("DB_CLIENTSECRET"))
	db.SetAccessToken(os.Getenv("DB_TOKEN"))

	_, err := db.UploadFile(src, dst, true, "")
	return err
}
