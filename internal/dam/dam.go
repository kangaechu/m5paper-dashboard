package dam

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kangaechu/m5paper-dashboard/internal/render"
)

// DefaultURL is the MLIT Kanto Regional Development Bureau page that
// publishes the current storage status of the four Arakawa river system
// dams (Futase, Takizawa, Urayama, Arakawa Reservoir).
const DefaultURL = "https://www.ktr.mlit.go.jp/river/shihon/river_shihon00000113.html"

// damNames are the individual dam labels expected in the table.
var damNames = []string{"二瀬ダム", "滝沢ダム", "浦山ダム", "荒川貯水池"}

var (
	tdRegexp         = regexp.MustCompile(`<td[^>]*>([\s\S]*?)</td>`)
	tagRegexp        = regexp.MustCompile(`<[^>]+>`)
	commentRegexp    = regexp.MustCompile(`<!--[\s\S]*?-->`)
	whitespaceRegexp = regexp.MustCompile(`\s+`)
	numberRegexp     = regexp.MustCompile(`-?[\d,]+(?:\.\d+)?`)
	observedAtRegexp = regexp.MustCompile(`令和\s*(\d+)\s*年\s*(\d+)\s*月\s*(\d+)\s*日\s*(\d+)\s*時現在`)
)

// Fetch retrieves current dam data from the given URL.
func Fetch(url string, now time.Time) (*render.DamData, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch dam data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dam data returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read dam data: %w", err)
	}

	return parseHTML(string(body), now)
}

func parseHTML(html string, now time.Time) (*render.DamData, error) {
	observedAt, err := parseObservedAt(html, now.Location())
	if err != nil {
		return nil, err
	}

	rows := strings.Split(html, "<tr")

	var reservoirs []render.DamReservoir
	var total *render.DamReservoir

	for _, row := range rows {
		cells := tdRegexp.FindAllStringSubmatch(row, -1)
		if len(cells) < 4 {
			continue
		}

		name := cleanCell(cells[0][1])
		if name == "" {
			continue
		}

		isDam := isDamRow(name)
		isTotal := !isDam && isTotalRow(name)
		if !isDam && !isTotal {
			continue
		}

		r := render.DamReservoir{
			Name:              name,
			EffectiveCapacity: parseFloat(cleanCell(cells[1][1])),
			Storage:           parseFloat(cleanCell(cells[2][1])),
			StorageRate:       parseFloat(cleanCell(cells[3][1])),
		}
		if r.EffectiveCapacity == 0 && r.Storage == 0 {
			continue
		}

		if isDam {
			reservoirs = append(reservoirs, r)
		} else {
			cp := r
			cp.Name = "4ダム合計"
			total = &cp
		}
	}

	if total == nil {
		return nil, fmt.Errorf("4ダム合計 row not found")
	}
	if len(reservoirs) == 0 {
		return nil, fmt.Errorf("no individual dam rows found")
	}

	return &render.DamData{
		SystemName:  "荒川水系",
		ObservedAt:  observedAt,
		Total:       *total,
		Reservoirs:  reservoirs,
		StorageRate: total.StorageRate,
	}, nil
}

func isDamRow(name string) bool {
	for _, n := range damNames {
		if strings.Contains(name, n) {
			return true
		}
	}
	return false
}

func isTotalRow(name string) bool {
	return strings.Contains(name, "合計")
}

func parseObservedAt(html string, loc *time.Location) (time.Time, error) {
	stripped := tagRegexp.ReplaceAllString(html, "")
	stripped = strings.ReplaceAll(stripped, "&nbsp;", " ")
	m := observedAtRegexp.FindStringSubmatch(stripped)
	if len(m) != 5 {
		return time.Time{}, fmt.Errorf("observed-at timestamp not found")
	}
	reiwa, _ := strconv.Atoi(m[1])
	month, _ := strconv.Atoi(m[2])
	day, _ := strconv.Atoi(m[3])
	hour, _ := strconv.Atoi(m[4])
	year := 2018 + reiwa
	return time.Date(year, time.Month(month), day, hour, 0, 0, 0, loc), nil
}

func cleanCell(s string) string {
	s = commentRegexp.ReplaceAllString(s, "")
	s = tagRegexp.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = whitespaceRegexp.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func parseFloat(s string) float64 {
	m := numberRegexp.FindString(s)
	if m == "" {
		return 0
	}
	m = strings.ReplaceAll(m, ",", "")
	v, _ := strconv.ParseFloat(m, 64)
	return v
}
