package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	sites := []string{
		"https://www.google.com",
		"https://drive.gogle.com", // there's a missing 'o' in google!!
		"https://maps.google.com",
		"https://hangouts.google.com",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(len(sites))

	for idx, site := range sites {
		go func(site string, idx int) {
			defer wg.Done()

			res, err := http.Get(site)
			if err != nil {
				cancel()
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Nanosecond):
				io.WriteString(os.Stdout, strconv.Itoa(idx)+": "+site+"\t\t "+res.Status+"\n")
			}
		}(site, idx)

		// time.Sleep(time.Second)
	}

	wg.Wait()
}
