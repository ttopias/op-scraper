package main

import (
	"flag"
)

func runFull() {
	toFile = false
	runBase()
	runMatchFull()
	combine()
}

func runDaily() {
	toFile = false
	runBaseDaily()
	runMatchFullDaily()
	combine()
}

func main() {
	toFile = true
	printLog("STARTING SCRAPER...")
	flag.StringVar(&mode, "m", "base", "Run mode: 'base', 'combine', 'match', 'full', 'daily', 'odds'")
	flag.StringVar(&url, "u", "https://www.oddsportal.com/hockey/usa/nhl-2022-2023/results/#/page/", "URL must end in ../#/page/")
	flag.StringVar(&saveAs, "s", "NHL_2023-2024_", "Filename/Dir for saving, will add 01.json")
	flag.StringVar(&filePath, "f", "", "Path to the JSON file for scraping the odds OR folder with jsons to combine")
	flag.BoolVar(&outputAsCSV, "o", false, "Output to CSV")
	flag.BoolVar(&strictMode, "strict", false, "Strict mode, only scrape wanted bookmakers")
	flag.BoolVar(&isDebug, "d", false, "Debug mode")
	flag.Parse()

	if mode == "base" {
		runBase()
	} else if mode == "combine" {
		combine()
	} else if mode == "match" {
		runMatch(url)
	} else if mode == "full" {
		runFull()
	} else if mode == "daily" {
		runDaily()
	} else if mode == "odds" {
		runMatchFull()
	} else {
		printLog("Error: Invalid mode. Please use '-h' to show options.")
	}

	printLog("SCRAPER FINISHED")
}
