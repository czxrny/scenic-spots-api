package models

import "time"

type Review struct {
	Id        string    `json:"id"`
	SpotId    string    `json:"spotId"`
	Rating    float32   `json:"rating"`
	Content   string    `json:"content"`
	AddedBy   string    `json:"addedBy"`
	CreatedAt time.Time `json:"createdAt"`
}

type NewReview struct {
	SpotId  string  `json:"spotId"`
	Rating  float32 `json:"rating"`
	Content string  `json:"content"`
}
