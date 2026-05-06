package dam

import (
	"testing"
	"time"

	"github.com/kangaechu/m5paper-dashboard/internal/render"
)

const sampleHTML = `<html><body>
<table>
<tbody>
<tr style="height: 27pt;">
  <td>ダム名</td>
  <td>有効容量<br />（万ｍ<sup>3</sup>）</td>
  <td>貯水量<br />（万ｍ<sup>3</sup>）</td>
  <td>貯水率<br />（％）</td>
  <td>前日補給量<br />（万ｍ<sup>3</sup>/日）</td>
  <td>平均値に<br />対する割合<br />（％）</td>
</tr>
<tr style="height: 13.5pt;">
  <td class="aly_tx_center">二瀬ダム</td>
  <td class="xl66 aly_tx_right">2,000</td>
  <td class="xl66 aly_tx_right">596</td>
  <td class="xl67 aly_tx_right">30</td>
  <td class="xl67 aly_tx_right">-1</td>
  <td class="xl67 aly_tx_right">44</td>
</tr>
<tr>
  <td>滝沢ダム</td>
  <td>5,800</td>
  <td>1,982</td>
  <td>34</td>
  <td>-7</td>
  <td>53</td>
</tr>
<tr>
  <td>浦山ダム</td>
  <td>5,600</td>
  <td>2,673</td>
  <td>48</td>
  <td>0</td>
  <td>71</td>
</tr>
<tr>
  <td>荒川貯水池</td>
  <td>1,020</td>
  <td>1,014</td>
  <td>99</td>
  <td>0</td>
  <td>101</td>
</tr>
<tr>
  <td>４ダム合計</td>
  <td>14,420<br />&nbsp;&nbsp;&nbsp;※1</td>
  <td>6,265<br />&nbsp;&nbsp;&nbsp;※2</td>
  <td>43<br />※3</td>
  <td>-8<br />&nbsp;※4</td>
  <td>64<br />&nbsp;※5</td>
</tr>
</tbody>
</table>
<table>
<tbody>
<tr>
  <td><span>&nbsp;令和8年5月1日</span><span>0時現在&nbsp;</span></td>
</tr>
</tbody>
</table>
</body></html>`

func TestParseHTML(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	now := time.Date(2026, 5, 1, 12, 0, 0, 0, loc)

	d, err := parseHTML(sampleHTML, now)
	if err != nil {
		t.Fatalf("parseHTML: %v", err)
	}

	if d.SystemName != "荒川水系" {
		t.Errorf("SystemName = %q, want 荒川水系", d.SystemName)
	}

	wantObserved := time.Date(2026, 5, 1, 0, 0, 0, 0, loc)
	if !d.ObservedAt.Equal(wantObserved) {
		t.Errorf("ObservedAt = %v, want %v", d.ObservedAt, wantObserved)
	}

	if got, want := d.Total.EffectiveCapacity, 14420.0; got != want {
		t.Errorf("Total.EffectiveCapacity = %v, want %v", got, want)
	}
	if got, want := d.Total.Storage, 6265.0; got != want {
		t.Errorf("Total.Storage = %v, want %v", got, want)
	}
	if got, want := d.Total.StorageRate, 43.0; got != want {
		t.Errorf("Total.StorageRate = %v, want %v", got, want)
	}
	if got, want := d.StorageRate, 43.0; got != want {
		t.Errorf("StorageRate shortcut = %v, want %v", got, want)
	}

	if len(d.Reservoirs) != 4 {
		t.Fatalf("Reservoirs length = %d, want 4", len(d.Reservoirs))
	}

	type expect struct {
		name string
		eff  float64
		stor float64
		rate float64
	}
	wantList := []expect{
		{"二瀬ダム", 2000, 596, 30},
		{"滝沢ダム", 5800, 1982, 34},
		{"浦山ダム", 5600, 2673, 48},
		{"荒川貯水池", 1020, 1014, 99},
	}
	for i, w := range wantList {
		got := d.Reservoirs[i]
		if got.Name != w.name {
			t.Errorf("Reservoirs[%d].Name = %q, want %q", i, got.Name, w.name)
		}
		if got.EffectiveCapacity != w.eff {
			t.Errorf("Reservoirs[%d].EffectiveCapacity = %v, want %v", i, got.EffectiveCapacity, w.eff)
		}
		if got.Storage != w.stor {
			t.Errorf("Reservoirs[%d].Storage = %v, want %v", i, got.Storage, w.stor)
		}
		if got.StorageRate != w.rate {
			t.Errorf("Reservoirs[%d].StorageRate = %v, want %v", i, got.StorageRate, w.rate)
		}
	}
}

func TestParseObservedAt(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")

	cases := []struct {
		name string
		html string
		want time.Time
	}{
		{
			name: "spans broken",
			html: `<span>令和8年5月1日</span><span>0時現在</span>`,
			want: time.Date(2026, 5, 1, 0, 0, 0, 0, loc),
		},
		{
			name: "plain",
			html: `令和7年12月31日23時現在`,
			want: time.Date(2025, 12, 31, 23, 0, 0, 0, loc),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseObservedAt(c.html, loc)
			if err != nil {
				t.Fatalf("parseObservedAt: %v", err)
			}
			if !got.Equal(c.want) {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestAverageHistory(t *testing.T) {
	avg := AverageHistory()
	// Expected count: full year minus 5 non-existent days
	// (2/30, 2/31, 4/31, 6/31, 9/31, 11/31) = 366 - 6 = ... wait: 2/30 and 2/31 are both
	// non-existent so two days for Feb. Plus 4/31, 6/31, 9/31, 11/31 = 6 missing.
	// Total = 12*31 - 6 = 366. So out length = 366.
	if got, want := len(avg), 366; got != want {
		t.Errorf("AverageHistory length = %d, want %d", got, want)
	}

	// Spot checks against the source PDF table.
	cases := []struct {
		date    string
		storage int
	}{
		{"2024-01-01", 8337},
		{"2024-05-01", 9964},
		{"2024-12-31", 8690},
		{"2024-02-29", 8286}, // leap-day average from H24/H28
	}
	for _, c := range cases {
		var found *render.DailyStorageRate
		for i := range avg {
			if avg[i].Date == c.date {
				found = &avg[i]
				break
			}
		}
		if found == nil {
			t.Errorf("entry for %s not found", c.date)
			continue
		}
		wantRate := float64(c.storage) / EffectiveCapacity * 100
		if found.StorageRate != wantRate {
			t.Errorf("rate for %s = %v, want %v", c.date, found.StorageRate, wantRate)
		}
	}
}

func TestParseFloat(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"2,000", 2000},
		{"14,420 ※1", 14420},
		{"-8 ※4", -8},
		{"43.5", 43.5},
		{"", 0},
		{"abc", 0},
	}
	for _, c := range cases {
		got := parseFloat(c.in)
		if got != c.want {
			t.Errorf("parseFloat(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
