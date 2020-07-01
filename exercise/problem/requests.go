package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	sites := []string{
		"https://www.google.com",
		"https://drive.google.com",
		"https://maps.google.com",
		"https://hangouts.google.com",
		"https://www.youtube.com",
		"https://meetup.com",
		"https://www.udc.es",
		"https://github.com/davorpa",
		"https://www.docker.com",
		"http://amazon.es",
		"https://twitter.com",
		"https://gobyexample.com",
	}

	for idx, site := range sites {
		wg.Add(1)

		go func(site string, idx int) {
			// `defer` is a way to DRY the sync notification, on an http response error or not.
			// it force the execution after goroutine function has exit
			defer wg.Done()

			res, err := http.Get(site)
			if err != nil {
				io.WriteString(os.Stderr, strconv.Itoa(idx)+": "+site+"\t\t FETCH FAIL: "+err.Error()+"\n")
				return
			}

			io.WriteString(os.Stdout, strconv.Itoa(idx)+": "+site+"\t\t "+res.Status+"\n")
		}(site, idx)
	}

	wg.Wait()
}
