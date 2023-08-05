package models

import (
	"time"
)

type Post struct {
	Id         int64     `json:"id"`
	User       User      `json:"user"`
	Body       string    `json:"body"`
	Images     []Media   `json:"media,omitempty"`
	Created_at time.Time `json:"created_at"`
}
