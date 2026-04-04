package render

const (
	Width  = 540
	Height = 960

	// Section Y positions
	headerY        = 0
	headerHeight   = 60
	weatherY       = headerHeight
	weatherHeight  = 195
	hourlyY        = weatherY + weatherHeight
	hourlyHeight   = 135
	trainY         = hourlyY + hourlyHeight
	trainHeight    = 110
	scheduleY      = trainY + trainHeight
	scheduleHeight = Height - scheduleY // = 460

	// Margins
	marginX      = 20
	contentWidth = Width - marginX*2

	// Separator
	separatorGray = 0.7
)
