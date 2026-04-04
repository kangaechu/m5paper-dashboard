package render

import (
	"fmt"

	"github.com/fogleman/gg"
)

// Weather Icon codepoints (from Weather Icons font)
var weatherIconMap = map[int]rune{
	0:  '\uf00d', // Clear sky → wi-day-sunny
	1:  '\uf00c', // Mainly clear → wi-day-cloudy
	2:  '\uf002', // Partly cloudy → wi-cloudy
	3:  '\uf013', // Overcast → wi-cloud
	45: '\uf014', // Fog
	48: '\uf014', // Rime fog
	51: '\uf01a', // Light drizzle
	53: '\uf01a', // Moderate drizzle
	55: '\uf01a', // Dense drizzle
	61: '\uf019', // Slight rain → wi-rain
	63: '\uf019', // Moderate rain
	65: '\uf019', // Heavy rain
	71: '\uf01b', // Slight snow → wi-snow
	73: '\uf01b', // Moderate snow
	75: '\uf01b', // Heavy snow
	80: '\uf009', // Slight rain showers → wi-showers
	81: '\uf009', // Moderate rain showers
	82: '\uf009', // Violent rain showers
	95: '\uf01e', // Thunderstorm → wi-thunderstorm
	96: '\uf01e', // Thunderstorm with hail
	99: '\uf01e', // Thunderstorm with heavy hail
}

func weatherIcon(code int) string {
	if r, ok := weatherIconMap[code]; ok {
		return string(r)
	}
	return string('\uf07b') // wi-na (not available)
}

func drawWeather(dc *gg.Context, w *WeatherData) {
	baseY := float64(weatherY)
	divX := float64(Width) * 0.55 // left 55%, right 45%

	// === Left: Today's weather ===
	leftCX := float64(marginX) + (divX-float64(marginX))/2

	// Weather icon (large)
	iconFace := fontFace(fontWeather, 64)
	dc.SetFontFace(iconFace)
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(weatherIcon(w.WeatherCode), leftCX-60, baseY+55, 0.5, 0.5)

	// Description
	descFace := fontFace(fontRegular, 26)
	dc.SetFontFace(descFace)
	dc.DrawStringAnchored(w.Description, leftCX+50, baseY+35, 0.5, 0.5)

	// Temperature
	tempFace := fontFace(fontRegular, 42)
	dc.SetFontFace(tempFace)
	dc.DrawStringAnchored(fmt.Sprintf("%.0f℃", w.Temperature), leftCX+50, baseY+80, 0.5, 0.5)

	// Detail info
	smallFace := fontFace(fontRegular, 16)
	dc.SetFontFace(smallFace)
	dc.SetRGB(0.3, 0.3, 0.3)
	infoY := baseY + 130
	if w.Humidity > 0 {
		dc.DrawStringAnchored(fmt.Sprintf("湿度 %d%%", w.Humidity), leftCX, infoY, 0.5, 0.5)
		infoY += 25
	}
	dc.DrawStringAnchored(fmt.Sprintf("降水確率 %d%%", w.PrecipChance), leftCX, infoY, 0.5, 0.5)
	infoY += 25
	if len(w.Forecasts) > 0 {
		f := w.Forecasts[0]
		if f.TempMin != 0 || f.TempMax != 0 {
			dc.DrawStringAnchored(fmt.Sprintf("%.0f℃ / %.0f℃", f.TempMin, f.TempMax), leftCX, infoY, 0.5, 0.5)
		}
	}

	// === Right: Tomorrow + Day after tomorrow ===
	halfH := float64(weatherHeight) / 2

	textLeft := divX + 75 // left-align text to the right of icon

	for i := 1; i < len(w.Forecasts) && i <= 2; i++ {
		f := w.Forecasts[i]
		cy := baseY + halfH*float64(i-1) + halfH/2

		// Day label (left-aligned)
		dc.SetRGB(0, 0, 0)
		labelFace := fontFace(fontRegular, 18)
		dc.SetFontFace(labelFace)
		dc.DrawString(f.DayLabel, textLeft, cy-22)

		// Icon
		iconSmall := fontFace(fontWeather, 34)
		dc.SetFontFace(iconSmall)
		dc.DrawStringAnchored(weatherIcon(f.WeatherCode), divX+40, cy+8, 0.5, 0.5)

		// Description (left-aligned, same baseline)
		descSmall := fontFace(fontRegular, 16)
		dc.SetFontFace(descSmall)
		dc.DrawString(f.Description, textLeft, cy+4)

		// Temp + precipitation (left-aligned)
		dc.SetRGB(0.3, 0.3, 0.3)
		dc.SetFontFace(descSmall)
		dc.DrawString(
			fmt.Sprintf("%.0f/%.0f℃  %d%%", f.TempMin, f.TempMax, f.PrecipChance),
			textLeft, cy+24)
		dc.SetRGB(0, 0, 0)
	}
}
