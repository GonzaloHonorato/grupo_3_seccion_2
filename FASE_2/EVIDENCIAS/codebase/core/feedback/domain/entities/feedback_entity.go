package entities

import "time"

type Feedback struct {
	ID              int        `json:"id"`
	FromUserID      string     `json:"fromUserId"`
	Comment         string     `json:"comment"`
	ResponseComment string     `json:"responseComment"`
	Rating          int        `json:"rating"` 
	CreatedAt       time.Time  `json:"createdAt"`
	ResponseAt      *time.Time `json:"responseAt,omitempty"`
}

type Feedbacks []Feedback
