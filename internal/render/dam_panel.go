package render

import (
	"fmt"
	"math"
	"time"

	"github.com/fogleman/gg"
)

const effectiveCapacity = 28420.0 // 有効貯水容量 (×10³m³)

var weekdayJP = [...]string{"日", "月", "火", "水", "木", "金", "土"}

func drawDamHeader(dc *gg.Context, now time.Time, dam *DamData) {
	dc.SetRGB(0, 0, 0)

	// Dam name (left)
	face := fontFace(fontRegular, 22)
	dc.SetFontFace(face)
	dc.DrawStringAnchored("宇連ダム貯水率", float64(marginX), float64(headerHeight)/2, 0, 0.5)

	// Date and time (center-right)
	dateStr := fmt.Sprintf("%d年%d月%d日(%s) %s",
		now.Year(), now.Month(), now.Day(), weekdayJP[now.Weekday()],
		now.Format("15:04"))
	dc.DrawStringAnchored(dateStr, float64(Width-marginX), float64(headerHeight)/2, 1, 0.5)

	// Observation time (right, smaller)
	if dam != nil {
		face2 := fontFace(fontRegular, 14)
		dc.SetFontFace(face2)
		dc.SetRGB(0.4, 0.4, 0.4)
		obsStr := fmt.Sprintf("観測: %s", dam.ObservedAt.Format("01/02 15:04"))
		dc.DrawStringAnchored(obsStr, float64(Width/2), float64(headerHeight)/2, 0.5, 0.5)
	}
}

func drawStorageRate(dc *gg.Context, dam *DamData) {
	centerX := float64(leftWidth) / 2
	baseY := float64(mainY)

	// "現在の貯水率は" label
	dc.SetRGB(0, 0, 0)
	face := fontFace(fontRegular, 20)
	dc.SetFontFace(face)
	dc.DrawStringAnchored("現在の貯水率は", centerX, baseY+50, 0.5, 0.5)

	// Large percentage number
	rateStr := fmt.Sprintf("%.1f%%", dam.StorageRate)
	faceLarge := fontFace(fontRegular, 80)
	dc.SetFontFace(faceLarge)
	dc.DrawStringAnchored(rateStr, centerX, baseY+150, 0.5, 0.5)

	// Storage volume below
	face3 := fontFace(fontRegular, 16)
	dc.SetFontFace(face3)
	dc.SetRGB(0.3, 0.3, 0.3)
	volStr := fmt.Sprintf("(%.0f / 28,420 千m³)", dam.EffectiveStorage)
	dc.DrawStringAnchored(volStr, centerX, baseY+210, 0.5, 0.5)
}

func drawYearlyChart(dc *gg.Context, now time.Time, history map[string][]DailyStorageRate) {
	if len(history) == 0 {
		return
	}

	chartLeft := float64(rightX) + 30
	chartRight := float64(Width-marginX) - 10
	chartTop := float64(mainY) + 25
	chartBottom := float64(mainY+mainHeight) - 20
	chartWidth := chartRight - chartLeft
	chartHeight := chartBottom - chartTop

	// Title
	dc.SetRGB(0, 0, 0)
	face := fontFace(fontRegular, 14)
	dc.SetFontFace(face)
	dc.DrawStringAnchored("年間貯水率 (%)", chartLeft+chartWidth/2, float64(mainY)+12, 0.5, 0.5)

	// Draw axes
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(1)
	dc.DrawLine(chartLeft, chartBottom, chartRight, chartBottom) // X axis
	dc.DrawLine(chartLeft, chartTop, chartLeft, chartBottom)     // Y axis
	dc.Stroke()

	// Y axis labels (0%, 50%, 100%)
	faceSmall := fontFace(fontRegular, 11)
	dc.SetFontFace(faceSmall)
	dc.SetRGB(0.4, 0.4, 0.4)
	for _, pct := range []float64{0, 25, 50, 75, 100} {
		y := chartBottom - (pct/100)*chartHeight
		dc.DrawStringAnchored(fmt.Sprintf("%.0f", pct), chartLeft-5, y, 1, 0.5)
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
	// dashPatterns[0]=oldest in the 4-year window
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

		// Determine style index: years are sorted ascending, last = current year
		styleIdx := i
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

		// Build point list
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
			// Solid line
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
			// Dashed/dotted line
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
		dc.DrawStringAnchored(yearKey, legendX+25, legendY, 0, 0.5)
	}
}

type pt struct{ x, y float64 }

