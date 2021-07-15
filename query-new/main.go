package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type httpQuery struct {
	Query string
}

func main() {
	token := os.Getenv("TOKEN")

	if token == "" {
		panic("TOKEN must be set as env var")
	}

	b, _ := ioutil.ReadFile("query.txt")
	u, _ := url.Parse("http://localhost:8086/api/v2/query?org=org")

	q := httpQuery{
		Query: string(b),
	}

	payload, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}

	var maxTime time.Duration
	minTime := 1 * time.Hour // hopefully a request takes less time than this!
	var totTime time.Duration
	numReqs := 100

	for i := 0; i < numReqs; i++ {
		req := &http.Request{
			URL:    u,
			Method: "POST",
			Header: map[string][]string{
				"Authorization": {fmt.Sprintf("Token %s", token)},
			},
			Body: ioutil.NopCloser(bytes.NewReader(payload)),
		}

		start := time.Now()
		res, err := http.DefaultClient.Do(req)
		elapsed := time.Since(start)
		fmt.Printf("Request took %s\n", elapsed)

		if err != nil {
			fmt.Println(err)
		} else {
			if res.StatusCode != 200 {
				fmt.Println(res.Status)
			}
		}

		totTime = totTime + elapsed
		if elapsed > maxTime {
			maxTime = elapsed
		}

		if elapsed < minTime {
			minTime = elapsed
		}

		if i == numReqs-1 {
			b, _ := ioutil.ReadAll(res.Body)
			fmt.Println("*** final response ***")
			fmt.Println(string(b))
		}
	}

	avg := time.Duration(totTime.Nanoseconds() / int64(numReqs))

	fmt.Printf("Did %d queries\n", numReqs)
	fmt.Println("Max time:", maxTime)
	fmt.Println("Min time:", minTime)
	fmt.Println("Avg time:", avg)
}
