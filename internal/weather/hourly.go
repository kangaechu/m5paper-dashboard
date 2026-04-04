package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kangaechu/m5paper-dashboard/internal/render"
)

const (
	openMeteoURL = "https://api.open-meteo.com/v1/jma?latitude=%s&longitude=%s&hourly=temperature_2m,weather_code,precipitation_probability,wind_speed_10m,wind_direction_10m&wind_speed_unit=ms&timezone=Asia%%2FTokyo&forecast_hours=25"
)

type openMeteoResponse struct {
	Hourly struct {
		Time                 []string  `json:"time"`
		Temperature2m        []float64 `json:"temperature_2m"`
		WeatherCode          []int     `json:"weather_code"`
		PrecipitationProb    []int     `json:"precipitation_probability"`
		WindSpeed10m         []float64 `json:"wind_speed_10m"`
		WindDirection10m     []int     `json:"wind_direction_10m"`
	} `json:"hourly"`
}

// FetchHourly retrieves hourly weather data from Open-Meteo JMA API.
// Returns 12 entries at 2-hour intervals covering 24 hours.
func FetchHourly(lat, lon string, now time.Time) ([]render.HourlyWeather, error) {
	url := fmt.Sprintf(openMeteoURL, lat, lon)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch open-meteo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("open-meteo API returned %d", resp.StatusCode)
	}

	var data openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode open-meteo: %w", err)
	}

	return parseHourly(data, now)
}

func parseHourly(data openMeteoResponse, now time.Time) ([]render.HourlyWeather, error) {
	h := data.Hourly
	if len(h.Time) == 0 {
		return nil, fmt.Errorf("no hourly data")
	}

	loc := now.Location()

	// Find the index closest to current hour
	currentHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)
	startIdx := -1
	for i, ts := range h.Time {
		t, err := time.ParseInLocation("2006-01-02T15:04", ts, loc)
		if err != nil {
			continue
		}
		if !t.Before(currentHour) {
			startIdx = i
			break
		}
	}
	if startIdx < 0 {
		startIdx = 0
	}

	// Pick every 2 hours, 12 entries
	var result []render.HourlyWeather
	for i := startIdx; i < len(h.Time) && len(result) < 12; i += 2 {
		t, err := time.ParseInLocation("2006-01-02T15:04", h.Time[i], loc)
		if err != nil {
			continue
		}

		hw := render.HourlyWeather{
			Hour: t.Hour(),
		}
		if i < len(h.Temperature2m) {
			hw.Temperature = h.Temperature2m[i]
		}
		if i < len(h.WeatherCode) {
			hw.WeatherCode = wmoForCode(fmt.Sprintf("%d", h.WeatherCode[i]))
			// Open-Meteo uses WMO codes directly, map them
			hw.WeatherCode = h.WeatherCode[i]
		}
		if i < len(h.PrecipitationProb) {
			hw.PrecipProb = h.PrecipitationProb[i]
		}
		if i < len(h.WindSpeed10m) {
			hw.WindSpeed = h.WindSpeed10m[i]
		}
		if i < len(h.WindDirection10m) {
			hw.WindDirection = h.WindDirection10m[i]
		}

		result = append(result, hw)
	}

	return result, nil
}
