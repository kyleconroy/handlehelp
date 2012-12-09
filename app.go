package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
"regexp"
)

func reddit(handle string) (string, bool) {

	if len(handle) < 3 || len(handle) > 20 {
		return "reddit", false
	}

	res, err := http.Get("http://www.reddit.com/user/" + handle)
	
	if err != nil {
		return "reddit", false
	}
	
	return "reddit", res.StatusCode == 404
}

func twitter(handle string) (string, bool) {

	valid, err := regexp.MatchString("[a-zA-Z0-9_]{1,15}", handle)

	if !valid || err != nil {
		return "twitter", false
	}

	res, err := http.Get("http://www.twitter.com/" + handle)
	
	if err != nil {
		return "twitter", false
	}

	return "twitter", res.StatusCode == 404
}

type handleResult struct {
	Site      string
	Available bool
}

func main() {
	tmpl, err := template.ParseFiles("index.html")

	if err != nil {
		log.Fatal(err)
	}

	sites := []func(string) (string, bool){reddit, twitter}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		f, ok := w.(http.Flusher)

		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		cm := make(chan handleResult)

		for _, checkHandle := range sites {
			go func(ch func(string) (string, bool)) {
				s, a := ch(r.FormValue("handle"))
				cm <- handleResult{Site: s, Available: a}
			}(checkHandle)
		}

		for _, _ = range sites {
			hr, ok := <-cm

			if !ok {
				continue
			}

			b, err := json.Marshal(hr)

			if err != nil {
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", string(b))
			f.Flush()
		}
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	pfx := "/static/"

	h := http.StripPrefix(pfx, http.FileServer(http.Dir("static")))

	http.Handle(pfx, h)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
