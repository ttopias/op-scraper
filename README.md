# op-scraper

Scrapes the results and odds from OddsPortal.com using chromedp, stores them in JSON or CSV. Can scrape either base data or full data (results + odds).

## Disclaimer

This program is for educational purposes only. I am not responsible for any misuse or damage caused by this program. Use it at your own risk.

By using this software, you agree that:

1. You are solely responsible for how you use this program
2. You will comply with OddsPortal's terms of service and any applicable laws
3. The author(s) assume no liability for any direct, indirect, incidental, special, exemplary, or consequential damages
4. You understand that web scraping may be against the target website's terms of service
5. This software comes with no warranty of any kind, express or implied

## Requirements

- Go >=1.23.2
- Chrome Browser

## Usage

Build the program

```bash
go build
```

or with some env vars

```bash
GOOS=linux GOARCH=amd64 go build -o op-scraper
```

Then run the program with the following options:

```bash
git clone https://github.com/ttopias/op-scraper
cd op-scraper
go build
./op-scraper -h
```

Then run something like this:

```bash
./op-scraper -m full -u https://www.oddsportal.com/hockey/usa/nhl-2022-2023/results/#/page/ -s NHL_2022-2023_ -f ./results/2022
```

or as CRON job:

```bash
chmod +x /path/to/your/scrape_daily.sh && crontab -e

# add this at the bottom
0 12 * * * /path/to/your/scrape_daily.sh
```

Check scripts for examples.

## Run options

```bash
-m base/combine/match/full/daily
```

Defines the run mode, options: 'base', 'combine'

```bash
-m base -u "https://www.oddsportal.com/hockey/usa/nhl-2022-2023/results/#/page/"
-m match -u "https://www.oddsportal.com/hockey/usa/nhl-2022-2023/florida-panthers-vegas-golden-knights-EeQklJzr/"
-m full -u "https://www.oddsportal.com/hockey/usa/nhl-2022-2023/results/#/page/"
-m odds -f "./results/2022" -strict true
```

'base', then URL must end in ../#/page/
'match', then URL must be exact path for the specific match to scrape, without the line suffix.
'full', same as base, but also scrapes odds data and combines them into single file.
'daily', then scrapes all matches within 48 hours.
'odds', then path to folder with scraped 'base data' and it scrapes odds data to it

```bash
-s "NHL_2022-2023_"
```

Filename/Dir for saving scraped odds, program will add '01.json' etc.

```bash
-f, "", "Path to the JSON file for scraping the odds OR folder with jsons to combine"
```

Path to the JSON file for scraping the odds OR folder with jsons to combine

```bash
-odds false
```

Output to CSV, default: false (saves as nested JSON). CSV will contain all data in flat format, e.g. multiple rows for the same match.

```bash
-strict false
```

Strict mode, only scrape oddsdata from 'wanted' bookmakers, default: false. No real performance gain from this, just a bit cleaner data.

The 'wanted' bookmakers; pinnacle, bet365, betfair, unibet.

```bash
-d false
```

Run in debug mode, default: false. Prints a lot of logs to console.

## LICENSE

MIT License

Copyright (c) 2024 ttopias

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
