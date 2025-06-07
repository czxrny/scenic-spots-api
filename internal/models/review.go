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

func (r *Review) SetId(id string) {
	r.Id = id
}

type NewReview struct {
	SpotId  string  `json:"spotId" validate:"required"`
	Rating  float32 `json:"rating" validate:"required,gte=0,lte=5"`
	Content string  `json:"content" validate:"max=300"`
}

type ReviewInfo struct {
	Rating  float32 `json:"rating" validate:"required,gte=0,lte=5"`
	Content string  `json:"content" validate:"max=300"`
}

type ReviewQueryParams struct {
	SpotId  string
	Limit   string
	AddedBy string
}
