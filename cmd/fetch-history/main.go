// fetch-history fetches historical Ure Dam storage rate data from opengov.jp
// and writes it to the dam_history.json cache file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kangaechu/m5paper-dashboard/internal/dam"
	"github.com/kangaechu/m5paper-dashboard/internal/render"
)

const opengovURL = "https://opengov.jp/geo/dam-reservoir/ure/"

func main() {
	cacheFile := flag.String("cache", "dam_history.json", "output cache file path")
	years := flag.String("years", "", "comma-separated years to fetch (default: all available)")
	flag.Parse()

	fmt.Println("Fetching historical data from opengov.jp...")

	data, err := fetchOpenGov()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch error: %v\n", err)
		os.Exit(1)
	}

	// Filter years if specified
	var yearFilter map[string]bool
	if *years != "" {
		yearFilter = make(map[string]bool)
		for _, y := range strings.Split(*years, ",") {
			yearFilter[strings.TrimSpace(y)] = true
		}
	}

	// Load existing cache
	history, err := dam.LoadHistory(*cacheFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cache load warning: %v (starting fresh)\n", err)
		history = make(map[string][]render.DailyStorageRate)
	}

	// Merge fetched data into history
	totalNew := 0
	for year, entries := range data {
		if yearFilter != nil && !yearFilter[year] {
			continue
		}

		existing := make(map[string]bool)
		for _, e := range history[year] {
			existing[e.Date] = true
		}

		added := 0
		for _, e := range entries {
			if !existing[e.Date] {
				history[year] = append(history[year], e)
				added++
			}
		}

		// Sort
		sort.Slice(history[year], func(i, j int) bool {
			return history[year][i].Date < history[year][j].Date
		})

		if added > 0 {
			fmt.Printf("  %s: %d entries added (total: %d)\n", year, added, len(history[year]))
			totalNew += added
		}
	}

	if totalNew == 0 {
		fmt.Println("No new data to add.")
		return
	}

	if err := dam.SaveHistory(*cacheFile, history); err != nil {
		fmt.Fprintf(os.Stderr, "save error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Saved %d new entries to %s\n", totalNew, *cacheFile)
}

// yearData is the JSON structure embedded in the opengov.jp page.
type yearData struct {
	Labels []string  `json:"labels"` // "MM-DD"
	Data   []float64 `json:"data"`   // storage rate %
}

// jsonDataRegexp finds the embedded JSON object with yearly dam data.
var jsonDataRegexp = regexp.MustCompile(`"(\d{4})":\{"labels":\[`)

func fetchOpenGov() (map[string][]render.DailyStorageRate, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(opengovURL)
	if err != nil {
		return nil, fmt.Errorf("fetch opengov: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("opengov returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read opengov: %w", err)
	}

	// The HTML contains &#34; encoded JSON. Unescape it.
	content := html.UnescapeString(string(body))

	// Find the JSON object containing yearly data.
	// The structure is: {"2005":{"labels":[...],"data":[...]}, "2006":...}
	// Find the start of this JSON by looking for the first year key
	loc := jsonDataRegexp.FindStringIndex(content)
	if loc == nil {
		return nil, fmt.Errorf("could not find yearly data in HTML")
	}

	// Back up to find the opening brace
	start := loc[0] - 1
	for start > 0 && content[start] != '{' {
		start--
	}

	// Find the matching closing brace
	depth := 0
	end := start
	for i := start; i < len(content); i++ {
		if content[i] == '{' {
			depth++
		} else if content[i] == '}' {
			depth--
			if depth == 0 {
				end = i + 1
				break
			}
		}
	}

	jsonStr := content[start:end]

	// Parse the JSON
	var rawData map[string]yearData
	if err := json.Unmarshal([]byte(jsonStr), &rawData); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	// Convert to our format
	result := make(map[string][]render.DailyStorageRate)
	for year, yd := range rawData {
		if len(yd.Labels) != len(yd.Data) {
			fmt.Fprintf(os.Stderr, "warning: %s has %d labels but %d data points\n",
				year, len(yd.Labels), len(yd.Data))
			continue
		}

		yearInt, err := strconv.Atoi(year)
		if err != nil {
			continue
		}

		var entries []render.DailyStorageRate
		for i, label := range yd.Labels {
			// label is "MM-DD"
			dateStr := fmt.Sprintf("%04d-%s", yearInt, strings.ReplaceAll(label, "-", "-"))
			// Validate date
			if _, err := time.Parse("2006-01-02", dateStr); err != nil {
				continue
			}
			entries = append(entries, render.DailyStorageRate{
				Date:        dateStr,
				StorageRate: yd.Data[i],
			})
		}

		if len(entries) > 0 {
			result[year] = entries
		}
	}

	return result, nil
}
