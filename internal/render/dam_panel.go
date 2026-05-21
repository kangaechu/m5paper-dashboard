package render

import (
	"fmt"
	"image"
	"math"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
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

// drawGraphImage renders the yearly storage chart PNG (fetched from the
// source page) into the right-hand area of the dashboard, scaled to fit
// while preserving aspect ratio, with a small caption underneath.
func drawGraphImage(dc *gg.Context, img image.Image) {
	areaLeft := float64(rightX) + 6
	areaRight := float64(Width - marginX)
	areaTop := float64(mainY) + 4
	areaBottom := float64(Height) - 18

	areaW := areaRight - areaLeft
	areaH := areaBottom - areaTop

	srcB := img.Bounds()
	srcW := float64(srcB.Dx())
	srcH := float64(srcB.Dy())
	if srcW <= 0 || srcH <= 0 {
		return
	}

	scale := math.Min(areaW/srcW, areaH/srcH)
	dstW := int(math.Round(srcW * scale))
	dstH := int(math.Round(srcH * scale))
	dstX := int(math.Round(areaLeft + (areaW-float64(dstW))/2))
	dstY := int(math.Round(areaTop))

	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, srcB, draw.Over, nil)
	dc.DrawImage(dst, dstX, dstY)

	// Source attribution caption
	face := fontFace(fontRegular, 11)
	dc.SetFontFace(face)
	dc.SetRGB(0.4, 0.4, 0.4)
	dc.DrawStringAnchored("出典: 東京都水道局", areaRight, centeredBaselineY(face, areaBottom+9), 1, 0)
}
