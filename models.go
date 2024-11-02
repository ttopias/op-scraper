package main

import "github.com/chromedp/cdproto/cdp"

type Match struct {
	ID                      int      `json:"id"`
	URL                     string   `json:"url"`
	IsDouble                bool     `json:"is-double"`
	Home                    int      `json:"home"`
	Away                    int      `json:"away"`
	HomeName                string   `json:"home-name"`
	AwayName                string   `json:"away-name"`
	HomeCountryTwoChartName string   `json:"home-country-two-chart-name"`
	AwayCountryTwoChartName string   `json:"away-country-two-chart-name"`
	HomeParticipantID       int      `json:"home-participant-id"`
	AwayParticipantID       int      `json:"away-participant-id"`
	StatusID                int      `json:"status-id"`
	EventStageID            int      `json:"event-stage-id"`
	EventStageName          string   `json:"event-stage-name"`
	TournamentStageID       int      `json:"tournament-stage-id"`
	TournamentStageTypeID   int      `json:"tournament-stage-type-id"`
	TournamentStageGroupID  int      `json:"tournament-stage-group-id"`
	TournamentStageName     string   `json:"tournament-stage-name"`
	SportID                 int      `json:"sport-id"`
	Cols                    string   `json:"cols"`
	CountryID               int      `json:"country-id"`
	CountryName             string   `json:"country-name"`
	CountryTwoChartName     string   `json:"country-two-chart-name"`
	CountryType             string   `json:"country-type"`
	TournamentID            int      `json:"tournament-id"`
	TournamentName          string   `json:"tournament-name"`
	TournamentURL           string   `json:"tournament-url"`
	HomeParticipantImages   []string `json:"home-participant-images"`
	AwayParticipantImages   []string `json:"away-participant-images"`
	SportURLName            string   `json:"sport-url-name"`
	Breadcrumbs             struct {
		Sport struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"sport"`
		Country struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"country"`
		Tournament struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"tournament"`
	} `json:"breadcrumbs"`
	EncodeEventID        string `json:"encodeEventId"`
	ColClassName         string `json:"colClassName"`
	HomeParticipantTypes []int  `json:"homeParticipantTypes"`
	AwayParticipantTypes []int  `json:"awayParticipantTypes"`
	DateStartBase        int    `json:"date-start-base"`
	DateStartTimestamp   int    `json:"date-start-timestamp"`
	Date                 string `json:"date"`
	Result               string `json:"result"`
	HomeResult           string `json:"homeResult"`
	AwayResult           string `json:"awayResult"`
	HomeWinner           string `json:"home-winner"`
	AwayWinner           string `json:"away-winner"`
	Partialresult        string `json:"partialresult"`
	BookmakersCount      int    `json:"bookmakersCount"`
	WinnerPost           int    `json:"winner_post"`
	BettingType          int    `json:"betting_type"`
	Odds                 []struct {
		AvgOdds           float64 `json:"avgOdds"`
		BettingTypeID     int     `json:"bettingTypeId"`
		EventID           int     `json:"eventId"`
		MaxOdds           float64 `json:"maxOdds"`
		OutcomeResultID   int     `json:"outcomeResultId"`
		ScopeID           int     `json:"scopeId"`
		OutcomeID         string  `json:"outcomeId"`
		MaxOddsProviderID int     `json:"maxOddsProviderId"`
		Active            bool    `json:"active"`
	} `json:"odds"`
	OddsData         map[string][]OddRow `json:"odds_data,omitempty"`
	Name             string              `json:"name"`
	ColClassNameTime string              `json:"colClassNameTime"`
}

type CSVMatch struct {
	ID                  int        `json:"id"`
	URL                 string     `json:"url"`
	HomeName            string     `json:"home-name"`
	AwayName            string     `json:"away-name"`
	Name                string     `json:"name"`
	EventStageName      string     `json:"event-stage-name"`
	TournamentStageName string     `json:"tournament-stage-name"`
	TournamentName      string     `json:"tournament-name"`
	Date                string     `json:"date"`
	DateStartTimestamp  int        `json:"date-start-timestamp"`
	Result              string     `json:"result"`
	HomeResult          string     `json:"homeResult"`
	AwayResult          string     `json:"awayResult"`
	Partialresult       string     `json:"partialresult"`
	Market              string     `json:"market"`
	Bookmaker           string     `json:"bookmaker"`
	Line                string     `json:"line"`
	LineValue           string     `json:"line_value"`
	Odd                 float64    `json:"odd"`
	OpeningOdd          OpeningOdd `json:"opening_odd"`
	OddsHistory         string     `json:"odds_history"`
}

// OddRow is the parsed odds row from the odds page
type OddRow struct {
	Bookmaker string     `json:"bookmaker"` // Pinnacle
	Line      string     `json:"line"`      // 1X2, -1.5, 5.5 etc.
	Payout    float64    `json:"payout"`    // 0.95, e.g margin
	OddsData  []OddsData `json:"oddsData"`
}

type OddsData struct {
	LineValue   string        `json:"lineValue"` // 1/X/2, -1.5, 5.5 etc.
	Odd         float64       `json:"odd"`       // 1.95
	OpeningOdd  OpeningOdd    `json:"openingOdd"`
	OddsHistory []OddsHistory `json:"oddsHistory"` // All odds history for the specific line type
}

// RawOddRow is the raw data from the odds page
type RawOddRow struct {
	Bookmaker  string `json:"bookmaker"`  // Pinnacle
	FirstCell  string `json:"firstCell"`  // 1.95 or -1.5 (AH etc)
	SecondCell string `json:"secondCell"` // 3.40
	ThirdCell  string `json:"thirdCell"`  // 3.90
	FourthCell string `json:"fourthCell"` // 1.95, appears only for #eh;2
}

// OpeningOdd is the opening odds for a match
type OpeningOdd struct {
	Date string  `json:"date"` // 2024-01-01
	Odds float64 `json:"odds"` // 1.95
}

// OddsHistory is the odds history for a match
type OddsHistory struct {
	Date   string  `json:"date"`   // 2024-01-01
	Odds   float64 `json:"odds"`   // 1.95
	Change string  `json:"change"` // +0.05
}

// OddPageNodes is the nodes from the odds page, representing a single row of odds
type OddPageNodes struct {
	Bookmakers  []*cdp.Node
	FirstCells  []*cdp.Node
	SecondCells []*cdp.Node
	ThirdCells  []*cdp.Node
}
