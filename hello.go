package main

import (
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
	"appengine/user"
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Config struct {
	ApiKey       string `json:"api_key"`
	ApiSecret    string `json:"api_secret"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
}

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/users", users)
	http.HandleFunc("/sign", signHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
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

func users(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
	greetings := make([]Greeting, 0, 10)
	if _, err := q.GetAll(c, &greetings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := guestbookTemplate.Execute(w, greetings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	fmt.Fprintf(w, "Hello, %v!", u)
}

func signHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := Greeting{
		Content: r.FormValue("content"),
		Date:    time.Now(),
	}

	if u := user.Current(c); u != nil {
		g.Author = u.String()
	}

	key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))
	_, err := datastore.Put(c, key, &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}

func guestbookKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

var guestbookTemplate = template.Must(template.New("sign").Parse(guestbookTemplateHTML))

const guestbookTemplateHTML = `
<html>
  <body>
    {{range .}}
      {{with .Author}}
        <p><b>{{.}}</b> wrote:</p>
      {{else}}
        <p>An anonymous person wrote:</p>
      {{end}}
      <pre>{{.Content}}</pre>
    {{end}}
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`

type Greeting struct {
	Author  string
	Content string
	Date    time.Time
}
