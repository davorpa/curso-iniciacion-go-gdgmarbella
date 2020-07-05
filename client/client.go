package main

import (
	"net/http"
	"time"
)

func newClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,

		Transport: &http.Transport{
			TLSClientConfig: tlsConfig(),
		},
	}
}
