package followers

import (
	"context"
	"social-network-api/internal/db/models"
	"social-network-api/internal/repository/followers"
	"social-network-api/internal/repository/users"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	Follow(ctx context.Context, userId int64, followerId int64) error
	Unfollow(ctx context.Context, userId int64, followerId int64) error
}

type service struct {
	followersRepo *followers.Repo
	usersRepo     *users.Repo
}

func NewService(db *pgxpool.Pool) Service {
	return &service{
		followersRepo: followers.NewRepo(db),
		usersRepo:     users.NewRepo(db),
	}
}

func (s *service) Follow(ctx context.Context, userId int64, followerId int64) error {
	if userId == followerId {
		return models.ErrCannotFollowYourself
	}

	// first check if given userId to follow exists
	_, err := s.usersRepo.CheckId(ctx, userId)

	if err != nil {
		return err
	}

	follow := &models.Follow{
		UserId:     userId,
		FollowerId: followerId,
	}

	return s.followersRepo.Follow(ctx, follow)
}

func (s *service) Unfollow(ctx context.Context, userId int64, followerId int64) error {
	if userId == followerId {
		return models.ErrCannotFollowYourself
	}

	// first check if given userId to unfollow exists
	_, err := s.usersRepo.CheckId(ctx, userId)

	if err != nil {
		return err
	}

	follow := &models.Follow{
		UserId:     userId,
		FollowerId: followerId,
	}

	return s.followersRepo.Unfollow(ctx, follow)
}
