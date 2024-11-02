package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func clickButton(btn *cdp.Node) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		printDebug(fmt.Sprintf("Before clicking button: %s", url))

		err := chromedp.WaitVisible(LINE_BUTTONS).Do(ctx)
		if err != nil {
			return fmt.Errorf("error waiting for line buttons: %v", err)
		}

		err = chromedp.Click(btn.FullXPath(), chromedp.NodeVisible).Do(ctx)
		if err != nil {
			return fmt.Errorf("error clicking line button: %v", err)
		}

		// Update location just in case we navigate to a new page
		err = chromedp.Location(&url).Do(ctx)
		if err != nil {
			return fmt.Errorf("error getting location: %v", err)
		}

		printDebug(fmt.Sprintf("After clicking button: %s", url))
		return nil
	})
}

func expandAllSections() chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		printDebug("Expanding sections")

		err := chromedp.WaitVisible(ODDS_TABLE).Do(ctx)
		if err != nil {
			return fmt.Errorf("error waiting for expanding buttons: %v", err)
		}

		err = chromedp.EvaluateAsDevTools(`
			document.querySelectorAll('div.bg-provider-arrow.h-4.w-4.bg-center.bg-no-repeat').forEach(button => button.click());
		`, nil).Do(ctx)
		if err != nil {
			return fmt.Errorf("error expanding sections: %v", err)
		}

		printDebug("Expanded sections")
		return nil
	})
}

func hoverOverCell(xPath string, waitTooltip bool) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		printDebug("Hovering over cell...")
		timeoutCtx, cancel := context.WithTimeout(ctx, MAX_SLEEP*time.Second)
		defer cancel()

		done := make(chan bool, 1)
		go func() {
			err := retry(func() error {
				return chromedp.EvaluateAsDevTools(fmt.Sprintf(`
					document.querySelector('%s').dispatchEvent(new MouseEvent('mouseover', {
						'view': window,
						'bubbles': true,
						'cancelable': true
					}));
				`, xPath), nil).Do(ctx)
			})
			if err != nil {
				done <- false
				return
			}
			if waitTooltip {
				err = retry(func() error {
					return chromedp.WaitVisible(`[class*="tooltip"]`).Do(ctx)
				})
				if err != nil {
					done <- false
					return
				}
			}

			done <- true
		}()

		select {
		case visible := <-done:
			if !visible {
				return fmt.Errorf("error hovering over cell")
			}
			printDebug("Hovered over cell")
			return nil
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout waiting for tooltip to become visible")
		}
	})
}

func scrapeOddPageNodes(o *OddPageNodes) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		s := parseURLSuffix(url)
		if slices.Contains(DONT_SCRAPE, s) {
			printDebug(fmt.Sprintf("Skipping suffix %s", s))
			return nil
		}

		printDebug("Scraping odd page nodes")
		err := retry(func() error {
			return chromedp.Nodes(`
				div[data-v-0e9f6ffa].border-black-borders.flex.h-9 > div:first-child
			`, &o.Bookmakers).Do(ctx)
		})
		if err != nil {
			return fmt.Errorf("error getting bookmakers: %v", err)
		}

		err = retry(func() error {
			return chromedp.Nodes(`
				div[data-v-0e9f6ffa].border-black-borders.flex.h-9 > div:nth-child(2)
			`, &o.FirstCells).Do(ctx)
		})
		if err != nil {
			return fmt.Errorf("error getting first cells: %v", err)
		}

		err = retry(func() error {
			return chromedp.Nodes(`
				div[data-v-0e9f6ffa].border-black-borders.flex.h-9 > div:nth-child(3)
			`, &o.SecondCells).Do(ctx)
		})
		if err != nil {
			return fmt.Errorf("error getting second cells: %v", err)
		}

		err = retry(func() error {
			return chromedp.Nodes(`
				div[data-v-0e9f6ffa].border-black-borders.flex.h-9 > div:nth-child(4)
			`, &o.ThirdCells).Do(ctx)
		})
		if err != nil {
			return fmt.Errorf("error getting third cells: %v", err)
		}

		printDebug("Scraped odd page nodes")
		return nil
	})
}

