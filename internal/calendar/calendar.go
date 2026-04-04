package calendar

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"time"

	"golang.org/x/oauth2/google"
	gcal "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/kangaechu/m5paper-dashboard/internal/render"
)

// Fetch retrieves calendar events from Google Calendar API using a service account.
// credentialsJSON is the base64-encoded service account JSON key.
// calendarIDs is a comma-separated list of calendar IDs.
func Fetch(ctx context.Context, credentialsJSON string, calendarIDs []string, now time.Time) ([]render.CalendarEvent, error) {
	jsonKey, err := base64.StdEncoding.DecodeString(credentialsJSON)
	if err != nil {
		return nil, fmt.Errorf("decode credentials: %w", err)
	}

	config, err := google.JWTConfigFromJSON(jsonKey, gcal.CalendarReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("create JWT config: %w", err)
	}

	srv, err := gcal.NewService(ctx, option.WithHTTPClient(config.Client(ctx)))
	if err != nil {
		return nil, fmt.Errorf("create calendar service: %w", err)
	}

	loc := now.Location()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	tomorrowEnd := todayStart.AddDate(0, 0, 2)

	var events []render.CalendarEvent

	for _, calID := range calendarIDs {
		items, err := fetchCalendarEvents(srv, calID, todayStart, tomorrowEnd)
		if err != nil {
			// Log but continue with other calendars
			continue
		}
		events = append(events, items...)
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].IsAllDay != events[j].IsAllDay {
			return events[i].IsAllDay
		}
		return events[i].StartTime.Before(events[j].StartTime)
	})

	return events, nil
}

func fetchCalendarEvents(srv *gcal.Service, calendarID string, timeMin, timeMax time.Time) ([]render.CalendarEvent, error) {
	call := srv.Events.List(calendarID).
		TimeMin(timeMin.Format(time.RFC3339)).
		TimeMax(timeMax.Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime")

	result, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("list events for %s: %w", calendarID, err)
	}

	var events []render.CalendarEvent
	for _, item := range result.Items {
		ev := render.CalendarEvent{
			Summary: item.Summary,
		}

		if item.Start.Date != "" {
			// All-day event
			ev.IsAllDay = true
			t, err := time.ParseInLocation("2006-01-02", item.Start.Date, timeMin.Location())
			if err != nil {
				continue
			}
			ev.StartTime = t
			if item.End.Date != "" {
				t2, _ := time.ParseInLocation("2006-01-02", item.End.Date, timeMin.Location())
				ev.EndTime = t2
			} else {
				ev.EndTime = ev.StartTime.AddDate(0, 0, 1)
			}
		} else {
			// Timed event
			start, err := time.Parse(time.RFC3339, item.Start.DateTime)
			if err != nil {
				continue
			}
			ev.StartTime = start.In(timeMin.Location())

			if item.End.DateTime != "" {
				end, err := time.Parse(time.RFC3339, item.End.DateTime)
				if err == nil {
					ev.EndTime = end.In(timeMin.Location())
				}
			}
			if ev.EndTime.IsZero() {
				ev.EndTime = ev.StartTime.Add(time.Hour)
			}
		}

		events = append(events, ev)
	}

	return events, nil
}
