package weight

import (
	"time"

	"github.com/spf13/cobra"
)

// weight should only be concerned with getting the first
func Weight(cmd *cobra.Command, args []string) {
	timeRange, _ := cmd.Flags().GetString("range")
	dateFrom := getDateFrom(timeRange)

	measureGroups := fetchMeasurements(dateFrom)

	chartPrintMeasurements(measureGroups)
}

func getDateFrom(timeRange string) time.Time {
	switch timeRange {
	case "month":
		return time.Now().AddDate(0, -1, 0)
	case "year":
		return time.Now().AddDate(-1, 0, 0)
	case "2year":
		return time.Now().AddDate(-2, 0, 0)
	case "all":
	default:
		return time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	return time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)
}
