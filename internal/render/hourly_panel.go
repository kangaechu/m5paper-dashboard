package render

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

// windDirectionArrow returns a Unicode arrow for wind direction (degrees).
var windDirLabels = [...]string{"↓", "↙", "←", "↖", "↑", "↗", "→", "↘"}

func windArrow(deg int) string {
	// Wind direction is where wind comes FROM, arrow shows where it blows TO
	// 0°=N wind (from north, blows south) → ↓
	idx := int(math.Round(float64(deg)/45.0)) % 8
	return windDirLabels[idx]
}

func drawHourly(dc *gg.Context, hourly []HourlyWeather) {
	if len(hourly) == 0 {
		return
	}

	baseY := float64(hourlyY)

	// Show 12 entries (2-hour intervals)
	count := len(hourly)
	if count > 12 {
		count = 12
	}

	cellW := float64(contentWidth) / float64(count)

	timeFace := fontFace(fontRegular, 12)
	iconFace := fontFace(fontWeather, 18)
	tempFace := fontFace(fontRegular, 13)
	windFace := fontFace(fontRegular, 11)

	for i := 0; i < count; i++ {
		h := hourly[i]
		cx := float64(marginX) + cellW*float64(i) + cellW/2

		// Hour
		dc.SetRGB(0.3, 0.3, 0.3)
		dc.SetFontFace(timeFace)
		dc.DrawStringAnchored(fmt.Sprintf("%d時", h.Hour), cx, baseY+15, 0.5, 0.5)

		// Weather icon
		dc.SetRGB(0, 0, 0)
		dc.SetFontFace(iconFace)
		dc.DrawStringAnchored(weatherIcon(h.WeatherCode), cx, baseY+35, 0.5, 0.5)

		// Temperature
		dc.SetFontFace(tempFace)
		dc.DrawStringAnchored(fmt.Sprintf("%.0f°", h.Temperature), cx, baseY+58, 0.5, 0.5)

		// Precipitation probability
		if h.PrecipProb > 0 {
			dc.SetRGB(0.4, 0.4, 0.4)
			dc.SetFontFace(timeFace)
			dc.DrawStringAnchored(fmt.Sprintf("%d%%", h.PrecipProb), cx, baseY+73, 0.5, 0.5)
		}

		// Wind: arrow + speed
		dc.SetRGB(0.3, 0.3, 0.3)
		dc.SetFontFace(windFace)
		windStr := fmt.Sprintf("%s%.0f", windArrow(h.WindDirection), h.WindSpeed)
		dc.DrawStringAnchored(windStr, cx, baseY+80, 0.5, 0.5)
	}
}
