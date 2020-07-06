package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/pabloos/http/greet"
)

func Test_Greet(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name                  string
		Greet                 greet.Greet
		Wanted                string
		ExpectedContentLength int64
	}{
		{
			Name: "Green case",
			Greet: greet.Greet{
				Name:     "John Doe",
				Location: "NY",
			},
			ExpectedContentLength: 24,
			Wanted:                "Hello John Doe, from NY\n",
		},
		{
			Name: "void case",
			Greet: greet.Greet{
				Name:     "",
				Location: "",
			},
			ExpectedContentLength: 55,
			Wanted:                "Tell us what is your name and where do you come from!\n",
		},
		{
			Name: "missing Name",
			Greet: greet.Greet{
				Name:     "",
				Location: "OH",
			},
			ExpectedContentLength: 55,
			Wanted:                "Tell us what is your name and where do you come from!\n",
		},
		{
			Name: "missing Location",
			Greet: greet.Greet{
				Name:     "David Doyle",
				Location: "",
			},
			ExpectedContentLength: 55,
			Wanted:                "Tell us what is your name and where do you come from!\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			t.Logf("Running '%s' with payload: %#v", tc.Name, tc.Greet)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.Greet)
			if err != nil {
				t.Fatalf("Problem while encoding JSON: " + err.Error())
			}

			req, err := http.NewRequest("POST", "localhost:8080/greet", buf)
			if err != nil {
				t.Fatalf("Couldn't create request: %v", err)
			}

			rec := httptest.NewRecorder()

			greetHandler(rec, req)

			res := rec.Result()
			res.Body.Close()

			msgBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Couldn't read response: %v", err)
			}

			msg := string(msgBytes)
			cl := res.ContentLength
			clh, clhErr := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
			t.Logf("Response Headers: %+v", res.Header)
			t.Logf("Response ContentLength: %d", cl)
			if msg != tc.Wanted {
				t.Errorf("%s failed: wanted %s, get %s", tc.Name, tc.Wanted, msg)
			}
			if cl != tc.ExpectedContentLength {
				t.Errorf("%s failed: ContentLength wanted %d, get %d", tc.Name, tc.ExpectedContentLength, cl)
			}
			if clhErr != nil {
				t.Errorf("%s failed: Header 'Content-Length' wanted %d, get ERROR %s", tc.Name, tc.ExpectedContentLength, clhErr.Error())
			} else if clh != tc.ExpectedContentLength {
				t.Errorf("%s failed: Header 'Content-Length' wanted %d, get %d", tc.Name, tc.ExpectedContentLength, clh)
			}
		})
	}
}
