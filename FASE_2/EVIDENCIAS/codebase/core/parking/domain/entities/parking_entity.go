package entities

type Parking struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Location string `json:"location"`
	Zone     string `json:"zone"`
	IsActive bool   `json:"isActive"`
}

type Parkings []Parking
