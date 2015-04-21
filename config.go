package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

// Config - configuration (tokens, secrets...)
type Config struct {
	// pocket settings
	PocketKey   string `json:"pocket_key"   env:"POCKET_KEY"`
	PocketToken string `json:"pocket_token" env:"POCKET_TOKEN"`

	// Dropbox settings
	DBClientID     string `json:"db_client_id"     env:"DB_CLIENTID"`
	DBClientSecret string `json:"db_client_secret" env:"DB_CLIENTSECRET"`
	DBToken        string `json:"db_token"         env:"DB_TOKEN"`

	// save favorites articles only
	Favorites bool `json:"favorites"`
}

// ----------------------------------------------------------------------------
func get_config() (Config, error) {
	cfg := Config{}
	cfg_json, err := ioutil.ReadFile(os.Getenv("HOME") + "/" + CONFIG_PATH)
	if err == nil {
		if err := json.Unmarshal(cfg_json, &cfg); err != nil {
			return cfg, err
		}
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

	// command line options
	favorites := flag.Bool("favorites", false, "save favorites articles only")
	flag.Usage = func() {
		fmt.Printf("usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	version := flag.Bool("version", false, "get version")
	flag.Parse()
	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// only true value rewrite config value
	if *favorites {
		cfg.Favorites = true
	}

	return cfg, nil
}

// ----------------------------------------------------------------------------
func save_config(cfg Config) error {
	json_cfg, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(os.Getenv("HOME")+"/"+CONFIG_PATH, json_cfg, 0600)
	if err != nil {
		return err
	}

	return nil
}
