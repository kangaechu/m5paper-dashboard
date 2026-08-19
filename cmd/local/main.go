package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kangaechu/m5paper-dashboard/internal/dam"
	"github.com/kangaechu/m5paper-dashboard/internal/render"
)

func main() {
	loadEnvFile(".env")

	output := flag.String("output", "output.jpg", "output JPEG file path (white background)")
	outputDark := flag.String("output-dark", "", "output JPEG file path (dark background); derived from --output if empty")
	damURL := flag.String("dam-url", envOrDefault("DAM_URL", dam.DefaultURL), "dam data URL")
	cacheFile := flag.String("cache", envOrDefault("DAM_CACHE_FILE", "dam_history.json"), "dam history cache file path")
	flag.Parse()

	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc)

	data := render.DamDashboardData{Now: now}

	// Fetch dam data
	d, err := dam.Fetch(*damURL, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dam error: %v\n", err)
	} else {
		data.Dam = d
	}

	// Load and update history cache
	history, err := dam.LoadHistory(*cacheFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cache load error: %v\n", err)
		history = make(map[string][]render.DailyStorageRate)
	}

	if data.Dam != nil {
		dam.UpdateHistory(history, now, data.Dam.StorageRate)
		if err := dam.SaveHistory(*cacheFile, history); err != nil {
			fmt.Fprintf(os.Stderr, "cache save error: %v\n", err)
		}
	}
	data.YearlyHistory = history

	img, err := render.Dashboard(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render error: %v\n", err)
		os.Exit(1)
	}

	// Derive dark output path from --output if not specified
	darkPath := *outputDark
	if darkPath == "" {
		ext := filepath.Ext(*output)
		darkPath = strings.TrimSuffix(*output, ext) + "_dark" + ext
	}

	saveJPEG(*output, img)
	saveJPEG(darkPath, render.Invert(img))

	fmt.Printf("Dashboard saved to %s and %s\n", *output, darkPath)
}

func saveJPEG(path string, img image.Image) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 90}); err != nil {
		fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		os.Exit(1)
	}
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
