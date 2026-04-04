package train

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kangaechu/m5paper-dashboard/internal/render"
)

const (
	trainInfoURL = "https://transit.yahoo.co.jp/traininfo/detail/%s/"
)

// LineConfig defines a train line to watch.
type LineConfig struct {
	Name string // display name (e.g. "都営三田線")
	Code string // Yahoo rail code (e.g. "129/0")
}

// mdServiceStatus section: <dl><dt>...<status></dt><dd class="..."><p>detail</p></dd></dl>
var statusRegexp = regexp.MustCompile(`id="mdServiceStatus".*?<dt>(.*?)</dt>`)
var detailRegexp = regexp.MustCompile(`id="mdServiceStatus".*?<dd[^>]*>(.*?)</dd>`)
var tagRegexp = regexp.MustCompile(`<[^>]*>`)

// Fetch retrieves train delay info for the specified lines from Yahoo Transit.
func Fetch(lines []LineConfig) ([]render.TrainInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	var result []render.TrainInfo

	for _, line := range lines {
		info, err := fetchLine(client, line)
		if err != nil {
			result = append(result, render.TrainInfo{
				LineName: line.Name,
				Status:   "取得エラー",
				IsDelay:  false,
			})
			continue
		}
		result = append(result, info)
	}

	return result, nil
}

func fetchLine(client *http.Client, line LineConfig) (render.TrainInfo, error) {
	url := fmt.Sprintf(trainInfoURL, line.Code)
	resp, err := client.Get(url)
	if err != nil {
		return render.TrainInfo{}, fmt.Errorf("fetch %s: %w", line.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return render.TrainInfo{}, fmt.Errorf("%s returned %d", line.Name, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return render.TrainInfo{}, fmt.Errorf("read %s: %w", line.Name, err)
	}

	return parseTrainStatus(line.Name, string(body))
}

func parseTrainStatus(name, html string) (render.TrainInfo, error) {
	info := render.TrainInfo{
		LineName: name,
		Status:   "平常運転",
		IsDelay:  false,
	}

	// Extract status from <dt> inside mdServiceStatus
	matches := statusRegexp.FindStringSubmatch(html)
	if len(matches) > 1 {
		status := stripTags(matches[1])
		status = strings.TrimSpace(status)
		if status != "" {
			info.Status = status
			if status != "平常運転" {
				info.IsDelay = true
			}
		}
	}

	return info, nil
}

func stripTags(s string) string {
	return tagRegexp.ReplaceAllString(s, "")
}
