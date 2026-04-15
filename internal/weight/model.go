package weight

import "math"

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

func (m Measure) RealValue() float64 {
	return float64(m.Value) * math.Pow(10, float64(m.Unit))
}
