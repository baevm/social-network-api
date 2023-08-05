package models

import (
	"social-network-api/pkg/password"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id         int64             `json:"id"`
	Email      string            `json:"email"`
	Password   password.Password `json:"-"`
	Firstname  string            `json:"first_name"`
	Lastname   string            `json:"last_name"`
	Username   string            `json:"username"`
	Avatar     string            `json:"avatar"`
	Activated  bool              `json:"activated"`
	Birthdate  pgtype.Date       `json:"birthdate" swaggertype:"string" format:"date" example:"2006-01-02"`
	Created_at time.Time         `json:"created_at"`
}