func scrapeOUorAH(r *RawOddRow, row int) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// suf := parseURLSuffix(url)

		printDebug(fmt.Sprintf("Scraping OU or AH for row %d...", row))

		var n []*cdp.Node
		err := retry(func() error {
			return chromedp.Nodes(BOOKMAKER_CELL_TC, &n).Do(ctx)
		})
		if err != nil {
			return fmt.Errorf("error getting bookmakers: %v", err)
		}
		err = retry(func() error {
			return chromedp.Text(n[row].FullXPath(), &r.Bookmaker).Do(ctx)
		})
		if err != nil {
			return fmt.Errorf("error getting bookmakers: %v", err)
		}

		err = retry(func() error {
			return chromedp.Nodes(FIRST_CELL_TC, &n).Do(ctx)
		})
		if err != nil {
			return fmt.Errorf("error getting odds: %v", err)
		}
		err = retry(func() error {
			return chromedp.Text(n[row].FullXPath(), &r.FirstCell).Do(ctx)
		})
		printDebug(fmt.Sprintf("Scraped first cell for row %d, got %s", row, r.FirstCell))
		if err != nil {
			return fmt.Errorf("error getting odds: %v", err)
		}

		err = retry(func() error {
			return chromedp.Nodes(SECOND_CELL_TC, &n).Do(ctx)
		})
		if err != nil {
			return fmt.Errorf("error getting odds: %v", err)
		}

		err = retry(func() error {
			return chromedp.Text(n[row].FullXPath(), &r.SecondCell).Do(ctx)
		})
		if err != nil {
			return fmt.Errorf("error getting odds: %v", err)
		}

		err = retry(func() error {
			return chromedp.Nodes(THIRD_CELL_TC, &n).Do(ctx)
		})
		if err != nil {
			return fmt.Errorf("error getting odds: %v", err)
		}

		err = retry(func() error {
			return chromedp.Text(n[row].FullXPath(), &r.ThirdCell).Do(ctx)
		})
		if err != nil {
			return fmt.Errorf("error getting odds: %v", err)
		}

		return nil
	})
}

func scrapeOddPageRow(r *RawOddRow, row int, url string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		printDebug(fmt.Sprintf("Scraping odd page row %d...", row))
		suf := parseURLSuffix(url)
		var err error
		if suf == "#over-under;1" || suf == "#over-under;2" || suf == "#ah;1" || suf == "#ah;2" {
			err = scrapeOUorAH(r, row).Do(ctx)
		} else {
			err = retry(func() error {
				return chromedp.Text(fmt.Sprintf(BOOKMAKER_CELL, row+2), &r.Bookmaker).Do(ctx)
			})
			if err != nil {
				return fmt.Errorf("error getting bookmakers: %v", err)
			}

			printDebug(fmt.Sprintf("Scraped bookmaker for row %d, got %s", row, r.Bookmaker))

			if isThreeColumn(suf) || suf == "#home-away;1" || suf == "#home-away;2" || suf == "#bts;2" || suf == "#dnb;2" {
				err = retry(func() error {
					return chromedp.Text(fmt.Sprintf(FIRST_CELL, row+2), &r.FirstCell).Do(ctx)
				})
				printDebug(fmt.Sprintf("Scraped first cell for row %d, got %s", row, r.FirstCell))
				if err != nil {
					return fmt.Errorf("error getting line: %v", err)
				}
			} else {
				if suf != "#over-under;1" && suf != "#over-under;2" && suf != "#ah;1" && suf != "#ah;2" {
					var n []*cdp.Node
					err = retry(func() error {
						return chromedp.Nodes(FIRST_CELL_TC, &n).Do(ctx)
					})
					if err != nil {
						return fmt.Errorf("error getting odds: %v", err)
					}
					err = retry(func() error {
						return chromedp.Text(n[row].FullXPath(), &r.FirstCell).Do(ctx)
					})
					printDebug(fmt.Sprintf("Scraped first cell for row %d, got %s", row, r.FirstCell))
					if err != nil {
						return fmt.Errorf("error getting odds: %v", err)
					}
				}
			}

			err = retry(func() error {
				return chromedp.Text(fmt.Sprintf(SECOND_CELL, row+2), &r.SecondCell).Do(ctx)
			})
			printDebug(fmt.Sprintf("Scraped second cell for row %d, got %s", row, r.SecondCell))
			if err != nil {
				return fmt.Errorf("error getting odds: %v", err)
			}

			// Process third cell if it exists
			if suf != "#home-away;1" && suf != "#home-away;2" && suf != "#bts;2" && suf != "#dnb;2" {
				err = retry(func() error {
					return chromedp.Text(fmt.Sprintf(THIRD_CELL, row+2), &r.ThirdCell).Do(ctx)
				})
				printDebug(fmt.Sprintf("Scraped third cell for row %d, got %s", row, r.ThirdCell))
			}
		}
		if err != nil {
			return fmt.Errorf("error getting odds: %v", err)
		}

		printDebug(fmt.Sprintf("Scraped values: Bookmaker: %s, 1st: %s, 2nd: %s, 3rd: %s, 4th: %s\n", r.Bookmaker, r.FirstCell, r.SecondCell, r.ThirdCell, r.FourthCell))
		return nil
	})
}

