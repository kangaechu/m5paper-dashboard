package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"image/png"
	"os"
	"strings"
	"time"

	"github.com/kangaechu/m5paper-dashboard/internal/calendar"
	"github.com/kangaechu/m5paper-dashboard/internal/render"
	"github.com/kangaechu/m5paper-dashboard/internal/train"
	"github.com/kangaechu/m5paper-dashboard/internal/weather"
)

func main() {
	loadEnvFile(".env")

	output := flag.String("output", "output.png", "output PNG file path")
	locationCode := flag.String("location", envOrDefault("LOCATION_CODE", "130000"), "JMA location code")
	lat := flag.String("lat", envOrDefault("LOCATION_LAT", "35.6895"), "latitude for hourly weather")
	lon := flag.String("lon", envOrDefault("LOCATION_LON", "139.6917"), "longitude for hourly weather")
	googleCreds := flag.String("google-creds", "", "base64-encoded Google service account JSON key")
	calendarIDs := flag.String("calendars", envOrDefault("GOOGLE_CALENDAR_IDS", "primary"), "comma-separated Google Calendar IDs")
	trainLines := flag.String("trains", envOrDefault("TRAIN_LINES", "山手線:21/0,都営三田線:129/0,都営浅草線:128/0,都営新宿線:130/0,都営大江戸線:131/0,東京メトロ銀座線:132/0,東京メトロ丸ノ内線:133/0,東京メトロ日比谷線:134/0,東京メトロ東西線:135/0,東京メトロ千代田線:136/0,東京メトロ有楽町線:137/0,東京メトロ半蔵門線:138/0,東京メトロ南北線:139/0,東京メトロ副都心線:540/0"), "comma-separated name:code pairs")
	flag.Parse()

	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc)

	data := render.DashboardData{Now: now}

	// Fetch weather (daily from JMA)
	w, err := weather.Fetch(*locationCode, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "weather error: %v\n", err)
	} else {
		data.Weather = w
	}

	// Fetch hourly weather (from Open-Meteo)
	hourly, err := weather.FetchHourly(*lat, *lon, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hourly weather error: %v\n", err)
	} else if data.Weather != nil {
		data.Weather.Hourly = hourly
	}

	// Fetch train delay info
	var lineConfigs []train.LineConfig
	for _, entry := range strings.Split(*trainLines, ",") {
		parts := strings.SplitN(entry, ":", 2)
		if len(parts) == 2 {
			lineConfigs = append(lineConfigs, train.LineConfig{Name: parts[0], Code: parts[1]})
		}
	}
	trains, err := train.Fetch(lineConfigs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "train error: %v\n", err)
	} else {
		data.Trains = trains
	}

	// Fetch calendar events
	creds := *googleCreds
	if creds == "" {
		creds = os.Getenv("GOOGLE_CREDENTIALS_JSON")
	}
	if creds != "" {
		ids := strings.Split(*calendarIDs, ",")
		events, err := calendar.Fetch(context.Background(), creds, ids, now)
		if err != nil {
			fmt.Fprintf(os.Stderr, "calendar error: %v\n", err)
		} else {
			data.Events = events
		}
	}

	img, err := render.RenderDashboard(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render error: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Dashboard saved to %s\n", *output)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// Don't override existing env vars
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}
