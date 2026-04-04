package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kangaechu/m5paper-dashboard/internal/calendar"
	"github.com/kangaechu/m5paper-dashboard/internal/render"
	"github.com/kangaechu/m5paper-dashboard/internal/train"
	"github.com/kangaechu/m5paper-dashboard/internal/uploader"
	"github.com/kangaechu/m5paper-dashboard/internal/weather"
)

func handler(ctx context.Context) error {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc)

	data := render.DashboardData{Now: now}

	locationCode := envOrDefault("LOCATION_CODE", "130000")
	lat := envOrDefault("LOCATION_LAT", "35.6895")
	lon := envOrDefault("LOCATION_LON", "139.6917")

	// Fetch weather
	w, err := weather.Fetch(locationCode, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "weather error: %v\n", err)
	} else {
		data.Weather = w
	}

	// Fetch hourly weather
	hourly, err := weather.FetchHourly(lat, lon, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hourly weather error: %v\n", err)
	} else if data.Weather != nil {
		data.Weather.Hourly = hourly
	}

	// Fetch train delay info
	trainLinesStr := envOrDefault("TRAIN_LINES", "")
	if trainLinesStr != "" {
		var lineConfigs []train.LineConfig
		for _, entry := range strings.Split(trainLinesStr, ",") {
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
	}

	// Fetch calendar events
	creds := os.Getenv("GOOGLE_CREDENTIALS_JSON")
	if creds != "" {
		calIDs := strings.Split(envOrDefault("GOOGLE_CALENDAR_IDS", "primary"), ",")
		events, err := calendar.Fetch(ctx, creds, calIDs, now)
		if err != nil {
			fmt.Fprintf(os.Stderr, "calendar error: %v\n", err)
		} else {
			data.Events = events
		}
	}

	// Render
	img, err := render.RenderDashboard(data)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	// Upload to S3
	bucket := os.Getenv("S3_BUCKET")
	key := os.Getenv("S3_OBJECT_KEY")
	if bucket == "" || key == "" {
		return fmt.Errorf("S3_BUCKET and S3_OBJECT_KEY must be set")
	}

	if err := uploader.Upload(ctx, bucket, key, img); err != nil {
		return fmt.Errorf("upload: %w", err)
	}

	fmt.Printf("Dashboard uploaded to s3://%s/%s\n", bucket, key)
	return nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	lambda.Start(handler)
}
