package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	ApiKey       string `json:"api_key"`
	ApiSecret    string `json:"api_secret"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
}

func LoadConfig(w http.ResponseWriter) Config {
	reader, err := os.Open("./config.json")
	if err != nil {
		fmt.Fprint(w, "opening config file: %s", err.Error())
	}

	dec := json.NewDecoder(reader)
	var cnf Config
	err = dec.Decode(&cnf)
	if err != nil {
		fmt.Fprint(w, "json decode error: %s", err.Error())
	}

	return cnf
}
