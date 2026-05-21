package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	graphURL := envOrDefault("DAM_GRAPH_URL", dam.DefaultGraphURL)

	// Fetch dam data
	d, err := dam.Fetch(damURL, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dam error: %v\n", err)
	} else {
		data.Dam = d
	}

	// Fetch storage chart image
	g, err := dam.FetchGraph(graphURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "graph error: %v\n", err)
	} else {
		data.GraphImage = g
	}

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

	// Upload dark version
	darkKey := envOrDefault("S3_OBJECT_KEY_DARK", "")
	if darkKey == "" {
		ext := filepath.Ext(key)
		darkKey = strings.TrimSuffix(key, ext) + "_dark" + ext
	}
	if err := uploader.Upload(ctx, bucket, darkKey, render.Invert(img)); err != nil {
		return fmt.Errorf("upload dark: %w", err)
	}
	fmt.Printf("Dashboard (dark) uploaded to s3://%s/%s\n", bucket, darkKey)
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
