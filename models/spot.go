package models

import "time"

type Spot struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Category  string    `json:"category"`
	Photos    []string  `json:"photos"`
	AddedBy   string    `json:"addedBy"`
	CreatedAt time.Time `json:"createdAt"`
}

type SpotMap struct {
	Spots map[string]Spot `json:"spots"`
}

type SpotQueryParams struct {
	Name      string
	Latitude  string
	Longitude string
	Radius    string
	Category  string
}
