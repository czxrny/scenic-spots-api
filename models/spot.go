package models

import "time"

type Spot struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Category    string    `json:"category"`
	Photos      []string  `json:"photos"`
	AddedBy     string    `json:"addedBy"`
	CreatedAt   time.Time `json:"createdAt"`
}

type NewSpot struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Category    string  `json:"category"`
}

type SpotQueryParams struct {
	Name      string
	Latitude  string
	Longitude string
	Radius    string
	Category  string
}
