package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pabloos/http/greet"
)

const (
	URL = "https://localhost:8080"
)

func main() {
	client := newClient()

	dumpIndex(client)

	dumpGreet(client, false, greet.Greet{
		Name:     "",
		Location: "",
	})

	dumpGreet(client, false, greet.Greet{
		Name:     "John Doe",
		Location: "USA",
	})

	dumpGreet(client, true, greet.Greet{
		Name:     "Silvia Saint",
		Location: "USA",
	})

	dumpGreet(client, true, greet.Greet{
		Name:     "Donald Kennet",
		Location: "USA",
	})

	dumpGreet(client, true, greet.Greet{
		Name:     "Donald Kennet",
		Location: "SWE",
	})

	dumpGreet(client, true, greet.Greet{
		Name:     "",
		Location: "USA",
	})

	dumpGreet(client, true, greet.Greet{
		Name:     "David Lin",
		Location: "OH",
	})

	dumpGreet(client, true, greet.Greet{
		Name:     "",
		Location: "OH",
	})
}

func dumpIndex(client *http.Client) {
	resp1, err := client.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp1.Body.Close()

	io.Copy(os.Stdout, resp1.Body)
}

func dumpGreet(client *http.Client, cached bool, entity greet.Greet) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(entity)

	var greetURL string
	if cached {
		greetURL = fmt.Sprintf("%s/%s", URL, "greet-cached")
	} else {
		greetURL = fmt.Sprintf("%s/%s", URL, "greet")
	}

	resp, err := client.Post(greetURL, "application/json; charset=utf-8", buf)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	defer io.Copy(os.Stdout, resp.Body)

	fmt.Fprintf(os.Stdout, "%s :: %+v -> Status code: %d\n\n", greetURL, entity, resp.StatusCode)
}
