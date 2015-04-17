package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"

	// "github.com/k0kubun/pp"
)

const (
	CONFIG_FILE = ".config/pocket2dropbox.cfg"
)

// Config - configuration (tokens, secrets...)
type Config struct {
	// pocket settings
	PocketKey   string `json:"pocket_key"   env:"POCKET_KEY"`
	PocketToken string `json:"pocket_token" env:"POCKET_TOKEN"`

	// Dropbox settings
	DBClientId     string `json:"db_client_id"     env:"DB_CLIENTID"`
	DBClientSecret string `json:"db_client_secret" env:"DB_CLIENTSECRET"`
	DBToken        string `json:"db_token"         env:"DB_TOKEN"`
}

// ----------------------------------------------------------------------------
func get_config() (Config, error) {
	cfg := Config{}
	cfg_json, err := ioutil.ReadFile(os.Getenv("HOME") + "/" + CONFIG_FILE)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(cfg_json, &cfg); err != nil {
		return cfg, err
	}

	// If need get cfg from env
	reflect_val := reflect.ValueOf(&cfg).Elem()
	for i := 0; i < reflect_val.NumField(); i++ {
		valueField := reflect_val.Field(i)
		tag := reflect_val.Type().Field(i).Tag

		// if empty string in json
		env_value := os.Getenv(tag.Get("env"))
		if tag.Get("env") != "" && valueField.String() == "" && env_value != "" {
			valueField.SetString(env_value)
		}
	}

	return cfg, nil
}

// ----------------------------------------------------------------------------
func save_config(cfg Config) error {
	json_cfg, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(os.Getenv("HOME")+"/"+CONFIG_FILE, json_cfg, 0600)
	if err != nil {
		return err
	}

	return nil
}
