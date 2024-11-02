package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func scraper(url string, pagecount int) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var filename string
	var url_ string
	if pagecount != 0 {
		filename = saveAs + fmt.Sprintf("%02d", pagecount) + ".json"
		url_ = url + fmt.Sprint(pagecount)
	} else {
		filename = saveAs + ".json"
		url_ = url
	}

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		printLog(fmt.Sprintf("File %s already exists. Skipping...", filename))
		return
	}

	chromedp.ListenTarget(
		ctx,
		func(ev interface{}) {
			if ev, ok := ev.(*network.EventResponseReceived); ok {

				// Check for HTTP 429 status, e.g too many requests
				if ev.Response.Status == 429 {
					fmt.Println("Received HTTP 429 - Too Many Requests. Waiting for 15 seconds...")
					time.Sleep(15 + time.Duration(rand.Intn(15)))
					return
				}

				if ev.Type != "XHR" {
					return
				}
				if !strings.Contains(ev.Response.URL, "ajax-sport-country-") {
					return
				}
				time.Sleep(time.Second * 3)

				go func() {
					c := chromedp.FromContext(ctx)
					rbp := network.GetResponseBody(ev.RequestID)
					body, err := rbp.Do(cdp.WithExecutor(ctx, c.Target))
					if err != nil {
						fmt.Println(err)
						pagecount--
						fmt.Println("RUN IT AGAIN..")
						return
					}

					var pageData struct {
						D struct {
							Rows    []json.RawMessage `json:"rows"`
							Total   int               `json:"total"`
							OnePage int               `json:"onepage"`
							Page    int               `json:"page"`
						} `json:"d"`
					}

					if err := json.Unmarshal(body, &pageData); err != nil {
						fmt.Println("Error unmarshaling JSON:", err)
						return
					}

					TotalPages = int(math.Ceil(float64(pageData.D.Total) / float64(pageData.D.OnePage)))
					fmt.Printf("Scraping Page %v out of %v..\n", pageData.D.Page, TotalPages)

					rowsJSON, err := json.Marshal(pageData.D.Rows)
					if err != nil {
						fmt.Println("Error marshaling rows data:", err)
						return
					}

					err = os.WriteFile(filename, rowsJSON, 0644)
					if err != nil {
						fmt.Println("Error writing file:", err)
						return
					}

					fmt.Printf("SAVED: %v \n", filename)
					fmt.Println("waiting..")
				}()
			}
		},
	)

	err := chromedp.Run(ctx,
		network.Enable(),
		network.SetExtraHTTPHeaders(HEADERS),
		chromedp.Navigate(url_),
		chromedp.Sleep(time.Second*time.Duration((10+rand.Intn(15)))),
	)
	if err != nil {
		fmt.Println(err)
	}
}

func runBase() {
	TotalPages = 1
	for i := 1; i <= TotalPages; i++ {
		printLog(fmt.Sprintf("CYCLE: %v.. TARGET: %v", i, url+fmt.Sprintf("%v", i)))
		scraper(url, i)
	}
}

func runBaseDaily() {
	printLog(fmt.Sprintf("CYCLE: %v.. TARGET: %v", 1, url+fmt.Sprintf("%v", 1)))
	scraper(url, 1)
}
