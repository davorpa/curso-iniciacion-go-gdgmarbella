package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// A cancelable context to allow comunicate goroutines one each other
	ctx, cancel := context.WithCancel(context.Background())
	// Best Practice: defer context `cancel` to avoid thread leaks
	defer cancel()

	respChan := make(chan bool)

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
		"https://unkdnasd.google.com", // don't exists so...
		"https://www.yahoo.com",       // this next url is not fetched
	}

	for idx, site := range sites {

		go func(site string, idx int, ctx context.Context, respChan chan<- bool) {
			res, err := http.Get(site)
			if err != nil {
				io.WriteString(os.Stderr, strconv.Itoa(idx)+": "+site+"\t\t FETCH FAIL: "+err.Error()+"\n")
				// notify error response to channel
				respChan <- false
				return
			}

			io.WriteString(os.Stdout, strconv.Itoa(idx)+": "+site+"\t\t "+res.Status+"\n")
			// notify success response to channel
			respChan <- true
		}(site, idx, ctx, respChan)

		// test if channel notification fails
		if !<-respChan {
			// avoid fetch next urls
			cancel()
			break
		}
	}
}
