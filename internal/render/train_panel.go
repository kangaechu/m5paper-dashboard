package render

import (
	"strings"

	"github.com/fogleman/gg"
)

// trainGroup groups lines by prefix for compact display.
type trainGroup struct {
	label       string
	normalLines []string
	delayed     []TrainInfo
}

func groupTrains(trains []TrainInfo) []trainGroup {
	groups := []struct {
		prefix string
		label  string
	}{
		{"山手線", "山手線"},
		{"都営", "都営"},
		{"東京メトロ", "東京メトロ"},
	}

	groupMap := make(map[string]*trainGroup)
	var order []string
	for _, g := range groups {
		groupMap[g.prefix] = &trainGroup{label: g.label}
		order = append(order, g.prefix)
	}

	for _, t := range trains {
		matched := false
		for _, g := range groups {
			if strings.HasPrefix(t.LineName, g.prefix) {
				tg := groupMap[g.prefix]
				shortName := strings.TrimPrefix(t.LineName, g.prefix)
				if shortName == "" {
					shortName = t.LineName
				}
				if t.IsDelay {
					tg.delayed = append(tg.delayed, TrainInfo{
						LineName: shortName,
						Status:   t.Status,
						IsDelay:  true,
					})
				} else {
					tg.normalLines = append(tg.normalLines, shortName)
				}
				matched = true
				break
			}
		}
		if !matched {
			// Ungrouped line
			prefix := t.LineName
			if _, ok := groupMap[prefix]; !ok {
				groupMap[prefix] = &trainGroup{label: prefix}
				order = append(order, prefix)
			}
			tg := groupMap[prefix]
			if t.IsDelay {
				tg.delayed = append(tg.delayed, t)
			} else {
				tg.normalLines = append(tg.normalLines, "")
			}
		}
	}

	var result []trainGroup
	for _, key := range order {
		if tg, ok := groupMap[key]; ok {
			if len(tg.normalLines) > 0 || len(tg.delayed) > 0 {
				result = append(result, *tg)
			}
		}
	}
	return result
}

func drawTrainInfo(dc *gg.Context, trains []TrainInfo) {
	baseY := float64(trainY)

	dc.SetRGB(0, 0, 0)
	titleFace := fontFace(fontRegular, 20)
	dc.SetFontFace(titleFace)
	dc.DrawString("運行情報", float64(marginX), baseY+20)

	if len(trains) == 0 {
		dc.SetRGB(0.4, 0.4, 0.4)
		smallFace := fontFace(fontRegular, 18)
		dc.SetFontFace(smallFace)
		dc.DrawString("情報を取得できません", float64(marginX), baseY+45)
		return
	}

	groups := groupTrains(trains)
	y := baseY + 42
	normalFace := fontFace(fontRegular, 14)
	delayFace := fontFace(fontRegular, 14)

	// Collect delayed and normal group labels
	var delayedItems []string
	var normalLabels []string

	for _, g := range groups {
		for _, d := range g.delayed {
			name := g.label + d.LineName
			delayedItems = append(delayedItems, name+": "+d.Status)
		}
		if len(g.normalLines) > 0 {
			normalLabels = append(normalLabels, g.label)
		}
	}

	// Show delayed lines individually
	for _, item := range delayedItems {
		if y > baseY+float64(trainHeight)-10 {
			break
		}
		dc.SetRGB(0.2, 0.0, 0.0)
		dc.SetFontFace(delayFace)
		dc.DrawString(item, float64(marginX+5), y)
		y += 22
	}

	// Show normal groups as one line
	if len(normalLabels) > 0 && y < baseY+float64(trainHeight)-10 {
		dc.SetRGB(0.3, 0.3, 0.3)
		dc.SetFontFace(normalFace)
		dc.DrawString("平常運転: "+strings.Join(normalLabels, "、"), float64(marginX+5), y)
	}
}
