package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	activityTagTemplate := "ActivityID=%s"
	dateTimeTagTemplate := "DateTime=%s"
	measurementName := "ExampleMeasurement"
	fieldKey := "val"

	actTagCount := 2000
	dateTagCount := 2000

	f, err := os.Create("lpdata.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	nowTime := time.Now()
	lower := nowTime.Add(-(24 * 365) * time.Hour)
	max := int(nowTime.UnixNano())
	min := int(lower.UnixNano())

	wroteLines := 0

	for o := 0; o < actTagCount; o++ {
		actTag := fmt.Sprintf(activityTagTemplate, String(30))
		for i := 0; i < dateTagCount; i++ {
			ts := rand.Intn(max-min) + min
			dateTag := fmt.Sprintf(dateTimeTagTemplate, String(30))
			lp := fmt.Sprintf("%s,%s,%s %s=%f %d\n", measurementName, actTag, dateTag, fieldKey, rand.Float64(), ts)
			f.WriteString(lp)
			wroteLines++
			if wroteLines > 0 && wroteLines%100000 == 00 {
				fmt.Printf("wrote %d lines\n", wroteLines)
			}
		}
	}

	f.Sync()
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
