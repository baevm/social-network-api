package models

type Follow struct {
	UserId     int64  `json:"user_id"`
	FollowerId int64  `json:"follower_id"`
	FollowedAt string `json:"followed_at"`
}
