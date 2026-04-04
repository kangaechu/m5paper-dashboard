package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kangaechu/m5paper-dashboard/internal/render"
)

const (
	forecastURL = "https://www.jma.go.jp/bosai/forecast/data/forecast/%s.json"
)

// jmaForecast represents the top-level JMA forecast response.
type jmaForecast struct {
	TimeSeries []timeSeries `json:"timeSeries"`
}

type timeSeries struct {
	TimeDefines []string `json:"timeDefines"`
	Areas       []area   `json:"areas"`
}

type area struct {
	Area struct {
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"area"`
	WeatherCodes []string `json:"weatherCodes"`
	Weathers     []string `json:"weathers"`
	Pops         []string `json:"pops"`
	Temps        []string `json:"temps"`
	TempsMin     []string `json:"tempsMin"`
	TempsMax     []string `json:"tempsMax"`
}

// weatherCodeToWMO maps JMA weather codes to simplified WMO-like codes for icons.
var weatherCodeToWMO = map[string]int{
	"100": 0, "101": 1, "102": 61, "103": 61, "104": 61,
	"110": 1, "111": 1, "112": 61, "113": 61,
	"114": 61, "115": 71, "116": 71,
	"117": 71, "118": 61, "119": 61,
	"120": 61, "121": 61, "122": 61, "123": 1, "124": 1,
	"125": 61, "126": 61, "127": 61, "128": 61,
	"130": 1, "131": 1, "132": 1,
	"140": 61, "160": 71, "170": 71,
	"200": 2, "201": 2, "202": 61, "203": 61, "204": 61,
	"205": 61, "206": 61, "207": 61, "208": 61,
	"209": 2, "210": 2, "211": 2, "212": 61, "213": 61,
	"214": 61, "215": 71, "216": 71,
	"217": 71, "218": 61, "219": 61,
	"220": 61, "221": 61, "222": 61, "223": 2, "224": 61,
	"225": 61, "226": 61, "228": 61,
	"230": 2, "231": 2, "240": 61, "250": 71,
	"260": 71, "270": 71,
	"300": 61, "301": 61, "302": 61, "303": 71, "304": 61,
	"306": 65, "308": 65, "309": 71,
	"311": 61, "313": 61, "314": 61, "315": 71,
	"316": 71, "317": 71, "320": 61,
	"321": 61, "322": 71, "323": 61, "324": 61, "325": 61,
	"326": 61, "327": 61, "328": 61,
	"329": 61, "340": 71, "350": 61,
	"361": 71, "371": 71,
	"400": 71, "401": 71, "402": 71, "403": 71,
	"405": 75, "406": 75, "407": 75,
	"409": 61, "411": 71, "413": 71,
	"414": 71, "420": 71, "421": 71,
	"422": 71, "423": 71, "425": 75,
	"426": 75, "427": 75, "450": 61,
}

// weatherCodeToDescription maps JMA codes to short Japanese descriptions.
var weatherCodeToDescription = map[string]string{
	"100": "晴れ", "101": "晴時々曇", "102": "晴時々雨", "110": "晴後曇", "112": "晴後雨",
	"200": "曇り", "201": "曇時々晴", "202": "曇時々雨", "210": "曇後晴", "212": "曇後雨",
	"300": "雨", "301": "雨時々晴", "302": "雨時々曇", "303": "雨時々雪", "311": "雨後晴", "313": "雨後曇",
	"400": "雪", "401": "雪時々晴", "402": "雪時々曇", "403": "雪時々雨",
}

func descriptionForCode(code string) string {
	if d, ok := weatherCodeToDescription[code]; ok {
		return d
	}
	return "---"
}

func wmoForCode(code string) int {
	if w, ok := weatherCodeToWMO[code]; ok {
		return w
	}
	return 0
}

// Fetch retrieves weather data from JMA API.
func Fetch(locationCode string, now time.Time) (*render.WeatherData, error) {
	url := fmt.Sprintf(forecastURL, locationCode)
	resp, err := http.Get(url) //nolint:gosec // URL is constructed from a constant format string
	if err != nil {
		return nil, fmt.Errorf("fetch JMA forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JMA API returned %d", resp.StatusCode)
	}

	var forecasts []jmaForecast
	if err := json.NewDecoder(resp.Body).Decode(&forecasts); err != nil {
		return nil, fmt.Errorf("decode JMA forecast: %w", err)
	}

	if len(forecasts) == 0 {
		return nil, fmt.Errorf("empty forecast response")
	}

	return parseForecasts(forecasts)
}

func parseForecasts(forecasts []jmaForecast) (*render.WeatherData, error) {
	data := &render.WeatherData{}

	// First report: short-term forecast
	if len(forecasts) > 0 {
		report := forecasts[0]

		// TimeSeries[0]: weather codes and descriptions (3 days)
		if len(report.TimeSeries) > 0 {
			ts := report.TimeSeries[0]
			if len(ts.Areas) > 0 {
				a := ts.Areas[0]
				dayLabels := []string{"今日", "明日", "明後日"}
				for i := 0; i < len(a.WeatherCodes) && i < 3; i++ {
					code := a.WeatherCodes[i]
					f := render.DayForecast{
						DayLabel:    dayLabels[min(i, len(dayLabels)-1)],
						Description: descriptionForCode(code),
						WeatherCode: wmoForCode(code),
					}
					data.Forecasts = append(data.Forecasts, f)
				}

				// Current weather from today's forecast
				if len(a.WeatherCodes) > 0 {
					data.WeatherCode = wmoForCode(a.WeatherCodes[0])
					data.Description = descriptionForCode(a.WeatherCodes[0])
				}
			}
		}

		// TimeSeries[1]: precipitation probability (6-hour intervals)
		if len(report.TimeSeries) > 1 {
			ts := report.TimeSeries[1]
			if len(ts.Areas) > 0 {
				a := ts.Areas[0]
				// Find the current/next pop
				for i, pop := range a.Pops {
					if pop == "" {
						continue
					}
					if v, err := strconv.Atoi(pop); err == nil {
						data.PrecipChance = v
						// Also assign to day forecasts
						if i < len(data.Forecasts) {
							data.Forecasts[0].PrecipChance = v
						}
						break
					}
				}
				// Assign max precip per day to forecasts
				assignPrecipToForecasts(data, ts)
			}
		}

		// TimeSeries[2]: temperature
		// temps format: [today_min, today_current, tomorrow_min, tomorrow_max] (varies by time of day)
		if len(report.TimeSeries) > 2 {
			ts := report.TimeSeries[2]
			if len(ts.Areas) > 0 {
				a := ts.Areas[0]
				// Parse all valid temps
				var parsedTemps []float64
				for _, t := range a.Temps {
					if v, err := strconv.ParseFloat(t, 64); err == nil {
						parsedTemps = append(parsedTemps, v)
					}
				}
				if len(parsedTemps) > 0 {
					data.Temperature = parsedTemps[0]
				}
				// Use first report temps for today's min/max
				if len(parsedTemps) >= 2 && len(data.Forecasts) > 0 {
					tMin := parsedTemps[0]
					tMax := parsedTemps[1]
					if tMin > tMax {
						tMin, tMax = tMax, tMin
					}
					data.Forecasts[0].TempMin = tMin
					data.Forecasts[0].TempMax = tMax
				}
			}
		}
	}

	// Second report: extended forecast (temps min/max)
	if len(forecasts) > 1 {
		report := forecasts[1]
		if len(report.TimeSeries) > 1 {
			ts := report.TimeSeries[1]
			if len(ts.Areas) > 0 {
				a := ts.Areas[0]
				for i := 0; i < len(data.Forecasts) && i < len(a.TempsMin); i++ {
					if a.TempsMin[i] != "" {
						if v, err := strconv.ParseFloat(a.TempsMin[i], 64); err == nil {
							data.Forecasts[i].TempMin = v
						}
					}
					if i < len(a.TempsMax) && a.TempsMax[i] != "" {
						if v, err := strconv.ParseFloat(a.TempsMax[i], 64); err == nil {
							data.Forecasts[i].TempMax = v
						}
					}
				}
			}
		}

		// Extended forecast precip probability
		if len(report.TimeSeries) > 0 {
			ts := report.TimeSeries[0]
			if len(ts.Areas) > 0 {
				a := ts.Areas[0]
				for i := 1; i < len(data.Forecasts) && i < len(a.Pops); i++ {
					if v, err := strconv.Atoi(a.Pops[i]); err == nil {
						data.Forecasts[i].PrecipChance = v
					}
				}
			}
		}
	}

	return data, nil
}

func assignPrecipToForecasts(data *render.WeatherData, ts timeSeries) {
	if len(ts.Areas) == 0 {
		return
	}
	a := ts.Areas[0]
	// Find max precip for each day
	for i, td := range ts.TimeDefines {
		t, err := time.Parse(time.RFC3339, td)
		if err != nil || i >= len(a.Pops) || a.Pops[i] == "" {
			continue
		}
		v, err := strconv.Atoi(a.Pops[i])
		if err != nil {
			continue
		}
		// Determine which forecast day this belongs to
		for fi := range data.Forecasts {
			if fi >= len(data.Forecasts) {
				break
			}
			// Simple heuristic: first 2 entries are today, next 2 tomorrow, etc.
			dayOffset := fi
			forecastDate := time.Now().AddDate(0, 0, dayOffset)
			if t.Day() == forecastDate.Day() && t.Month() == forecastDate.Month() {
				if v > data.Forecasts[fi].PrecipChance {
					data.Forecasts[fi].PrecipChance = v
				}
			}
		}
	}
}
