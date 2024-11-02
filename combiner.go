package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func combine() {
	path := filepath.FromSlash(filePath)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	var allMatches []Match

	for _, f := range files {
		if !f.Type().IsRegular() || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}

		fp := filepath.Join(path, f.Name())

		file, err := os.ReadFile(fp)
		if err != nil {
			log.Printf("error reading file: %v", err)
			continue
		}

		fmt.Printf("CHECKING AND MERGING: %v \n", fp)
		var matches []Match
		if err := json.Unmarshal(file, &matches); err != nil {
			log.Printf("error unmarshalling JSON: %v", err)
			continue
		}

		allMatches = append(allMatches, matches...)
	}

	if outputAsCSV {
		if err := processCSVFile(allMatches); err != nil {
			log.Printf("Error processing CSV file: %v", err)
		}
	} else {
		if err := processJSONFile(allMatches); err != nil {
			log.Printf("Error processing JSON file: %v", err)
		}
	}
}

func processCSVFile(matches []Match) error {
	var csvRows []CSVMatch

	for _, match := range matches {
		for marketType, marketData := range match.OddsData {
			for _, lineData := range marketData {
				for _, odd := range lineData.OddsData {
					csvRow := CSVMatch{
						ID:                  match.ID,
						URL:                 match.URL,
						HomeName:            retroTeamId(match.HomeName),
						AwayName:            retroTeamId(match.AwayName),
						Name:                match.Name,
						EventStageName:      match.EventStageName,
						TournamentStageName: match.TournamentStageName,
						TournamentName:      match.TournamentName,
						Date:                parseMatchDate(int64(match.DateStartBase)),
						DateStartTimestamp:  match.DateStartTimestamp,
						Result:              match.Result,
						HomeResult:          match.HomeResult,
						AwayResult:          match.AwayResult,
						Partialresult:       match.Partialresult,
						Market:              marketType,
						Bookmaker:           lineData.Bookmaker,
						Line:                lineData.Line,
						LineValue:           odd.LineValue,
						Odd:                 odd.Odd,
						OpeningOdd:          odd.OpeningOdd,
						OddsHistory:         formatOddsHistory(odd.OddsHistory),
					}
					csvRows = append(csvRows, csvRow)
				}
			}
		}
	}

	// Write to CSV file
	fn := saveAs + ".csv"
	if mode == "daily" {
		fn = saveAs + "daily.csv"
	}

	if err := writeCSV(fn, csvRows); err != nil {
		return fmt.Errorf("error writing CSV: %w", err)
	}

	fmt.Printf("Compiled %v rows of data\n", len(csvRows))
	fmt.Println("Wrote CSV File")
	return nil
}

func formatOddsHistory(history []OddsHistory) string {
	if history == nil {
		return ""
	}
	bytes, err := json.Marshal(history)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func writeCSV(filename string, rows []CSVMatch) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"OddsportalID", "URL", "HomeTeam", "AwayTeam", "Name",
		"EventStageName", "TournamentStageName", "TournamentName",
		"Date", "DateStartTimestamp", "Result", "HomeResult", "AwayResult",
		"Partialresult", "Market", "Bookmaker", "Line", "LineValue",
		"Odd", "OpeningOdd", "OddsHistory",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}

	// Write data rows
	for _, row := range rows {
		record := []string{
			fmt.Sprint(row.ID),
			row.URL,
			row.HomeName,
			row.AwayName,
			row.Name,
			row.EventStageName,
			row.TournamentStageName,
			row.TournamentName,
			row.Date,
			fmt.Sprint(row.DateStartTimestamp),
			row.Result,
			row.HomeResult,
			row.AwayResult,
			row.Partialresult,
			row.Market,
			row.Bookmaker,
			row.Line,
			row.LineValue,
			fmt.Sprint(row.Odd),
			fmt.Sprint(row.OpeningOdd),
			row.OddsHistory,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	return nil
}

func processJSONFile(matches []Match) error {
	countRows := 0
	var content []Match

	for _, match := range matches {
		match.HomeName = retroTeamId(match.HomeName)
		match.AwayName = retroTeamId(match.AwayName)
		match.Date = parseMatchDate(int64(match.DateStartBase))

		countRows++
		content = append(content, match)
	}

	if err := writeJSON(saveAs, content, countRows); err != nil {
		log.Fatalf("Error writing JSON: %v", err)
	}

	return nil
}

func writeJSON(saveAs string, ResultsCompiled_ []Match, countRows int) error {
	output, err := json.Marshal(ResultsCompiled_)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	if err := os.WriteFile(saveAs+".json", output, 0644); err != nil {
		return fmt.Errorf("error writing JSON file: %w", err)
	}

	fmt.Printf("Compiled %v rows of games \n", countRows)
	fmt.Println("Wrote JSON File")
	return nil
}
