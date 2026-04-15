package weight

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func saveMeasurements(mResp MeasureResponse) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic("x failed to get home dir")
	}

	var configDirPath = filepath.Join(home, ".config")
	var withingsStorePath = filepath.Join(configDirPath, "withings/measurements.json")

	file, err := os.ReadFile(withingsStorePath)

	var measurements []Measure

	if err := json.Unmarshal(file, &measurements); err != nil {
		log.Panic()
	}

	for _, measure := range measurements {
		measurements = append(measurements, measure)
	}

	// os.WriteFile(withingsStorePath, )
}

// in the future when there is a local cache it should use that first
func fetchMeasurements(from int64, accessToken string, offset int) MeasureResponse {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://wbsapi.withings.net/measure?action=getmeas&meastypes=1,6&category=1&lastupdate=%d&offset=%d", from, offset), nil)
	if err != nil {
		log.Panic("✗ Failed to create request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Panic("✗ Failed to fetch measurements")
	}
	defer resp.Body.Close()

	var result MeasureResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	if result.Status != 0 {
		log.Fatalf("API error, status: %d", result.Status)
	}

	return result
}
