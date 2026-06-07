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

	// Main section (full-width yearly chart with overlaid storage rate)
	mainY      = headerHeight
	mainHeight = Height - headerHeight

	// Separator
	separatorGray = 0.7
)
