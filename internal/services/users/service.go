package users

import (
	"context"
	"social-network-api/internal/db/models"
	"social-network-api/internal/repository/users"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindById(ctx context.Context, id int64) (*models.User, error)
	CheckPasswordHash(user *models.User) (bool, error)
}

type service struct {
	userRepo *users.Repo
}

func NewService(db *pgxpool.Pool) Service {
	return &service{
		userRepo: users.NewRepo(db),
	}
}

func (s *service) Create(ctx context.Context, user *models.User) error {
	// check if user with this email already exists
	existEmail, err := s.userRepo.IsEmailUnique(ctx, user.Email)
	if err != nil {
		return err
	}

	if existEmail {
		return models.ErrDuplicateEmail
	}

	// check if user with this username already exists
	existUsername, err := s.userRepo.IsUsernameUnique(ctx, user.Username)
	if err != nil {
		return err
	}

	if existUsername {
		return models.ErrDuplicateUsername
	}

	// hash password
	err = user.Password.HashPassword(user.Password.PlainTextPass)
	if err != nil {
		return err
	}

	user.Activated = false

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *service) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *service) FindById(ctx context.Context, userId int64) (*models.User, error) {
	return s.userRepo.GetById(ctx, userId)
}

func (s *service) CheckPasswordHash(user *models.User) (bool, error) {
	return user.Password.Matches(user.Password.PlainTextPass)
}