// drawDashedLine draws a polyline with a dash pattern (e.g. {6,4} = 6px on, 4px off).
func drawDashedLine(dc *gg.Context, points []pt, pattern []float64) {
	if len(points) < 2 || len(pattern) < 2 {
		return
	}

	drawing := true         // true=pen down, false=pen up
	patIdx := 0             // index into pattern
	remaining := pattern[0] // remaining length in current segment

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

		ux, uy := dx/segLen, dy/segLen // unit vector
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

func drawHourlyDelta(dc *gg.Context, history []DamObservation) {
	if len(history) < 2 {
		return
	}

	baseY := float64(hourlyDeltaY) + 10

	dc.SetRGB(0, 0, 0)
	face := fontFace(fontRegular, 14)
	dc.SetFontFace(face)
	dc.DrawStringAnchored("24時間の貯水量変化 (千m³)", float64(marginX), baseY+8, 0, 0.5)

	// Calculate deltas between consecutive observations
	type delta struct {
		timeStr string
		diff    float64
		storage float64
		rate    float64
	}

	// Collect deltas at 2-hour intervals
	var deltas []delta
	for i := 0; i < len(history)-1 && i < 24; i += 2 {
		end := i + 2
		if end >= len(history) {
			end = len(history) - 1
		}
		diff := history[i].EffectiveStorage - history[end].EffectiveStorage
		rate := history[i].EffectiveStorage / effectiveCapacity * 100
		deltas = append(deltas, delta{
			timeStr: history[i].Time.Format("15:04"),
			diff:    diff,
			storage: history[i].EffectiveStorage,
			rate:    rate,
		})
	}

	if len(deltas) == 0 {
		return
	}

	// Draw as a table with columns
	colCount := len(deltas)

	availWidth := float64(contentWidth)
	colWidth := availWidth / float64(colCount)

	faceLabel := fontFace(fontRegular, 16)

	for i := 0; i < colCount; i++ {
		d := deltas[i]
		x := float64(marginX) + float64(i)*colWidth + colWidth/2

		// Time label
		dc.SetRGB(0.3, 0.3, 0.3)
		dc.SetFontFace(faceLabel)
		dc.DrawStringAnchored(d.timeStr, x, baseY+38, 0.5, 0.5)

		// Delta value with color coding
		var diffStr string
		if d.diff >= 0 {
			diffStr = fmt.Sprintf("+%.0f", d.diff)
			dc.SetRGB(0, 0, 0) // black for positive
		} else {
			diffStr = fmt.Sprintf("%.0f", d.diff)
			dc.SetRGB(0.5, 0.5, 0.5) // gray for negative
		}

		faceVal := fontFace(fontRegular, 20)
		dc.SetFontFace(faceVal)
		dc.DrawStringAnchored(diffStr, x, baseY+65, 0.5, 0.5)

		// Storage value
		dc.SetRGB(0.5, 0.5, 0.5)
		dc.SetFontFace(faceLabel)
		dc.DrawStringAnchored(fmt.Sprintf("%.0f", d.storage), x, baseY+90, 0.5, 0.5)

		// Storage rate
		dc.SetRGB(0.3, 0.3, 0.3)
		dc.SetFontFace(faceLabel)
		dc.DrawStringAnchored(fmt.Sprintf("%.1f%%", d.rate), x, baseY+115, 0.5, 0.5)
	}

	// Column separators
	dc.SetRGB(0.9, 0.9, 0.9)
	dc.SetLineWidth(0.5)
	for i := 1; i < colCount; i++ {
		x := float64(marginX) + float64(i)*colWidth
		dc.DrawLine(x, baseY+25, x, baseY+125)
	}
	dc.Stroke()

	// Row labels
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.SetFontFace(faceLabel)
	dc.DrawStringAnchored("時刻", float64(marginX)-2, baseY+38, 1, 0.5)
	dc.DrawStringAnchored("差異", float64(marginX)-2, baseY+65, 1, 0.5)
	dc.DrawStringAnchored("貯水量", float64(marginX)-2, baseY+90, 1, 0.5)
	dc.DrawStringAnchored("貯水率", float64(marginX)-2, baseY+115, 1, 0.5)
}

func drawDamFooter(dc *gg.Context, dam *DamData) {
	y := float64(footerY) + float64(footerHeight)/2

	dc.SetRGB(0, 0, 0)
	face := fontFace(fontRegular, 16)
	dc.SetFontFace(face)

	stats := fmt.Sprintf("貯水位: %.2fm    貯水量: %.0f千m³    流入量: %.2fm³/s    放流量: %.2fm³/s",
		dam.WaterLevel, dam.EffectiveStorage, dam.Inflow, dam.Outflow)
	dc.DrawStringAnchored(stats, float64(Width)/2, y, 0.5, 0.5)
}
