package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Website struct {
	Pattern     string
	UserURL     string
	Name        string
	NotFound    string
	RegisterURL string
}

type Config struct {
	Sites []Website
}

type handleResult struct {
	Site      Website
	Available bool
	Profile string
}

func readConfig() Config {
	if len(os.Args) != 2 {
		log.Fatal("You must supply a configuration filename")
	}
	filename := os.Args[1]
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var c Config
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		panic(err)
	}
	return c
}

func checkHandle(handle string, site Website) bool {

	valid, err := regexp.MatchString(site.Pattern, handle)

	if !valid || err != nil {
		return false
	}

	res, err := http.Get(fmt.Sprintf(site.UserURL, handle))

	if err != nil {
		return false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return false
	}

	if site.NotFound == "" {
		return res.StatusCode == 404 || string(body) == ""
	}

	return strings.Contains(string(body), site.NotFound)
}

func main() {
	tmpl, err := template.ParseFiles("index.html")

	if err != nil {
		log.Fatal(err)
	}

	config := readConfig()

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

		go func() {
			datum := map[string]interface{}{
				"event": "search",
				"properties": map[string]interface{}{
					"token":  os.Getenv("HANDLE_MIXPANEL_TOKEN"),
					"ip":     r.Header.Get("X-Forwarded-For"),
					"handle": r.FormValue("handle"),
				},
			}

			b, err := json.Marshal(datum)

			if err != nil {
				return
			}

			v := url.Values{}
			v.Set("data", base64.StdEncoding.EncodeToString(b))
			_, _ = http.Get("http://api.mixpanel.com/track/?" + v.Encode())
		}()

		for _, site := range config.Sites {
			go func(site Website) {
				a := checkHandle(r.FormValue("handle"), site)
				profile := fmt.Sprintf(site.UserURL, r.FormValue("handle"))
				cm <- handleResult{Site: site, Available: a, Profile: profile}
			}(site)
		}

		for _, _ = range config.Sites {
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
