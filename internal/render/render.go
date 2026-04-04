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

// DashboardData holds all data needed to render the dashboard.
type DashboardData struct {
	Now     time.Time
	Weather *WeatherData
	Trains  []TrainInfo
	Events  []CalendarEvent
}

// WeatherData holds weather information.
type WeatherData struct {
	Description  string
	Temperature  float64
	Humidity     int
	PrecipChance int
	WeatherCode  int
	Forecasts    []DayForecast
	Hourly       []HourlyWeather
}

// DayForecast holds a single day's forecast.
type DayForecast struct {
	DayLabel     string
	Description  string
	WeatherCode  int
	TempMax      float64
	TempMin      float64
	PrecipChance int
}

// HourlyWeather holds one hour's weather data.
type HourlyWeather struct {
	Hour          int
	Temperature   float64
	WeatherCode   int
	PrecipProb    int
	WindSpeed     float64 // m/s
	WindDirection int     // degrees (0=N, 90=E, 180=S, 270=W)
}

// TrainInfo holds train delay information.
type TrainInfo struct {
	LineName string
	Status   string
	IsDelay  bool
}

// CalendarEvent holds a single calendar event.
type CalendarEvent struct {
	Summary   string
	StartTime time.Time
	EndTime   time.Time
	IsAllDay  bool
}

var (
	fontRegular *truetype.Font
	fontWeather *truetype.Font
	fontMDI     *truetype.Font
)

func init() {
	var err error
	fontRegular, err = loadFont("NotoSansJP-Regular.ttf")
	if err != nil {
		panic("failed to load NotoSansJP: " + err.Error())
	}
	fontWeather, err = loadFont("weathericons.ttf")
	if err != nil {
		panic("failed to load Weather Icons: " + err.Error())
	}
	fontMDI, err = loadFont("materialdesignicons.ttf")
	if err != nil {
		panic("failed to load Material Design Icons: " + err.Error())
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

// Dashboard generates the dashboard image.
func Dashboard(data DashboardData) (image.Image, error) {
	dc := gg.NewContext(Width, Height)

	// White background
	dc.SetColor(color.White)
	dc.Clear()

	drawHeader(dc, data.Now)
	drawSeparator(dc, float64(weatherY))

	if data.Weather != nil {
		drawWeather(dc, data.Weather)
		drawHourly(dc, data.Weather.Hourly)
	}
	drawSeparator(dc, float64(trainY))

	drawTrainInfo(dc, data.Trains)
	drawSeparator(dc, float64(scheduleY))

	drawSchedule(dc, data.Now, data.Events)

	return toGrayscale(dc.Image()), nil
}

func drawSeparator(dc *gg.Context, y float64) {
	dc.SetRGB(separatorGray, separatorGray, separatorGray)
	dc.SetLineWidth(1)
	dc.DrawLine(float64(marginX), y, float64(Width-marginX), y)
	dc.Stroke()
}

func toGrayscale(src image.Image) *image.Gray {
	bounds := src.Bounds()
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			lum := uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
			// Quantize to 4-bit (16 shades) for e-ink
			lum = (lum / 17) * 17
			gray.SetGray(x, y, color.Gray{Y: lum})
		}
	}
	return gray
}
