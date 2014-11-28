package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	ApiKey       string `json:"api_key"`
	ApiSecret    string `json:"api_secret"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
}

func LoadConfig() (Config, error) {
	var cnf Config
	reader, err := os.Open("./config.json")
	if err != nil {
		err := errors.New(fmt.Sprintf("opening config file: ", err.Error()))
		return cnf, err
	}

	dec := json.NewDecoder(reader)
	err = dec.Decode(&cnf)
	if err != nil {
		err := errors.New(fmt.Sprintf("json decode error: ", err.Error()))
		return cnf, err
	}

	return cnf, nil
}
