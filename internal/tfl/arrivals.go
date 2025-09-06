package tfl

import (
	"encoding/json"
	"fmt"
	"iter"
	"time"
)

type Arrival struct {
	Type                string        `json:"$type"`
	ID                  string        `json:"id"`
	OperationType       int           `json:"operationType"`
	VehicleID           string        `json:"vehicleId"`
	NaptanID            string        `json:"naptanId"`
	StationName         string        `json:"stationName"`
	LineID              string        `json:"lineId"`
	LineName            string        `json:"lineName"`
	PlatformName        string        `json:"platformName"`
	Direction           string        `json:"direction"`
	Bearing             string        `json:"bearing"`
	TripID              string        `json:"tripId"`
	BaseVersion         string        `json:"baseVersion"`
	DestinationNaptanID string        `json:"destinationNaptanId"`
	DestinationName     string        `json:"destinationName"`
	Timestamp           time.Time     `json:"timestamp"`
	TimeToStation       int           `json:"timeToStation"`
	CurrentLocation     string        `json:"currentLocation"`
	Towards             string        `json:"towards"`
	ExpectedArrival     time.Time     `json:"expectedArrival"`
	TimeToLive          time.Time     `json:"timeToLive"`
	ModeName            string        `json:"modeName"`
	Timing              ArrivalTiming `json:"timing"`
}

type ArrivalTiming struct {
	Type                      string    `json:"$type"`
	CountdownServerAdjustment string    `json:"countdownServerAdjustment"`
	Source                    time.Time `json:"source"`
	Insert                    time.Time `json:"insert"`
	Read                      time.Time `json:"read"`
	Sent                      time.Time `json:"sent"`
	Received                  time.Time `json:"received"`
}

// GET https://api.tfl.gov.uk/StopPoint/490011187S/Arrivals?app_id={app_id}
func (c *Client) Arrivals(stopPointID string) (iter.Seq[Arrival], error) {
	resp, err := c.get(fmt.Sprintf("/StopPoint/%s/Arrivals", stopPointID))
	if err != nil {
		return nil, fmt.Errorf("error getting arrivals: %w", err)
	}
	defer resp.Body.Close()

	var arrivals []Arrival
	err = json.NewDecoder(resp.Body).Decode(&arrivals)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	arrivalsIter := func(yield func(Arrival) bool) {
		for _, arrival := range arrivals {
			if !yield(arrival) {
				return
			}
		}
	}
	return arrivalsIter, nil
}
