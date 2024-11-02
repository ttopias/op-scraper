package main

import "github.com/chromedp/cdproto/network"

const (
	// Constants
	BASEURL         = "https://www.oddsportal.com"
	MAX_RETRIES     = 10
	MAX_SLEEP       = 3
	MIN_MICRO_SLEEP = 50
	MAX_MICRO_SLEEP = 300

	// XPATH's
	ODDS_TABLE         = `div[data-v-49199a7b]`
	LINE_BUTTONS       = `ul.visible-links.bg-black-main.odds-tabs.flex.w-full > li.text-white-main.odds-item`
	FT_LINE_BUTTON     = `Array.from(document.querySelectorAll('div.tab-wrapper > div.flex-center.bg-gray-medium.h-\\[30px\\].cursor-pointer.px-3')).find(el => el.textContent.trim() === 'Full Time')?.click();`
	HIDDENLINE_BUTTONS = `ul.hidden-links.no-scrollbar.links-invisible > li`
	MORE_BUTTON        = `div.text-white-main.ml-auto.flex.items-center.p-3.pb-\\[14px\\].pl-3.pr-1.text-xs > .drop-arrow`
	BOOKMAKER_CELL     = `div[data-v-0e9f6ffa].border-black-borders.flex.h-9:nth-child(%d) > div:nth-child(1) > :nth-child(2) > p`
	BOOKMAKER_CELL_TC  = `div[data-v-0e9f6ffa].border-black-borders.flex.h-9 > div:nth-child(1) > :nth-child(2) > p`
	FIRST_CELL         = `div[data-v-0e9f6ffa].border-black-borders.flex.h-9:nth-child(%d) > div:nth-child(2) > div > div > p`
	SECOND_CELL        = `div[data-v-0e9f6ffa].border-black-borders.flex.h-9:nth-child(%d) > div:nth-child(3) > div > div > p`
	THIRD_CELL         = `div[data-v-0e9f6ffa].border-black-borders.flex.h-9:nth-child(%d) > div:nth-child(4) > div > div > p`
	FIRST_CELL_TC      = `div[data-v-0e9f6ffa].border-black-borders.flex.h-9 > :nth-child(2)`
	SECOND_CELL_TC     = `div[data-v-0e9f6ffa].border-black-borders.flex.h-9 > div:nth-child(3) > div > div > p`
	THIRD_CELL_TC      = `div[data-v-0e9f6ffa].border-black-borders.flex.h-9 > div:nth-child(4) > div > div > p`
	FOURTH_CELL_TC     = `div[data-v-0e9f6ffa].border-black-borders.flex.h-9 > div:nth-child(5) > div > div > p`
)

var BOOKMAKERS_TO_SCRAPE = []string{"pinnacle", "bet365", "betfair", "unibet"}
var TWO_COL_CELLS = []string{"#home-away;1", "#over-under;1", "#over-under;2", "#ah;1", "#ah;2", "#bts;2", "#dnb;2", "#odd-even;1"}
var THREE_COL_CELLS = []string{"#1X2", "#double", "#eh"}
var DONT_SCRAPE = []string{
	"#eh;1",
	"#eh;2",
	"#cs;1",
	"#cs;2",
	"#odd-even;1",
	"#odd-even;2",
}

var HEADERS = network.Headers{
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"Accept-Language":           "en-US,en;q=0.5",
	"Cache-Control":             "no-cache",
	"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Sec-Ch-Ua":                 `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
	"Sec-Ch-Ua-Mobile":          "?0",
	"Sec-Ch-Ua-Platform":        `"Windows"`,
	"Sec-Fetch-Dest":            "document",
	"Sec-Fetch-Mode":            "navigate",
	"Sec-Fetch-Site":            "none",
	"Sec-Fetch-User":            "?1",
	"Upgrade-Insecure-Requests": "1",
}

// Scraping vars
var TotalPages int
var url string
var saveAs string
var mode string
var filePath string
var strictMode bool
var outputAsCSV bool
var toFile bool
var isDebug bool
