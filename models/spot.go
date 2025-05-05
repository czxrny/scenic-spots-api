package models

import "time"

type SpotMap struct {
	Spots map[string]Spot `json:"spots"`
}

// Used for returning complete info about a spot in "GET" /spot method.
type Spot struct {
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Category  string    `json:"category"`
	Photos    []string  `json:"photos"`
	AddedBy   string    `json:"addedBy"`
	CreatedAt time.Time `json:"createdAt"`
}

// Used for adding a new spot in "POST" /spot method.
type NewSpot struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Category  string  `json:"category"`
}

type SpotQueryParams struct {
	Name      string
	Latitude  string
	Longitude string
	Radius    string
	Category  string
}
