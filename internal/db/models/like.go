package models

import "time"

type Like struct {
	UserId     int64     `json:"user_id"`
	PostId     int64     `json:"post_id"`
	Created_at time.Time `json:"created_at"`
}
