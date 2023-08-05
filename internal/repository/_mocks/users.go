package mocks

import (
	"context"
	"social-network-api/internal/db/models"
	"social-network-api/internal/repository/users"
	"time"

	"github.com/stretchr/testify/mock"
)

type UserModel struct {
	mock.Mock
}

func (u *UserModel) Create(ctx context.Context, user *models.User) error {
	ret := u.Called(ctx, user)

	switch {
	case user.Email == "duplicate@email.com":
		return users.ErrDuplicateEmail
	case user.Username == "duplicate_username":
		return users.ErrDuplicateUsername
	default:
		return nil
	}
}

func (u *UserModel) GetById(id int64) (*models.User, error) {
	if id == 1 {
		u := &models.User{
			Id:         1,
			Username:   "User",
			Email:      "user@test.com",
			Created_at: time.Now(),
		}

		return u, nil
	}

	return nil, users.ErrRecordNotFound
}

func (u *UserModel) Exists(id int64) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
