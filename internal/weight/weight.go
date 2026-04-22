package weight

import (
	"fmt"
	"math"
	"time"

	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/spf13/cobra"
)

// weight should only be concerned with getting the first
func Weight(cmd *cobra.Command, args []string) {
	measureGroups := fetchMeasurements()


	// add a flag for --verbose to just print the verbose data
	chartPrintMeasurements(measureGroups)
}


func chartPrintMeasurements(measureGroups []MeasureGroup) {
	tableData := []timeserieslinechart.TimePoint{}

	minVal, maxVal := math.MaxFloat64, -math.MaxFloat64

	for _, grp := range measureGroups {
		for _, m := range grp.Measures {
			if m.Type == 1 {
				v := m.RealValue()
				minVal = min(minVal, v)
				maxVal = max(maxVal, v)

				tableData = append(tableData, timeserieslinechart.TimePoint{
					Time:  time.Unix(grp.Date, 0),
					Value: m.RealValue(),
				})
			}
		}
	}

	margin := (maxVal - minVal) * 0.05
	minVal -= margin
	maxVal += margin

	tslc := timeserieslinechart.New(50, 15, timeserieslinechart.WithYRange(minVal, maxVal), timeserieslinechart.WithYLabelFormatter(func(i int, v float64) string {
		return fmt.Sprintf("%.1f kg", v)
	}))

	tslc.XLabelFormatter = func(i int, v float64) string {
		t := time.Unix(int64(v), 0)
		return t.Format("01/06")
	}

	for _, point := range tableData {
		tslc.Push(point)
	}

	tslc.Draw()
	fmt.Println(tslc.View())
}
