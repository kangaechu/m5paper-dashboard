package render

import (
	"image"
	"image/color"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/kangaechu/m5paper-dashboard/fonts"
)

// DamDashboardData holds all data needed to render the dam dashboard.
type DamDashboardData struct {
	Now           time.Time
	Dam           *DamData
	YearlyHistory map[string][]DailyStorageRate // key: "2026", "2025", ...
}

// DamData holds current dam observation data.
type DamData struct {
	Name             string
	ObservedAt       time.Time
	WaterLevel       float64 // 貯水位 (EL.m)
	EffectiveStorage float64 // 有効貯水量 (×10³m³)
	StorageRate      float64 // 貯水率 (%)
	Inflow           float64 // 流入量 (m³/s)
	Outflow          float64 // 放流量 (m³/s)
	Rainfall         float64 // ダム地点雨量 (mm/h)
	History          []DamObservation
}

// DamObservation holds one hourly observation.
type DamObservation struct {
	Time             time.Time
	WaterLevel       float64
	EffectiveStorage float64
	Inflow           float64
	Outflow          float64
}

// DailyStorageRate holds one day's storage rate for the yearly chart.
type DailyStorageRate struct {
	Date        string  `json:"date"`         // "2026-01-15"
	StorageRate float64 `json:"storage_rate"` // percentage
}

var fontRegular *truetype.Font

func init() {
	var err error
	fontRegular, err = loadFont("NotoSansJP-SemiBold.ttf")
	if err != nil {
		panic("failed to load NotoSansJP: " + err.Error())
	}
}

func loadFont(name string) (*truetype.Font, error) {
	data, err := fonts.FS.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return truetype.Parse(data)
}

func fontFace(f *truetype.Font, size float64) font.Face {
	return truetype.NewFace(f, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// Dashboard generates the dam dashboard image.
func Dashboard(data DamDashboardData) (image.Image, error) {
	dc := gg.NewContext(Width, Height)

	// White background
	dc.SetColor(color.White)
	dc.Clear()

	drawDamHeader(dc, data.Now, data.Dam)
	drawSeparator(dc, float64(mainY))

	if data.Dam != nil {
		drawStorageRate(dc, data.Dam)
		drawHourlyDelta(dc, data.Dam.History)
		drawDamFooter(dc, data.Dam)
	}

	drawYearlyChart(dc, data.Now, data.YearlyHistory)

	drawSeparator(dc, float64(hourlyDeltaY))
	drawSeparator(dc, float64(footerY))

	return toGrayscale(dc.Image()), nil
}

func drawSeparator(dc *gg.Context, y float64) {
	dc.SetRGB(separatorGray, separatorGray, separatorGray)
	dc.SetLineWidth(1)
	dc.DrawLine(float64(marginX), y, float64(Width-marginX), y)
	dc.Stroke()
}

// toGrayscale converts src to a grayscale image stored as NRGBA.
// Using NRGBA (not image.Gray) so jpeg.Encode produces RGB JPEG,
// which is required by M5EPD's TJpgDec decoder.
func toGrayscale(src image.Image) *image.NRGBA {
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			lum := uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
			// Quantize to 4-bit (16 shades) for e-ink
			lum = (lum / 17) * 17
			i := dst.PixOffset(x, y)
			dst.Pix[i+0] = lum
			dst.Pix[i+1] = lum
			dst.Pix[i+2] = lum
			dst.Pix[i+3] = 0xFF
		}
	}
	return dst
}
