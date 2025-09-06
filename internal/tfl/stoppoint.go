package tfl

import (
	"encoding/json"
	"fmt"
	"iter"
)

type StopPoint struct {
	ID         string `json:"id"`
	Letter     string `json:"stopLetter"`
	CommonName string `json:"commonName"`
	Children   []any  `json:"children"`
}

// GET https://api.tfl.gov.uk/Line/{route}/StopPoints?app_id={app_id}
func (c *Client) StopPoints(route string) (iter.Seq[StopPoint], error) {
	resp, err := c.get(fmt.Sprintf("/Line/%s/StopPoints", route))
	if err != nil {
		return nil, fmt.Errorf("error getting stop points: %w", err)
	}
	defer resp.Body.Close()

	var stopPoints []StopPoint
	err = json.NewDecoder(resp.Body).Decode(&stopPoints)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	stopPointsIter := func(yield func(StopPoint) bool) {
		for _, sp := range stopPoints {
			if len(sp.Children) == 0 {
				if !yield(sp) {
					return
				}
			}
		}
	}

	return stopPointsIter, nil
}
