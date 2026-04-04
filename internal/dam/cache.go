package dam

import (
	"encoding/json"
	"os"
	"sort"
	"time"

	"github.com/kangaechu/m5paper-dashboard/internal/render"
)

// LoadHistory reads the yearly history cache from a JSON file.
// Returns an empty map if the file does not exist.
func LoadHistory(path string) (map[string][]render.DailyStorageRate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string][]render.DailyStorageRate), nil
		}
		return nil, err
	}

	var history map[string][]render.DailyStorageRate
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}
	return history, nil
}

// SaveHistory writes the yearly history cache to a JSON file.
func SaveHistory(path string, history map[string][]render.DailyStorageRate) error {
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// UpdateHistory adds or updates today's storage rate in the history.
func UpdateHistory(history map[string][]render.DailyStorageRate, now time.Time, rate float64) {
	yearKey := now.Format("2006")
	dateStr := now.Format("2006-01-02")

	entries := history[yearKey]

	// Update existing entry or append
	found := false
	for i, e := range entries {
		if e.Date == dateStr {
			entries[i].StorageRate = rate
			found = true
			break
		}
	}
	if !found {
		entries = append(entries, render.DailyStorageRate{
			Date:        dateStr,
			StorageRate: rate,
		})
	}

	// Sort by date
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date < entries[j].Date
	})

	history[yearKey] = entries
}
