package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kangaechu/m5paper-dashboard/internal/dam"
	"github.com/kangaechu/m5paper-dashboard/internal/render"
	"github.com/kangaechu/m5paper-dashboard/internal/uploader"
)

func handler(ctx context.Context) error {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc)

	data := render.DamDashboardData{Now: now}

	damURL := envOrDefault("DAM_URL", dam.DefaultURL)
	cacheFile := envOrDefault("DAM_CACHE_FILE", "/tmp/dam_history.json")

	// Fetch dam data
	d, err := dam.Fetch(damURL, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dam error: %v\n", err)
	} else {
		data.Dam = d
	}

	// Load and update history cache
	history, err := dam.LoadHistory(cacheFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cache load error: %v\n", err)
		history = make(map[string][]render.DailyStorageRate)
	}

	if data.Dam != nil {
		dam.UpdateHistory(history, now, data.Dam.StorageRate)
		if err := dam.SaveHistory(cacheFile, history); err != nil {
			fmt.Fprintf(os.Stderr, "cache save error: %v\n", err)
		}
	}
	data.YearlyHistory = history

	// Render
	img, err := render.Dashboard(data)
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
