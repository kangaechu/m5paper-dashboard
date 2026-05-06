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
	mainHeight = Height - headerHeight
	leftWidth  = 320
	rightX     = leftWidth
	rightWidth = Width - leftWidth

	// Separator
	separatorGray = 0.7
)
