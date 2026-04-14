package weight

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/ojparkinson/withings/internal/auth"
	"github.com/spf13/cobra"
)

func Weight(cmd *cobra.Command, args []string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configDirPath := filepath.Join(home, ".config")
	withingsPath := filepath.Join(configDirPath, "withings-cli.toml")

	withingsConfigBytes, _ := os.ReadFile(withingsPath)

	withingsConfig, _ := auth.DecodeConfig(withingsConfigBytes)

	if withingsConfig.ExpiresAt < time.Now().Unix() {
		// Turn into auth fun and use refresh rather than error. Need to check if expired and if refresh exists
		log.Panic("You are not logged in")
		return
	}

	moreMeasurements := true
	offset := 0
	initialStart := time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()

	var measureGroups []MeasureGroup

	for moreMeasurements {
		result := fetchMeasurements(initialStart, withingsConfig.AccessToken, offset)

		moreMeasurements = result.Body.More > 0
		offset = result.Body.Offset

		measureGroups = append(measureGroups, result.Body.MeasureGrps...)
	}

	chartPrintMeasurements(measureGroups)
}

// in the future when there is a local cache it should use that first
func fetchMeasurements(from int64, accessToken string, offset int) *MeasureResponse {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://wbsapi.withings.net/measure?action=getmeas&meastypes=1,6&category=1&lastupdate=%d&offset=%d", from, offset), nil)
	if err != nil {
		log.Panic("✓ Failed to create request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Panic("✓ Failed to fetch measurements")
	}

	var result MeasureResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	if result.Status != 0 {
		log.Fatalf("API error, status: %d", result.Status)
	}

	return &result
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

	tslc := timeserieslinechart.New(150, 15, timeserieslinechart.WithYLabelFormatter(func(i int, v float64) string {
		return fmt.Sprintf("%.1f kg", v)
	}))

	tslc.XLabelFormatter = func(i int, v float64) string {
		t := time.Unix(int64(v), 0)
		return t.Format("01/06")
	}

	for _, grp := range measureGroups {
		for _, m := range grp.Measures {
			if m.Type == 1 {
				tslc.Push(timeserieslinechart.TimePoint{
					Time:  time.Unix(grp.Date, 0),
					Value: m.RealValue(),
				})
			}
		}
	}

	tslc.AutoMinY = true
	tslc.DrawBraille()
	fmt.Println(tslc.View())
}

type MeasureResponse struct {
	Status int         `json:"status"`
	Body   MeasureBody `json:"body"`
}

type MeasureBody struct {
	UpdateTime  int64          `json:"updatetime"`
	Timezone    string         `json:"timezone"`
	MeasureGrps []MeasureGroup `json:"measuregrps"`
	More        int            `json:"more"`
	Offset      int            `json:"offset"`
}

type MeasureGroup struct {
	GrpID    int64     `json:"grpid"`
	Attrib   int       `json:"attrib"`
	Date     int64     `json:"date"`
	Created  int64     `json:"created"`
	Modified int64     `json:"modified"`
	Category int       `json:"category"`
	DeviceID string    `json:"deviceid"`
	Measures []Measure `json:"measures"`
	ModelID  *int      `json:"modelid"`
	Model    *string   `json:"model"`
	Comment  *string   `json:"comment"`
}

type Measure struct {
	Value int `json:"value"`
	Type  int `json:"type"`
	Unit  int `json:"unit"`
	Algo  int `json:"algo"`
	FM    int `json:"fm"`
}

// RealValue converts the Withings value/unit encoding to a float.
// e.g. value=118235, unit=-3 → 118.235
func (m Measure) RealValue() float64 {
	return float64(m.Value) * math.Pow(10, float64(m.Unit))
}