func scrapeOddPageRows(rows *[]OddRow, nodes *OddPageNodes, s *string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		*s = parseURLSuffix(url)

		if slices.Contains(DONT_SCRAPE, *s) {
			printDebug(fmt.Sprintf("Skipping suffix %s", *s))
			return nil
		}

		for i := 0; i < len(nodes.Bookmakers); i++ {
			printDebug(fmt.Sprintf("Scraping odds row %d/%d", i+1, len(nodes.Bookmakers)))
			var rRow RawOddRow
			err := chromedp.Run(ctx,
				scrapeOddPageRow(&rRow, i, url),
			)
			if err != nil {
				return err
			}

			if strictMode && !isWantedBookmaker(rRow.Bookmaker) {
				continue
			}

			o := parseRowData(&rRow, *s)

			if *s == "OU-FT" || *s == "AH-FT" || *s == "OU-ML" || *s == "AH-ML" {
				o.OddsData[0].LineValue = "1"
				o.OddsData[1].LineValue = "2"
			}
			*rows = append(*rows, o)
		}

		return nil
	})
}

func scrapeURL(ctx context.Context, btn *cdp.Node, mode string) ([]OddRow, string, error) {
	var err error
	var o []OddRow
	var nodes OddPageNodes
	var s string
	if mode == "hidden" {
		err = chromedp.Run(ctx,
			hoverOverCell(MORE_BUTTON, false),
			clickButton(btn),
			expandAllSections(),
			scrapeOddPageNodes(&nodes),
			scrapeOddPageRows(&o, &nodes, &s),
		)
	} else if mode == "subpage" {
		err = chromedp.Run(ctx,
			expandAllSections(),
			scrapeOddPageNodes(&nodes),
			scrapeOddPageRows(&o, &nodes, &s),
		)
	} else {
		err = chromedp.Run(ctx,
			clickButton(btn),
			expandAllSections(),
			scrapeOddPageNodes(&nodes),
			scrapeOddPageRows(&o, &nodes, &s),
		)
	}
	if err != nil {
		return o, parseLineValue(s), err
	}
	return o, parseLineValue(s), nil
}

