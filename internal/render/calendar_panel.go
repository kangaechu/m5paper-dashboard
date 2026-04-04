package render

import (
	"fmt"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

func drawSchedule(dc *gg.Context, now time.Time, events []CalendarEvent) {
	baseY := float64(scheduleY)

	dc.SetRGB(0, 0, 0)
	titleFace := fontFace(fontRegular, 20)
	dc.SetFontFace(titleFace)
	dc.DrawString("今日の予定", float64(marginX), baseY+20)

	if len(events) == 0 {
		dc.SetRGB(0.4, 0.4, 0.4)
		smallFace := fontFace(fontRegular, 18)
		dc.SetFontFace(smallFace)
		dc.DrawString("予定はありません", float64(marginX), baseY+45)
		return
	}

	eventFace := fontFace(fontRegular, 17)
	timeFace := fontFace(fontRegular, 15)
	y := baseY + 42

	todayEvents := filterEventsForDay(events, now)
	tomorrowEvents := filterEventsForDay(events, now.AddDate(0, 0, 1))

	maxToday := int((float64(Height) - y - 80) / 25) // 明日の予定分も残す
	if len(tomorrowEvents) == 0 {
		maxToday = int((float64(Height) - y - 20) / 25)
	}
	y = drawEventList(dc, todayEvents, eventFace, timeFace, y, maxToday)

	if len(tomorrowEvents) > 0 && y < float64(Height)-50 {
		y += 10
		dc.SetRGB(0, 0, 0)
		dc.SetFontFace(titleFace)
		dc.DrawString("明日の予定", float64(marginX), y)
		y += 25
		drawEventList(dc, tomorrowEvents, eventFace, timeFace, y, 3)
	}
}

func drawEventList(dc *gg.Context, events []CalendarEvent, eventFace, timeFace font.Face, startY float64, maxEvents int) float64 {
	y := startY
	lineH := 28.0

	for i, ev := range events {
		if i >= maxEvents || y > float64(Height)-20 {
			break
		}

		var timeStr string
		if ev.IsAllDay {
			timeStr = "終日"
		} else {
			timeStr = fmt.Sprintf("%s-%s",
				ev.StartTime.Format("15:04"),
				ev.EndTime.Format("15:04"))
		}

		dc.SetRGB(0.3, 0.3, 0.3)
		dc.SetFontFace(timeFace)
		dc.DrawString(timeStr, float64(marginX+5), y)

		dc.SetRGB(0, 0, 0)
		dc.SetFontFace(eventFace)
		summary := truncateString(ev.Summary, 25)
		dc.DrawString(summary, float64(marginX+110), y)

		y += lineH
	}
	return y
}

func filterEventsForDay(events []CalendarEvent, day time.Time) []CalendarEvent {
	var result []CalendarEvent
	dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	dayEnd := dayStart.AddDate(0, 0, 1)

	for _, ev := range events {
		if ev.IsAllDay {
			evDay := time.Date(ev.StartTime.Year(), ev.StartTime.Month(), ev.StartTime.Day(), 0, 0, 0, 0, ev.StartTime.Location())
			if evDay.Equal(dayStart) {
				result = append(result, ev)
			}
		} else if ev.StartTime.Before(dayEnd) && ev.EndTime.After(dayStart) {
			result = append(result, ev)
		}
	}
	return result
}

func truncateString(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes-1]) + "…"
}
