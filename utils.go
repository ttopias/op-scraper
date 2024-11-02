package main

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"time"
)

func printLog(s string) {
	fmt.Printf("%s - LOG:\t%s\n", time.Now().Format("2024/01/01 13:45:00"), s)
}

func printDebug(s string) {
	if isDebug {
		fmt.Printf("%s - DEBUG:\t%s\n", time.Now().Format("2024/01/01 13:45:00"), s)
	}
}

// Retry a function a number of times with a random sleep between attempts,
// returns an error if the function fails after MAX_RETRIES attempts.
func retry(f func() error) (err error) {
	for i := 0; i < MAX_RETRIES; i++ {
		if i > 0 {
			printDebug(fmt.Sprintf("DEBUG: Sleeping for %d microseconds", MIN_MICRO_SLEEP+rand.Intn(MAX_MICRO_SLEEP)))
			microSleep()
		}

		done := make(chan error)
		go func() {
			done <- f()
		}()

		select {
		case err = <-done:
			if err == nil {
				return nil
			}
		case <-time.After(MIN_MICRO_SLEEP * time.Millisecond):
			err = fmt.Errorf("function execution timed out")
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", MAX_RETRIES, err)
}

func microSleep() {
	n := rand.Intn(MAX_MICRO_SLEEP)
	printDebug(fmt.Sprintf("DEBUG: Sleeping for %d microseconds", MIN_MICRO_SLEEP+n))
	time.Sleep(MIN_MICRO_SLEEP + time.Duration(n))
}

func calculatePayout(odds []OddsData) float64 {
	if len(odds) == 0 {
		return 0
	}

	product := 1.0
	for _, odd := range odds {
		if odd.Odd == 0 {
			continue
		}

		product *= odd.Odd
	}
	return 1 - (1 / product)
}

func parseLineValue(name string) string {
	switch name {
	case "#1X2;2":
		return "1X2"
	case "#home-away;1":
		return "ML"
	case "#over-under;1":
		return "OU-ML"
	case "#over-under;2":
		return "OU-FT"
	case "#ah;1":
		return "AH-ML"
	case "#ah;2":
		return "AH-FT"
	case "#bts;2":
		return "BTTS"
	case "#double;2":
		return "DC"
	case "#eh;2":
		return "EH"
	case "#dnb;2":
		return "DNB"
	}
	return ""
}

func getLineValue(line string, i int) string {
	switch line {
	case "1X2":
		switch i {
		case 1:
			return "1"
		case 2:
			return "X"
		case 3:
			return "2"
		}
	case "ML":
		switch i {
		case 1:
			return "1"
		case 2:
			return "2"
		}
	case "ML-OU":
		switch i {
		case 2:
			return "Over"
		case 3:
			return "Under"
		}
	case "FT-OU":
		switch i {
		case 2:
			return "Over"
		case 3:
			return "Under"
		}
	case "AH-ML":
		switch i {
		case 2:
			return "1"
		case 3:
			return "2"
		}
	case "AH-FT":
		switch i {
		case 2:
			return "1"
		case 3:
			return "2"
		}
	case "DNB":
		switch i {
		case 1:
			return "1"
		case 2:
			return "2"
		}
	case "BTTS":
		switch i {
		case 1:
			return "Yes"
		case 2:
			return "No"
		}
	case "DC":
		switch i {
		case 1:
			return "1X"
		case 2:
			return "12"
		case 3:
			return "X2"
		}
	case "EH":
		switch i {
		case 1:
			return "1"
		case 2:
			return "X"
		case 3:
			return "2"
		}
	}

	return ""
}

func parseURLSuffix(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func isWantedBookmaker(s string) bool {
	return slices.Contains(BOOKMAKERS_TO_SCRAPE, strings.ToLower(s))
}

func isTwoColumn(s string) bool {
	return slices.Contains(TWO_COL_CELLS, strings.ToLower(s))
}

func isThreeColumn(s string) bool {
	return slices.Contains(THREE_COL_CELLS, strings.ToLower(s))
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func getCellValue(r *RawOddRow, cell int) float64 {
	switch cell {
	case 1:
		return parseFloat(r.FirstCell)
	case 2:
		return parseFloat(r.SecondCell)
	case 3:
		return parseFloat(r.ThirdCell)
	}

	return 0
}

func parseRowData(r *RawOddRow, s string) OddRow {
	o := OddRow{}

	o.Bookmaker = r.Bookmaker
	o.Line = parseLineValue(s)
	if s == "#home-away;1" || s == "#home-away;2" || s == "#bts;2" || s == "#dnb;2" {
		o.OddsData = append(o.OddsData, getOddsData(parseLineValue(s), r, 1))
		o.OddsData = append(o.OddsData, getOddsData(parseLineValue(s), r, 2))
	} else if isTwoColumn(s) {
		o.Line = r.FirstCell
		o.OddsData = append(o.OddsData, getOddsData(parseLineValue(s), r, 2))
		o.OddsData = append(o.OddsData, getOddsData(parseLineValue(s), r, 3))
	} else {
		for i := 1; i <= 3; i++ {
			o.OddsData = append(o.OddsData, getOddsData(parseLineValue(s), r, i))
		}
	}

	o.Payout = calculatePayout(o.OddsData)
	return o
}

func getOddsData(s string, r *RawOddRow, cell int) OddsData {
	return OddsData{
		LineValue: getLineValue(s, cell),
		Odd:       getCellValue(r, cell),
	}
}

func isWithinTwoDays(date int64) bool {
	matchDate := time.Unix(date, 0)

	return time.Since(matchDate) <= 2*24*time.Hour
}

func filterMatches(matches []Match) []Match {
	filteredMatches := []Match{}
	for _, match := range matches {
		if isWithinTwoDays(int64(match.DateStartTimestamp)) {
			filteredMatches = append(filteredMatches, match)
		}
	}
	return filteredMatches
}

// date is in format '1537317000' and target format is '2006-01-02 15:04:00'
func parseMatchDate(date int64) string {
	return time.Unix(date, 0).Format("2006-01-02 15:00")
}

func retroTeamId(name string) string {
	switch name {
	// NHL
	case "Anaheim Ducks":
		return "ANA"
	case "Arizona Coyotes":
		return "AZN"
	case "Boston Bruins":
		return "BOS"
	case "Buffalo Sabres":
		return "BUF"
	case "Calgary Flames":
		return "CGY"
	case "Carolina Hurricanes":
		return "CAR"
	case "Chicago Blackhawks":
		return "CHI"
	case "Colorado Avalanche":
		return "COL"
	case "Columbus Blue Jackets":
		return "CBJ"
	case "Dallas Stars":
		return "DAL"
	case "Detroit Red Wings":
		return "DET"
	case "Edmonton Oilers":
		return "EDM"
	case "Florida Panthers":
		return "FLA"
	case "Los Angeles Kings":
		return "LAK"
	case "Minnesota Wild":
		return "MIN"
	case "Montreal Canadiens":
		return "MTL"
	case "Nashville Predators":
		return "NSH"
	case "New Jersey Devils":
		return "NJD"
	case "New York Islanders":
		return "NYI"
	case "New York Rangers":
		return "NYR"
	case "Ottawa Senators":
		return "OTT"
	case "Philadelphia Flyers":
		return "PHI"
	case "Pittsburgh Penguins":
		return "PIT"
	case "San Jose Sharks":
		return "SJS"
	case "Seattle Kraken":
		return "SEA"
	case "St. Louis Blues":
		return "STL"
	case "Tampa Bay Lightning":
		return "TBL"
	case "Toronto Maple Leafs":
		return "TOR"
	case "Vancouver Canucks":
		return "VAN"
	case "Vegas Golden Knights":
		return "VGK"
	case "Washington Capitals":
		return "WSH"
	case "Winnipeg Jets":
		return "WPG"
	case "Utah":
		return "UTA"
	// Default if not found
	default:
		return name
	}
}