func scrapeOdds(url string) (map[string][]OddRow, error) {
	printLog(fmt.Sprintf("Starting to scrape odds for URL: %s", url))

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", !isDebug),
		chromedp.Flag("disable-gpu", !isDebug),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	if !isDebug {
		ctx, cancel = context.WithTimeout(ctx, MAX_SLEEP*time.Minute)
		defer cancel()
	}

	// Add rate limit handling
	var gotResponse bool
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*network.EventResponseReceived); ok {
			if ev.Response.Status == 429 {
				gotResponse = true
				printLog("Received HTTP 429 - Too Many Requests. Waiting before retry...")
				time.Sleep(15*time.Second + time.Duration(rand.Intn(15))*time.Second)
			}
		}
	})

	err := chromedp.Run(ctx,
		network.Enable(),
		network.SetExtraHTTPHeaders(HEADERS),
		chromedp.Navigate(url),
	)
	if err != nil || gotResponse {
		printLog(fmt.Sprintf("Error navigating to page: %v", err))
		return scrapeOdds(url) // Recursive retry
	}

	var lineButtons []*cdp.Node
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(LINE_BUTTONS),
		chromedp.Nodes(LINE_BUTTONS, &lineButtons),
	)
	if err != nil {
		printLog(fmt.Sprintf("Error getting suffixes: %v", err))
	}

	oddsData := make(map[string][]OddRow)
	for _, b := range lineButtons {
		o, s, err := scrapeURL(ctx, b, "visible")
		if err != nil {
			printLog(fmt.Sprintf("Error scraping URL %s: %v", url, err))
			// return nil, err
		}

		if o != nil {
			if s == "OU-ML" || s == "AH-ML" {
				for j := range o {
					o[j].OddsData[0].LineValue = "1"
					o[j].OddsData[1].LineValue = "2"
				}
			}
			oddsData[s] = append(oddsData[s], o...)
		}

		// SUBPAGE LOGIC

		if s == "OU-ML" || s == "AH-ML" {
			// Check if there is a subpage for this line
			printDebug(fmt.Sprintf("Checking for subpage button for %s\n", s))
			var subpageBtn []*cdp.Node
			err = chromedp.Run(ctx, chromedp.Evaluate(FT_LINE_BUTTON, &subpageBtn))
			if err != nil {
				printLog(fmt.Sprintf("Error navigating to subpage: %v", err))
				// continue
			}

			// Get the location of the subpage button
			var loc string
			err = chromedp.Run(ctx, chromedp.Location(&loc))
			if err != nil {
				printLog(fmt.Sprintf("Error getting location: %v", err))
				// continue
			}
			printDebug(fmt.Sprintf("\t\tLocation: %s\n", loc))
			lv := parseLineValue(parseURLSuffix(loc))

			// Scrape subpage
			od, _, err := scrapeURL(ctx, b, "subpage")
			if err != nil {
				printLog(fmt.Sprintf("Error scraping subpage: %v", err))
				// continue
			}

			if od != nil {
				for j := range od {
					od[j].OddsData[0].LineValue = "1"
					od[j].OddsData[1].LineValue = "2"
				}
				printDebug(fmt.Sprintf("\t\tScraped subpage %s..., \t %+v\n", lv, od[0]))
				oddsData[lv] = append(oddsData[lv], od...)
			}
		}
	}

	// Check if the site has MORE_BUTTON
	var hasMoreButton bool
	err = chromedp.Run(ctx, chromedp.Evaluate(fmt.Sprintf("document.querySelector('%s') !== null", MORE_BUTTON), &hasMoreButton))
	if err != nil {
		printLog(fmt.Sprintf("Error checking for MORE_BUTTON: %v", err))
	}

	if hasMoreButton {
		var hiddenLineButtons []*cdp.Node
		err = chromedp.Run(ctx,
			chromedp.Nodes(HIDDENLINE_BUTTONS, &hiddenLineButtons),
		)
		if err != nil {
			return nil, fmt.Errorf("error getting hidden suffixes: %v", err)
		}

		if len(hiddenLineButtons) > 0 {
			for _, b := range hiddenLineButtons[:len(hiddenLineButtons)-1] {
				o, s, err := scrapeURL(ctx, b, "hidden")
				if err != nil {
					return nil, err
				}

				if o != nil {
					oddsData[s] = append(oddsData[s], o...)
				}

				// SUBPAGE LOGIC

				if s == "OU-ML" || s == "OU-FT" {
					// Check if there is a subpage for this line
					printDebug(fmt.Sprintf("Checking for subpage button for %s\n", s))
					var subpageBtn []*cdp.Node
					err = chromedp.Run(ctx, chromedp.Evaluate(FT_LINE_BUTTON, &subpageBtn))
					if err != nil {
						printLog(fmt.Sprintf("Error navigating to subpage: %v", err))
						continue
					}

					// Get the location of the subpage button
					var loc string
					err = chromedp.Run(ctx, chromedp.Location(&loc))
					if err != nil {
						printLog(fmt.Sprintf("Error getting location: %v", err))
						continue
					}
					printDebug(fmt.Sprintf("\t\tLocation: %s\n", loc))
					lv := parseLineValue(parseURLSuffix(loc))

					// Scrape subpage
					od, _, err := scrapeURL(ctx, b, "subpage")
					if err != nil {
						return nil, err
					}

					if od != nil {
						for j := range od {
							od[j].OddsData[0].LineValue = "1"
							od[j].OddsData[1].LineValue = "2"
						}
						printDebug(fmt.Sprintf("\t\tScraped subpage %s..., \t %+v\n", lv, od[0]))
						oddsData[lv] = append(oddsData[lv], od...)
					}
				}
			}
		}
	}

	printLog(fmt.Sprintf("Successfully scraped and saved odds for URL: %s", url))
	return oddsData, nil
}

func runMatch(url string) {
	oddsData, err := scrapeOdds(url)
	if err != nil {
		printLog(fmt.Sprintf("Error scraping odds: %v", err))
	}

	data, err := json.MarshalIndent(oddsData, "", "  ")
	if err != nil {
		printLog(fmt.Sprintf("Error marshaling scraped data to JSON: %v", err))
	}

	err = os.WriteFile(saveAs, data, 0644)
	if err != nil {
		printLog(fmt.Sprintf("Error writing scraped data to file: %v", err))
	}
}

