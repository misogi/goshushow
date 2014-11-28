package main

import (
	"appengine"
	"appengine/urlfetch"
	"appengine/user"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"html/template"
	"net/http"
	"net/url"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/sign", signHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	anaconda.SetConsumerKey("-")
	anaconda.SetConsumerSecret("-")
	api := anaconda.NewTwitterApi("-", "-")
	c := appengine.NewContext(r)
	api.HttpClient.Transport = &urlfetch.Transport{Context: c}
	v := url.Values{}
	v.Set("count", "30")
	searchResult, _ := api.GetSearch("ご冥福をお祈り", v)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, tweet := range searchResult {
		fmt.Fprintln(w, tweet.RetweetCount)
		fmt.Fprintln(w, tweet.Text)
		fmt.Fprintln(w, "<br />")
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
	signTemplate.Execute(w, r.FormValue("content"))
}

var signTemplate = template.Must(template.New("sign").Parse(signTemplateHTML))

const signTemplateHTML = `
<html>
  <body>
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
    <p>You wrote:</p>
    <pre>{{.}}</pre>
  </body>
</html>
`
