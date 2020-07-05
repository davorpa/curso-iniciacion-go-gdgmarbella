package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/pabloos/http/greet"
)

func POST(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func Debug(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer h.ServeHTTP(w, r)

		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(dump))
	}
}

func Delay(delay time.Duration, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer h.ServeHTTP(w, r)

		time.Sleep(delay)
	}
}

// Cached decorator inspect a Greet resource into a cache to suggests
// some equivalents when some data are missing
func Cached(h http.HandlerFunc) http.HandlerFunc {
	// Create the cache struct
	c, err := greet.NewCache()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// decode body request using a temporal reader
		// and rewriting it to avoid handler chain errors
		var t greet.Greet
		b := bytes.NewBuffer(make([]byte, 0))
		reader := io.TeeReader(r.Body, b)
		r.Body = ioutil.NopCloser(b)
		err := json.NewDecoder(reader).Decode(&t)
		if err != nil {
			log.Print("Error decoding body: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// do handler chain
		h.ServeHTTP(w, r)

		// inspect cache and append custom messages depending on it status
		store := c.GetAll()
		gr, prs := c.SetIfAbsent(t)
		if prs {
			log.Printf("Resource %+v cached as: %+v", t, gr)
			fmt.Fprintf(w, "Hey %s!! Your cached location was: %s\n", t.Name, gr.Location)
		} else {
			log.Printf("Resource %+v is not cached", t)
			// use and slice to resolve suggestions
			var suggestions = make([]greet.Greet, 0)
			for _, v := range store {
				if v.Name != t.Name && strings.Contains(v.Name, t.Name) {
					suggestions = append(suggestions, v)
				}
			}
			n := len(suggestions)
			log.Printf("%d suggested resources: %+v", n, suggestions)
			if n > 0 {
				if n == 1 {
					gr = suggestions[0]
					fmt.Fprintf(w, "Are you %s from %s?\n", gr.Name, gr.Location)
				} else {
					fmt.Fprintf(w, "Are you some of these suggestions? %s\n", strings.Trim(fmt.Sprintf("%+v", suggestions), "[]"))
				}
			}
		}
		log.Println("CACHE STORE: ", store)
	}
}
