package main

import (
	"appengine"
	"appengine/urlfetch"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/http"
	"net/url"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	cnf, err := LoadConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	anaconda.SetConsumerKey(cnf.ApiKey)
	anaconda.SetConsumerSecret(cnf.ApiSecret)
	api := anaconda.NewTwitterApi(cnf.AccessKey, cnf.AccessSecret)
	c := appengine.NewContext(r)
	api.HttpClient.Transport = &urlfetch.Transport{Context: c}
	v := url.Values{}
	v.Set("count", "30")
	searchResult, err := api.GetSearch("ご冥福をお祈り", v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, tweet := range searchResult {
		fmt.Fprintln(w, tweet.RetweetCount)
		fmt.Fprintln(w, tweet.Text)
		fmt.Fprintln(w, "<br />")
	}
}
