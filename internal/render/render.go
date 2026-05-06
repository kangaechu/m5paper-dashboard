package render

import (
	"image"
	"image/color"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/kangaechu/m5paper-dashboard/fonts"
)

// DamDashboardData holds all data needed to render the dam dashboard.
type DamDashboardData struct {
	Now           time.Time
	Dam           *DamData
	YearlyHistory map[string][]DailyStorageRate // key: "2026", "2025", ...
}

// DamData holds current dam observation data for the entire river system.
type DamData struct {
	SystemName  string         // e.g. "荒川水系"
	ObservedAt  time.Time      // observation timestamp from the source page
	Total       DamReservoir   // aggregated total of all reservoirs
	Reservoirs  []DamReservoir // individual reservoirs (display order)
	StorageRate float64        // shortcut for Total.StorageRate
}

// DamReservoir holds the headline figures for one dam (or the system total).
type DamReservoir struct {
	Name              string
	EffectiveCapacity float64 // 有効容量 (万m³)
	Storage           float64 // 貯水量 (万m³)
	StorageRate       float64 // 貯水率 (%)
}

// DailyStorageRate holds one day's storage rate for the yearly chart.
type DailyStorageRate struct {
	Date        string  `json:"date"`         // "2026-01-15"
	StorageRate float64 `json:"storage_rate"` // percentage
}

var fontRegular *opentype.Font

func init() {
	var err error
	fontRegular, err = loadFont("NotoSansJP-SemiBold.ttf")
	if err != nil {
		panic("failed to load NotoSansJP: " + err.Error())
	}
}

func loadFont(name string) (*opentype.Font, error) {
	data, err := fonts.FS.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return opentype.Parse(data)
}

func fontFace(f *opentype.Font, size float64) font.Face {
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic("failed to create font face: " + err.Error())
	}
	return face
}

// centeredBaselineY returns the baseline Y that visually centers text
// of the given font.Face around centerY.
func centeredBaselineY(face font.Face, centerY float64) float64 {
	m := face.Metrics()
	ascent := float64(m.Ascent) / 64.0
	descent := float64(m.Descent) / 64.0
	return centerY + (ascent-descent)/2
}

// Dashboard generates the dam dashboard image.
func Dashboard(data DamDashboardData) (*image.NRGBA, error) {
	dc := gg.NewContext(Width, Height)

	// White background
	dc.SetColor(color.White)
	dc.Clear()

	drawDamHeader(dc, data.Now, data.Dam)
	drawSeparator(dc, float64(mainY))

	if data.Dam != nil {
		drawStorageRate(dc, data.Dam)
	}

	drawYearlyChart(dc, data.Now, data.YearlyHistory)

	return toGrayscale(dc.Image()), nil
}

// Invert returns a new image with all pixel values inverted (255 - v).
func Invert(src *image.NRGBA) *image.NRGBA {
	dst := image.NewNRGBA(src.Bounds())
	for i := 0; i < len(src.Pix); i += 4 {
		dst.Pix[i+0] = 255 - src.Pix[i+0] // R
		dst.Pix[i+1] = 255 - src.Pix[i+1] // G
		dst.Pix[i+2] = 255 - src.Pix[i+2] // B
		dst.Pix[i+3] = src.Pix[i+3]       // A
	}
	return dst
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
