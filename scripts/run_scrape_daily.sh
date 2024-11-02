#! /bin/bash

set -e
exec 1> >(logger -s -t $(basename $0)) 2>&1

echo "Scraping the data"
./oddsportal-scraper -m "daily" -u https://www.oddsportal.com/hockey/usa/nhl/results/#/page/ -s "./daily/" -f "./daily" -o true

# Run some R / Python code to process the data and save to the database or so
echo "Processing the data"
/usr/bin/Rscript daily_oddsportal.R "./daily/daily.csv"

# Clean up the daily folder for the next day
rm -rf "./daily"
