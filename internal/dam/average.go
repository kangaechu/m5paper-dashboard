package dam

import (
	"fmt"
	"time"

	"github.com/kangaechu/m5paper-dashboard/internal/render"
)

// EffectiveCapacity is the non-flood-period utility capacity for the
// 4-dam Arakawa system: 二瀬 + 滝沢 + 浦山 + 荒川第一調節池（彩湖）.
// Source: MLIT Kanto Regional Development Bureau PDF
// /ktr_content/content/000722250.pdf
const EffectiveCapacity = 14420.0 // 万m³

// averageStorageH22H30 is the 9-year (平成22-30, 2010-2018) average daily
// total storage volume of the 4 Arakawa dams in 万m³.
//
// Indexing: averageStorageH22H30[month-1][day-1].
// Zero values are placeholders for non-existent calendar days
// (e.g. 2/30, 4/31). 2/29 is included — it is averaged over the
// 2 leap years (2012, 2016) within the source range.
//
// Transcribed verbatim from the PDF table at /ktr_content/content/000722250.pdf.
var averageStorageH22H30 = [12][31]int{
	// 1月
	{8337, 8346, 8351, 8358, 8364, 8364, 8364, 8359, 8359, 8364, 8360, 8356, 8352, 8350, 8348, 8349, 8347, 8345, 8347, 8342, 8334, 8336, 8349, 8361, 8373, 8378, 8378, 8378, 8375, 8369, 8364},
	// 2月 (29 days max, 30/31 are zero placeholders)
	{8356, 8345, 8334, 8329, 8319, 8310, 8305, 8312, 8310, 8298, 8288, 8279, 8279, 8275, 8291, 8297, 8300, 8297, 8298, 8288, 8281, 8285, 8279, 8280, 8281, 8281, 8279, 8284, 8286, 0, 0},
	// 3月
	{8293, 8319, 8346, 8370, 8385, 8405, 8438, 8464, 8490, 8578, 8616, 8641, 8666, 8683, 8704, 8725, 8740, 8752, 8766, 8795, 8818, 8833, 8857, 8879, 8905, 8929, 8953, 8976, 8997, 9014, 9050},
	// 4月 (30 days max, 31 is zero)
	{9087, 9113, 9144, 9199, 9235, 9264, 9295, 9353, 9391, 9427, 9453, 9492, 9529, 9565, 9606, 9651, 9675, 9698, 9750, 9786, 9810, 9828, 9842, 9855, 9861, 9879, 9895, 9915, 9938, 9953, 0},
	// 5月
	{9964, 9982, 10001, 10076, 10107, 10118, 10128, 10135, 10135, 10131, 10113, 10104, 10088, 10084, 10079, 10057, 10018, 9994, 9969, 9932, 9887, 9855, 9809, 9765, 9727, 9681, 9642, 9597, 9550, 9542, 9526},
	// 6月 (30 days max, 31 is zero)
	{9473, 9424, 9369, 9306, 9248, 9189, 9140, 9114, 9063, 9005, 8939, 8890, 8860, 8823, 8781, 8745, 8719, 8674, 8624, 8589, 8532, 8444, 8358, 8252, 8134, 8026, 7919, 7815, 7710, 7626, 0},
	// 7月
	{7546, 7508, 7458, 7408, 7364, 7347, 7326, 7307, 7266, 7237, 7204, 7168, 7133, 7082, 7043, 7016, 7010, 6984, 6939, 6900, 6863, 6842, 6820, 6795, 6764, 6735, 6729, 6711, 6753, 6733, 6711},
	// 8月
	{6677, 6648, 6642, 6634, 6622, 6596, 6554, 6511, 6496, 6455, 6400, 6354, 6317, 6268, 6222, 6175, 6144, 6131, 6122, 6109, 6103, 6093, 6191, 6183, 6188, 6166, 6155, 6146, 6125, 6109, 6184},
	// 9月 (30 days max, 31 is zero)
	{6135, 6140, 6127, 6138, 6102, 6065, 6040, 6012, 5999, 6072, 6028, 6000, 5974, 5961, 5950, 5943, 6127, 6159, 6184, 6200, 6225, 6303, 6307, 6323, 6325, 6320, 6320, 6333, 6342, 6339, 0},
	// 10月
	{6374, 6558, 6639, 6705, 6794, 6865, 6985, 7043, 7078, 7118, 7146, 7173, 7192, 7209, 7263, 7290, 7411, 7450, 7467, 7491, 7526, 7556, 7651, 7851, 7905, 7950, 8033, 8066, 8087, 8129, 8183},
	// 11月 (30 days max, 31 is zero)
	{8227, 8262, 8290, 8309, 8326, 8338, 8342, 8345, 8342, 8339, 8337, 8341, 8342, 8344, 8344, 8344, 8342, 8338, 8334, 8342, 8344, 8346, 8353, 8358, 8365, 8371, 8383, 8391, 8400, 8409, 0},
	// 12月
	{8419, 8427, 8437, 8464, 8481, 8495, 8503, 8516, 8530, 8541, 8547, 8558, 8569, 8579, 8589, 8597, 8604, 8609, 8617, 8624, 8630, 8638, 8648, 8658, 8663, 8667, 8672, 8678, 8678, 8679, 8690},
}

// AverageHistoryKey is the special map key used to store the H22-H30
// daily average storage rate in the YearlyHistory cache.
const AverageHistoryKey = "average"

// AverageHistory returns the 9-year (H22-H30) daily average storage
// rate (%) of the Arakawa 4-dam system as a sequence ordered by date.
// Days that don't exist (2/30, 4/31, 6/31, 9/31, 11/31) are skipped.
// 2/29 is included.
//
// Date strings use a leap representative year ("2024-MM-DD") so 2/29
// is preserved; the chart renderer plots by day-of-year.
func AverageHistory() []render.DailyStorageRate {
	const repYear = 2024 // leap year so 2/29 is a valid date
	loc := time.UTC
	var out []render.DailyStorageRate

	for m := 1; m <= 12; m++ {
		for d := 1; d <= 31; d++ {
			storage := averageStorageH22H30[m-1][d-1]
			if storage == 0 {
				continue
			}
			// Validate the calendar day exists (skips 2/30, 4/31, etc.)
			t := time.Date(repYear, time.Month(m), d, 0, 0, 0, 0, loc)
			if int(t.Month()) != m || t.Day() != d {
				continue
			}
			out = append(out, render.DailyStorageRate{
				Date:        fmt.Sprintf("%04d-%02d-%02d", repYear, m, d),
				StorageRate: float64(storage) / EffectiveCapacity * 100,
			})
		}
	}
	return out
}
