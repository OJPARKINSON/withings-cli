package weight

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/ojparkinson/withings/internal/auth"
	"github.com/spf13/cobra"
)

// weight should only be concerned with getting the first 
func Weight(cmd *cobra.Command, args []string) {

	accessToken, err := auth.LoadToken()
	if err != nil {
		log.Fatal("Failed to load Authentication")
	}

	moreMeasurements := true
	offset := 0
	initialStart := time.Date(1999, 01, 01, 01, 01, 01,01, time.UTC).Unix()


	// initialStart := time.Now().Add(time.Duration(-2000) * time.Hour).Unix()

	var measureGroups []MeasureGroup

	for moreMeasurements {
		result := fetchMeasurements(initialStart, accessToken, offset)

		moreMeasurements = result.Body.More > 0
		offset = result.Body.Offset

		measureGroups = append(measureGroups, result.Body.MeasureGrps...)
	}

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

	tslc := timeserieslinechart.New(50, 5, timeserieslinechart.WithYRange(minVal, maxVal), timeserieslinechart.WithYLabelFormatter(func(i int, v float64) string {
		return fmt.Sprintf("%.1f kg", v)
	}))

	tslc.XLabelFormatter = func(i int, v float64) string {
		t := time.Unix(int64(v), 0)
		return t.Format("01/06")
	}

	for _, point := range tableData {
		tslc.Push(point)
	}

	tslc.DrawBraille()
	fmt.Println(tslc.View())
}
