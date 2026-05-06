package render

import (
	"fmt"
	"math"
	"time"

	"github.com/fogleman/gg"
)

var weekdayJP = [...]string{"日", "月", "火", "水", "木", "金", "土"}

func drawDamHeader(dc *gg.Context, now time.Time, dam *DamData) {
	dc.SetRGB(0, 0, 0)

	face := fontFace(fontRegular, 22)
	dc.SetFontFace(face)
	baselineY := centeredBaselineY(face, float64(headerHeight)/2)

	title := "荒川水系 貯水率"
	if dam != nil && dam.SystemName != "" {
		title = dam.SystemName + " 貯水率"
	}
	dc.DrawStringAnchored(title, float64(marginX), baselineY, 0, 0)

	// Right: current rendering time
	dateStr := fmt.Sprintf("%d年%d月%d日(%s) %s",
		now.Year(), now.Month(), now.Day(), weekdayJP[now.Weekday()],
		now.Format("15:04"))
	dc.DrawStringAnchored(dateStr, float64(Width-marginX), baselineY, 1, 0)

	// Center: observation time from source page
	if dam != nil {
		face2 := fontFace(fontRegular, 14)
		dc.SetFontFace(face2)
		dc.SetRGB(0.4, 0.4, 0.4)
		obsStr := fmt.Sprintf("観測: %s 現在", dam.ObservedAt.Format("01/02 15:04"))
		dc.DrawStringAnchored(obsStr, float64(Width/2), baselineY, 0.5, 0)
	}
}

func drawStorageRate(dc *gg.Context, dam *DamData) {
	centerX := float64(leftWidth) / 2
	baseY := float64(mainY)

	// "4ダム合計貯水率" label
	dc.SetRGB(0, 0, 0)
	face := fontFace(fontRegular, 20)
	dc.SetFontFace(face)
	dc.DrawStringAnchored("4ダム合計貯水率", centerX, centeredBaselineY(face, baseY+50), 0.5, 0)

	// Large percentage number
	rateStr := fmt.Sprintf("%.0f%%", dam.StorageRate)
	faceLarge := fontFace(fontRegular, 100)
	dc.SetFontFace(faceLarge)
	dc.DrawStringAnchored(rateStr, centerX, centeredBaselineY(faceLarge, baseY+150), 0.5, 0)

	// Storage volume below
	face3 := fontFace(fontRegular, 14)
	dc.SetFontFace(face3)
	dc.SetRGB(0.3, 0.3, 0.3)
	volStr := fmt.Sprintf("%s / %s 万m³",
		formatThousands(dam.Total.Storage),
		formatThousands(dam.Total.EffectiveCapacity))
	dc.DrawStringAnchored(volStr, centerX, centeredBaselineY(face3, baseY+220), 0.5, 0)

	// Per-dam mini list
	if len(dam.Reservoirs) > 0 {
		listFace := fontFace(fontRegular, 14)
		dc.SetFontFace(listFace)
		dc.SetRGB(0.2, 0.2, 0.2)

		lineHeight := 22.0
		listTop := baseY + 260
		nameX := float64(marginX) + 20
		rateX := float64(leftWidth) - 20

		for i, r := range dam.Reservoirs {
			y := listTop + float64(i)*lineHeight
			bl := centeredBaselineY(listFace, y)
			dc.DrawStringAnchored(r.Name, nameX, bl, 0, 0)
			dc.DrawStringAnchored(fmt.Sprintf("%.0f%%", r.StorageRate), rateX, bl, 1, 0)
		}
	}
}

// formatThousands renders an integer-valued float with comma thousands separators.
func formatThousands(v float64) string {
	n := int64(math.Round(v))
	neg := n < 0
	if neg {
		n = -n
	}
	s := fmt.Sprintf("%d", n)
	// Insert commas
	var out []byte
	for i, c := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, c)
	}
	if neg {
		return "-" + string(out)
	}
	return string(out)
}

