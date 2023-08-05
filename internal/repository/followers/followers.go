package followers

import (
	"context"
	"errors"
	"social-network-api/internal/db/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	DB *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{DB: db}
}

func (r *Repo) Follow(ctx context.Context, follow *models.Follow) error {
	query := `
	INSERT INTO followers (user_id, follower_id)
	VALUES ($1, $2)`

	args := []any{follow.UserId, follow.FollowerId}

	ct, err := r.DB.Exec(ctx, query, args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrRecordNotFound
	}

	if ct.RowsAffected() == 0 {
		return models.ErrRecordNotFound
	}

	return err
}

func (r *Repo) Unfollow(ctx context.Context, follow *models.Follow) error {
	query := `
	DELETE FROM followers
	WHERE user_id = $1 AND follower_id = $2`

	args := []any{follow.UserId, follow.FollowerId}

	ct, err := r.DB.Exec(ctx, query, args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrRecordNotFound
	}

	if ct.RowsAffected() == 0 {
		return models.ErrRecordNotFound
	}

	return err
}
