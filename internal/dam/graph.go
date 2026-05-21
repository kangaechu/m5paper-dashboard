package dam

import (
	"fmt"
	"image"
	"image/png"
	"net/http"
	"time"
)

// DefaultGraphURL is the Tokyo Waterworks Bureau page that returns a PNG
// chart of the Arakawa river system dam storage (current year, last year,
// 平年, and a notable past year), updated daily.
const DefaultGraphURL = "https://www.waterworks.metro.tokyo.lg.jp/documents/d/waterworks/suigen_g_arakawa"

// FetchGraph downloads and decodes the storage-rate chart PNG.
func FetchGraph(url string) (image.Image, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch graph: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("graph returned %d", resp.StatusCode)
	}

	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("decode graph PNG: %w", err)
	}
	return img, nil
}
