package main

import (
	"encoding/json"
	"fmt"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

type rest struct {
	Url string `json:"url"`
}

func main() {
	url := "https://httpbin.org/delay/1"

	// Instantiate default collector
	c := colly.NewCollector()

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		10, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("visiting", r.URL)
	//})

	c.OnResponse(func(r *colly.Response) {
		res := rest{}
		_ = json.Unmarshal(r.Body, &res)
		fmt.Println(res)
		fmt.Println("-----------------------------")
	})

	for i := 0; i < 1000; i++ {
		// Add URLs to the queue
		q.AddURL(fmt.Sprintf("%s?n=%d", url, i))
	}
	// Consume URLs
	q.Run(c)

}
