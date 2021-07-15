package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	urlBase := os.Getenv("URL_BASE")

	if urlBase == "" {
		panic("TOKEN and URL_BASE must be set as env vars")
	}

	u, _ := url.Parse(fmt.Sprintf("%s/write?db=mydb", urlBase))

	activityTagTemplate := "ActivityID=%s"
	dateTimeTagTemplate := "DateTime=%s"
	measurementName := "ExampleMeasurement"
	fieldKey := "val"

	actTagCount := 2000
	dateTagCount := 2000

	nowTime := time.Now()
	lower := nowTime.Add(-(24 * 365) * time.Hour)
	max := int(nowTime.UnixNano())
	min := int(lower.UnixNano())

	wroteLines := 0
	var batchedLp string

	for o := 0; o < actTagCount; o++ {
		actTag := fmt.Sprintf(activityTagTemplate, String(30))
		for i := 0; i < dateTagCount; i++ {
			ts := rand.Intn(max-min) + min
			dateTag := fmt.Sprintf(dateTimeTagTemplate, String(30))
			lp := fmt.Sprintf("%s,%s,%s %s=%f %d\n", measurementName, actTag, dateTag, fieldKey, rand.Float64(), ts)

			batchedLp = batchedLp + lp

			wroteLines++
			if wroteLines > 0 && wroteLines%10000 == 0 {
				req := &http.Request{
					URL:    u,
					Method: "POST",
					Body:   ioutil.NopCloser(strings.NewReader(batchedLp)),
				}
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					panic(err)
				} else {
					if res.StatusCode != 204 {
						fmt.Println(res.Status)
						b, _ := ioutil.ReadAll(res.Body)
						fmt.Println(string(b))
					}
				}

				batchedLp = ""

				fmt.Printf("wrote %d lines (%d total)\n", wroteLines, actTagCount*dateTagCount)
			}
		}
	}

	fmt.Printf("done - wrote %d lines\n", wroteLines)
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
