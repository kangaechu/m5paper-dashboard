package render

const (
	Width  = 960
	Height = 540

	// Margins
	marginX      = 20
	contentWidth = Width - marginX*2

	// Header
	headerY      = 0
	headerHeight = 40

	// Main section (left: storage rate, right: yearly chart)
	mainY      = headerHeight
	mainHeight = 280
	leftWidth  = 320 // left panel for storage rate
	rightX     = leftWidth
	rightWidth = Width - leftWidth

	// Hourly delta section
	hourlyDeltaY      = mainY + mainHeight
	hourlyDeltaHeight = 180

	// Footer (stats summary)
	footerY      = hourlyDeltaY + hourlyDeltaHeight
	footerHeight = Height - footerY

	// Separator
	separatorGray = 0.7
)