func runMatchFull() {
	path := filepath.FromSlash(filePath)
	files, err := os.ReadDir(path)
	if err != nil {
		printLog(fmt.Sprintf("Error finding JSON files: %v", err))
		return
	}

	if len(files) == 0 {
		printLog(fmt.Sprintf("No JSON files found matching pattern: %s*.json", saveAs))
		return
	}

	for i, f := range files {
		if !f.Type().IsRegular() || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}

		file := filepath.Join(path, f.Name())
		printLog(fmt.Sprintf("Processing file %d/%d: %s", i+1, len(files), file))

		data, err := os.ReadFile(file)
		if err != nil {
			printLog(fmt.Sprintf("Error reading file %s: %v", file, err))
			continue
		}

		var matches []Match
		err = json.Unmarshal(data, &matches)
		if err != nil {
			printLog(fmt.Sprintf("Error unmarshaling JSON from file %s: %v", file, err))
			continue
		}

		for j := range matches {
			printLog(fmt.Sprintf("Scraping odds for match %d/%d in file %s", j+1, len(matches), file))
			if len(matches[j].OddsData) > 0 {
				printLog(fmt.Sprintf("Odds data already exists for match %s, skipping", matches[j].URL))
				continue
			}

			oddsData, err := scrapeOdds(BASEURL + matches[j].URL)
			if err != nil {
				printLog(fmt.Sprintf("Error scraping odds for %s: %v", BASEURL+matches[j].URL, err))
			}
			matches[j].OddsData = oddsData
			matches[j].Date = parseMatchDate(int64(matches[j].DateStartTimestamp))

			// Save after each match
			updatedData, err := json.MarshalIndent(matches, "", "  ")
			if err != nil {
				printLog(fmt.Sprintf("Warning: Failed to marshal with indentation after match %d: %v. Trying without indentation...", j+1, err))
				updatedData, err = json.Marshal(matches)
				if err != nil {
					printLog(fmt.Sprintf("Error: Failed to marshal data even without indentation after match %d: %v", j+1, err))
					continue
				}
			}

			err = os.WriteFile(file, updatedData, 0644)
			if err != nil {
				printLog(fmt.Sprintf("Error writing updated data to file after match %d: %v", j+1, err))
				continue
			}
			printLog(fmt.Sprintf("Successfully saved progress after match %d/%d", j+1, len(matches)))

			microSleep()
		}

		printLog(fmt.Sprintf("Successfully processed file %d/%d: %s", i+1, len(files), file))
	}

	printLog("Finished processing all files")
}

func runMatchFullDaily() {
	printLog(fmt.Sprintf("Processing file %s", saveAs+"01.json"))

	// Read the file
	data, err := os.ReadFile(saveAs + "01.json")
	if err != nil {
		printLog(fmt.Sprintf("Error reading file %s: %v", saveAs+"01.json", err))
		return
	}

	// Unmarshal JSON data
	var matches []Match
	err = json.Unmarshal(data, &matches)
	if err != nil {
		printLog(fmt.Sprintf("Error unmarshaling JSON from file %s: %v", saveAs+"01.json", err))
		return
	}

	// Filter matches within two days
	matches = filterMatches(matches)

	// Process each match
	for j := range matches {
		printLog(fmt.Sprintf("Scraping odds for match %d/%d in file %s", j+1, len(matches), saveAs+"01.json"))
		if len(matches[j].OddsData) > 0 {
			printLog(fmt.Sprintf("Odds data already exists for match %s, skipping", matches[j].URL))
			continue
		}

		oddsData, err := scrapeOdds(BASEURL + matches[j].URL)
		if err != nil {
			printLog(fmt.Sprintf("Error scraping odds for %s: %v", BASEURL+matches[j].URL, err))
			continue
		}
		matches[j].OddsData = oddsData

		// Save after each match
		updatedData, err := json.MarshalIndent(matches, "", "  ")
		if err != nil {
			printLog(fmt.Sprintf("Warning: Failed to marshal with indentation after match %d: %v. Trying without indentation...", j+1, err))
			updatedData, err = json.Marshal(matches)
			if err != nil {
				printLog(fmt.Sprintf("Error: Failed to marshal data even without indentation after match %d: %v", j+1, err))
				continue
			}
		}

		err = os.WriteFile(saveAs+"01.json", updatedData, 0644)
		if err != nil {
			printLog(fmt.Sprintf("Error writing updated data to file after match %d: %v", j+1, err))
			continue
		}
		printLog(fmt.Sprintf("Successfully saved progress after match %d/%d", j+1, len(matches)))

		microSleep()
	}

	printLog("Finished processing all files")
}