func drawYearlyChart(dc *gg.Context, now time.Time, history map[string][]DailyStorageRate) {
	if len(history) == 0 {
		return
	}

	chartLeft := float64(rightX) + 30
	chartRight := float64(Width-marginX) - 10
	chartTop := float64(mainY) + 25
	chartBottom := float64(mainY+mainHeight) - 30
	chartWidth := chartRight - chartLeft
	chartHeight := chartBottom - chartTop

	// Title
	dc.SetRGB(0, 0, 0)
	face := fontFace(fontRegular, 14)
	dc.SetFontFace(face)
	dc.DrawStringAnchored("年間貯水率 (%)", chartLeft+chartWidth/2, centeredBaselineY(face, float64(mainY)+12), 0.5, 0)

	// Draw axes
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(1)
	dc.DrawLine(chartLeft, chartBottom, chartRight, chartBottom) // X axis
	dc.DrawLine(chartLeft, chartTop, chartLeft, chartBottom)     // Y axis
	dc.Stroke()

	// Y axis labels (0%, 25%, 50%, 75%, 100%)
	faceSmall := fontFace(fontRegular, 11)
	dc.SetFontFace(faceSmall)
	dc.SetRGB(0.4, 0.4, 0.4)
	for _, pct := range []float64{0, 25, 50, 75, 100} {
		y := chartBottom - (pct/100)*chartHeight
		dc.DrawStringAnchored(fmt.Sprintf("%.0f", pct), chartLeft-5, centeredBaselineY(faceSmall, y), 1, 0)
		// Grid line
		if pct > 0 && pct < 100 {
			dc.SetRGB(0.85, 0.85, 0.85)
			dc.SetLineWidth(0.5)
			dc.DrawLine(chartLeft, y, chartRight, y)
			dc.Stroke()
			dc.SetRGB(0.4, 0.4, 0.4)
		}
	}

	// X axis labels (month names)
	for m := 1; m <= 12; m++ {
		x := chartLeft + (float64(m)-0.5)/12*chartWidth
		dc.DrawStringAnchored(fmt.Sprintf("%d月", m), x, chartBottom+14, 0.5, 0)
	}

	// Show only current year + 3 previous years
	currentYear := now.Format("2006")
	currentYearInt := now.Year()
	years := make([]string, 0, 4)
	for y := currentYearInt - 3; y <= currentYearInt; y++ {
		key := fmt.Sprintf("%d", y)
		if _, ok := history[key]; ok {
			years = append(years, key)
		}
	}

	// Line styles per year: current=solid thick, previous years=dashed/dotted
	dashPatterns := [][]float64{
		{2, 4}, // dotted (3 years ago)
		{6, 4}, // short dash (2 years ago)
		nil,    // solid (last year)
		nil,    // solid (current year)
	}

	for i, yearKey := range years {
		entries := history[yearKey]
		if len(entries) == 0 {
			continue
		}

		// Anchor styles to the right of the pattern list so the current
		// year (always last in `years`) maps to the solid pattern even
		// when fewer than 4 years are available.
		styleIdx := len(dashPatterns) - (len(years) - i)
		if styleIdx < 0 {
			styleIdx = 0
		}
		if styleIdx >= len(dashPatterns) {
			styleIdx = len(dashPatterns) - 1
		}

		var lineWidth float64
		if yearKey == currentYear {
			dc.SetRGB(0, 0, 0)
			lineWidth = 2
		} else {
			gray := 0.3 + float64(len(years)-1-i)*0.15
			if gray > 0.7 {
				gray = 0.7
			}
			dc.SetRGB(gray, gray, gray)
			lineWidth = 1.5
		}

		var points []pt
		for _, e := range entries {
			t, err := time.Parse("2006-01-02", e.Date)
			if err != nil {
				continue
			}
			dayOfYear := t.YearDay()
			daysInYear := 365
			if isLeapYear(t.Year()) {
				daysInYear = 366
			}
			x := chartLeft + float64(dayOfYear)/float64(daysInYear)*chartWidth
			y := chartBottom - (e.StorageRate/100)*chartHeight
			y = math.Max(chartTop, math.Min(chartBottom, y))
			points = append(points, pt{x, y})
		}

		dash := dashPatterns[styleIdx]
		if dash == nil {
			dc.SetLineWidth(lineWidth)
			for j, p := range points {
				if j == 0 {
					dc.MoveTo(p.x, p.y)
				} else {
					dc.LineTo(p.x, p.y)
				}
			}
			dc.Stroke()
		} else {
			dc.SetLineWidth(lineWidth)
			drawDashedLine(dc, points, dash)
		}

		// Legend
		legendY := float64(mainY) + 25 + float64(i)*16
		legendX := chartRight - 60
		if dash == nil {
			dc.DrawLine(legendX, legendY, legendX+20, legendY)
		} else {
			drawDashedLine(dc, []pt{{legendX, legendY}, {legendX + 20, legendY}}, dash)
		}
		dc.Stroke()
		dc.SetFontFace(faceSmall)
		dc.DrawStringAnchored(yearKey, legendX+25, centeredBaselineY(faceSmall, legendY), 0, 0)
	}

	// Average line (H22-H30 平年値). Stored under the "average" key.
	if avg, ok := history["average"]; ok && len(avg) > 0 {
		avgDash := []float64{4, 3}
		dc.SetRGB(0.5, 0.5, 0.5)
		dc.SetLineWidth(1.5)

		var points []pt
		for _, e := range avg {
			t, err := time.Parse("2006-01-02", e.Date)
			if err != nil {
				continue
			}
			dayOfYear := t.YearDay()
			x := chartLeft + float64(dayOfYear)/365.0*chartWidth
			y := chartBottom - (e.StorageRate/100)*chartHeight
			y = math.Max(chartTop, math.Min(chartBottom, y))
			points = append(points, pt{x, y})
		}
		drawDashedLine(dc, points, avgDash)

		// Legend (placed below the per-year legends)
		legendY := float64(mainY) + 25 + float64(len(years))*16
		legendX := chartRight - 60
		drawDashedLine(dc, []pt{{legendX, legendY}, {legendX + 20, legendY}}, avgDash)
		dc.Stroke()
		dc.SetFontFace(faceSmall)
		dc.DrawStringAnchored("平年", legendX+25, centeredBaselineY(faceSmall, legendY), 0, 0)
	}
}

