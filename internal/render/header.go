package render

import (
	"fmt"
	"time"

	"github.com/fogleman/gg"
)

var weekdayJP = [...]string{"日", "月", "火", "水", "木", "金", "土"}

func drawHeader(dc *gg.Context, now time.Time) {
	dc.SetRGB(0, 0, 0)

	// Date: 2026年4月4日(土)
	face := fontFace(fontRegular, 28)
	dc.SetFontFace(face)
	dateStr := fmt.Sprintf("%d年%d月%d日(%s)",
		now.Year(), now.Month(), now.Day(), weekdayJP[now.Weekday()])
	dc.DrawStringAnchored(dateStr, float64(marginX), float64(headerHeight)/2, 0, 0.5)

	// Time: 15:30
	timeStr := now.Format("15:04")
	dc.DrawStringAnchored(timeStr, float64(Width-marginX), float64(headerHeight)/2, 1, 0.5)
}
