package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type httpQuery struct {
	Query string
}

func main() {
	u, err := url.Parse("http://localhost:8086/query")
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Set("db", "mydb")
	q.Set("epoch", "ms")
	q.Set("q", `SHOW TAG VALUES FROM "ExampleMeasurement" WITH KEY = "ActivityID"`)
	u.RawQuery = q.Encode()

	var maxTime time.Duration
	minTime := 1 * time.Hour // hopefully a request takes less time than this!
	var totTime time.Duration
	numReqs := 100

	for i := 0; i < numReqs; i++ {
		req := &http.Request{
			URL:    u,
			Method: "GET",
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
