package dam

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kangaechu/m5paper-dashboard/internal/render"
)

const (
	// DefaultURL is the water.go.jp real-time data page for Ure Dam.
	DefaultURL = "https://www.water.go.jp/mizu/chubu/realtime/p020201_60/301_1.html"

	// Ure Dam constants
	EffectiveCapacity = 28420.0 // 有効貯水容量 (×10³m³)
	NormalWaterLevel  = 229.15  // 常時満水位 (EL.m)
	MinWaterLevel     = 178.85  // 最低水位 (EL.m)
)

// tdRegexp extracts content from a <td> tag.
var tdRegexp = regexp.MustCompile(`<td[^>]*>(.*?)</td>`)

// timeWithDateRegexp matches "MM/DD HH:MM"
var timeWithDateRegexp = regexp.MustCompile(`(\d{2}/\d{2})\s+(\d{2}:\d{2})`)

// timeOnlyRegexp matches "HH:MM" alone
var timeOnlyRegexp = regexp.MustCompile(`^(\d{2}:\d{2})$`)

// commentRegexp matches HTML comments like <!---->
var commentRegexp = regexp.MustCompile(`<!--.*?-->`)

// Fetch retrieves current dam data from the given URL.
func Fetch(url string, now time.Time) (*render.DamData, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch dam data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dam data returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read dam data: %w", err)
	}

	return parseHTML(string(body), now)
}

func parseHTML(html string, now time.Time) (*render.DamData, error) {
	year := now.Year()
	loc := now.Location()

	// Split by <tr to process each row
	rows := strings.Split(html, "<tr")

	var history []render.DamObservation
	var currentDate string // tracks MM/DD across rows

	for _, row := range rows {
		// Extract all <td> cells from this row
		cells := tdRegexp.FindAllStringSubmatch(row, -1)
		if len(cells) < 7 {
			continue
		}

		// First cell should be the time cell
		timeCell := cleanCell(cells[0][1])
		if timeCell == "" {
			continue
		}

		var datePart, timePart string

		if m := timeWithDateRegexp.FindStringSubmatch(timeCell); len(m) == 3 {
			datePart = m[1]
			timePart = m[2]
			currentDate = datePart
		} else if m := timeOnlyRegexp.FindStringSubmatch(timeCell); len(m) == 2 {
			timePart = m[1]
			datePart = currentDate
		} else {
			continue
		}

		if datePart == "" || timePart == "" {
			continue
		}

		t, err := time.ParseInLocation("2006/01/02 15:04",
			fmt.Sprintf("%d/%s %s", year, datePart, timePart), loc)
		if err != nil {
			continue
		}
		// Handle year boundary
		if t.After(now.Add(24 * time.Hour)) {
			t = t.AddDate(-1, 0, 0)
		}

		// cells[1] = 貯水位, cells[2] = 有効貯水量
		// cells[3], cells[4] = empty
		// cells[5] = 流入量, cells[6] = 放流量(利水)
		waterLevel := parseFloat(cleanCell(cells[1][1]))
		storage := parseFloat(cleanCell(cells[2][1]))
		inflow := parseFloat(cleanCell(cells[5][1]))
		outflow := parseFloat(cleanCell(cells[6][1]))

		if waterLevel == 0 && storage == 0 {
			continue
		}

		history = append(history, render.DamObservation{
			Time:             t,
			WaterLevel:       waterLevel,
			EffectiveStorage: storage,
			Inflow:           inflow,
			Outflow:          outflow,
		})
	}

	if len(history) == 0 {
		return nil, fmt.Errorf("no valid dam observations parsed")
	}

	// The last row is the most recent observation (data is chronological)
	latest := history[len(history)-1]
	storageRate := latest.EffectiveStorage / EffectiveCapacity * 100

	// Reverse history so newest is first
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return &render.DamData{
		Name:             "宇連ダム",
		ObservedAt:       latest.Time,
		WaterLevel:       latest.WaterLevel,
		EffectiveStorage: latest.EffectiveStorage,
		StorageRate:      storageRate,
		Inflow:           latest.Inflow,
		Outflow:          latest.Outflow,
		History:          history,
	}, nil
}

func cleanCell(s string) string {
	s = commentRegexp.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	return s
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
