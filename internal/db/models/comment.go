package models

import "time"

type Comment struct {
	PostId     int64     `json:"post_id"`
	UserId     int64     `json:"user_id"`
	Body       string    `json:"body"`
	Images     []Media   `json:"images"`
	Created_at time.Time `json:"created_at"`
}
