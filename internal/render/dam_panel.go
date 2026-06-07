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

// drawStorageRateOverlay draws the total storage rate as a compact card
// overlaid on the bottom-right corner of the chart. The card has an opaque
// white background so the figures stay legible over the graph.
func drawStorageRateOverlay(dc *gg.Context, dam *DamData) {
	rateStr := fmt.Sprintf("%.0f%%", dam.StorageRate)
	volStr := fmt.Sprintf("%s / %s 万m³",
		formatThousands(dam.Total.Storage),
		formatThousands(dam.Total.EffectiveCapacity))

	faceRate := fontFace(fontRegular, 48)
	faceVol := fontFace(fontRegular, 12)

	dc.SetFontFace(faceRate)
	rateW, _ := dc.MeasureString(rateStr)
	dc.SetFontFace(faceVol)
	volW, _ := dc.MeasureString(volStr)

	rateH := float64(faceRate.Metrics().Height) / 64
	volH := float64(faceVol.Metrics().Height) / 64

	const (
		pad = 12.0
		gap = 2.0
	)
	boxW := math.Max(rateW, volW) + pad*2
	boxH := rateH + volH + gap + pad*2

	boxRight := float64(Width - marginX)
	boxBottom := float64(Height - marginX)
	boxLeft := boxRight - boxW
	boxTop := boxBottom - boxH
	centerX := boxLeft + boxW/2

	// Opaque white card with a light border.
	dc.SetRGB(1, 1, 1)
	dc.DrawRoundedRectangle(boxLeft, boxTop, boxW, boxH, 10)
	dc.Fill()
	dc.SetRGB(separatorGray, separatorGray, separatorGray)
	dc.SetLineWidth(1.5)
	dc.DrawRoundedRectangle(boxLeft, boxTop, boxW, boxH, 10)
	dc.Stroke()

	// Large percentage number
	dc.SetRGB(0, 0, 0)
	dc.SetFontFace(faceRate)
	rateCY := boxTop + pad + rateH/2
	dc.DrawStringAnchored(rateStr, centerX, centeredBaselineY(faceRate, rateCY), 0.5, 0)

	// Storage volume
	dc.SetRGB(0.3, 0.3, 0.3)
	dc.SetFontFace(faceVol)
	volCY := boxTop + pad + rateH + gap + volH/2
	dc.DrawStringAnchored(volStr, centerX, centeredBaselineY(faceVol, volCY), 0.5, 0)
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
// source page) across the full main area of the dashboard, scaled to fit
// while preserving aspect ratio, with a small caption underneath.
func drawGraphImage(dc *gg.Context, img image.Image) {
	areaLeft := float64(marginX)
	areaRight := float64(Width - marginX)
	areaTop := float64(mainY) + 4
	areaBottom := float64(Height) - 4

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
}
