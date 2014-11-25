package main

import (
	"fmt"
    "github.com/ChimeraCoder/anaconda"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    anaconda.SetConsumerKey("-")
    anaconda.SetConsumerSecret("your-consumer-secret")
    api := anaconda.NewTwitterApi("your-access-token", "your-access-token-secret")
    searchResult, _ := api.GetSearch("golang", nil)
    for _ , tweet := range searchResult {
        fmt.Println(tweet.Text)
    }
}
