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

func Weight(cmd *cobra.Command, args []string) {

	accessToken, err := auth.LoadToken()
	if err != nil {
		log.Fatal("Failed to load Authentication")
	}

	moreMeasurements := true
	offset := 0
	initialStart := time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()

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

func verbosePrintMeasurements(result *MeasureResponse) {
	for _, grp := range result.Body.MeasureGrps {
		t := time.Unix(grp.Date, 0)
		for _, m := range grp.Measures {
			switch m.Type {
			case 1:
				fmt.Printf("%s  Weight: %.1f kg\n", t.Format("2006-01-02"), m.RealValue())
			case 6:
				fmt.Printf("%s  Fat:    %.1f%%\n", t.Format("2006-01-02"), m.RealValue())
			}
		}
	}
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

	tslc := timeserieslinechart.New(150, 15, timeserieslinechart.WithYRange(minVal, maxVal), timeserieslinechart.WithYLabelFormatter(func(i int, v float64) string {
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
