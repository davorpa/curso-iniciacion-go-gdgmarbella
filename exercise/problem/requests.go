package main

import (
	"io"
	"net/http"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	sites := []string{
		"https://www.google.com",
		"https://drive.google.com",
		"https://maps.google.com",
		"https://hangouts.google.com",
	}

	for _, site := range sites {
		wg.Add(1)

		go func(site string) {
			// `defer` is a way to DRY the sync notification, on an http response error or not.
			// it force the execution after goroutine function has exit
			defer wg.Done()

			res, err := http.Get(site)
			if err != nil {
				io.WriteString(os.Stderr, site+"\t\t FETCH FAIL: "+err.Error()+"\n")
				return
			}

			io.WriteString(os.Stdout, site+"\t\t "+res.Status+"\n")
		}(site)
	}

	wg.Wait()
}