type pt struct{ x, y float64 }

// drawDashedLine draws a polyline with a dash pattern (e.g. {6,4} = 6px on, 4px off).
func drawDashedLine(dc *gg.Context, points []pt, pattern []float64) {
	if len(points) < 2 || len(pattern) < 2 {
		return
	}

	drawing := true
	patIdx := 0
	remaining := pattern[0]

	dc.MoveTo(points[0].x, points[0].y)
	prev := points[0]

	for i := 1; i < len(points); i++ {
		cur := points[i]
		dx := cur.x - prev.x
		dy := cur.y - prev.y
		segLen := math.Hypot(dx, dy)
		if segLen == 0 {
			continue
		}

		ux, uy := dx/segLen, dy/segLen
		consumed := 0.0

		for consumed < segLen {
			step := math.Min(remaining, segLen-consumed)
			nx := prev.x + ux*(consumed+step)
			ny := prev.y + uy*(consumed+step)

			if drawing {
				dc.LineTo(nx, ny)
			} else {
				dc.MoveTo(nx, ny)
			}

			consumed += step
			remaining -= step

			if remaining <= 0 {
				if drawing {
					dc.Stroke()
				}
				patIdx = (patIdx + 1) % len(pattern)
				drawing = !drawing
				remaining = pattern[patIdx]
				if !drawing {
					dc.MoveTo(nx, ny)
				}
			}
		}
		prev = cur
	}

	if drawing {
		dc.Stroke()
	}
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
