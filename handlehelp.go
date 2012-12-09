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

type Website struct {
	Pattern string
	UserURL string
	Name string
}

type handleResult struct {
	Site      string
	Available bool
}

func checkHandle(handle string, site Website) (bool) {

	valid, err := regexp.MatchString(site.Pattern, handle)

	if !valid || err != nil {
		return false
	}

	res, err := http.Get(fmt.Sprintf(site.UserURL, handle))
	
	if err != nil {
		return false
	}
	
	return res.StatusCode == 404
}


func main() {
	tmpl, err := template.ParseFiles("index.html")

	if err != nil {
		log.Fatal(err)
	}

	sites := []Website{
		Website{Name: "twitter", UserURL: "http://twitter.com/%s", Pattern: "^[a-zA-Z0-9_]{1,15}$"},
		Website{Name: "reddit", UserURL: "http://reddit.com/user/%s", Pattern: "^[a-zA-Z0-9_-]{3,20}$"},
		Website{Name: "hacker news", UserURL: "http://news.ycombinator.com/user?id=%s", Pattern: "^[a-zA-Z0-9_-]{2,15}$"},
		Website{Name: "lobsters", UserURL: "https://lobste.rs/u/%s", Pattern: "^[a-zA-Z0-9_-]{2,20}$"},
		Website{Name: "dribble", UserURL: "https://dribble.com/%s", Pattern: "^[a-zA-Z0-9_-]{2,20}$"},
		Website{Name: "forrst", UserURL: "https://forrst.com/people/%s", Pattern: "^[a-zA-Z0-9_]{1,20}$"},
		Website{Name: "facebook", UserURL: "https://facebook.com/%s", Pattern: "^[a-zA-Z0-9\\.]{3,50}$"},
		Website{Name: "youtube", UserURL: "https://www.youtube.com/user/%s", Pattern: "^[a-zA-Z0-9]{3,50}$"},
	}

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

		for _, site := range sites {
			go func(site Website) {
				a := checkHandle(r.FormValue("handle"), site)
				cm <- handleResult{Site: site.Name, Available: a}
			}(site)
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
