package posts

import (
	"context"
	"errors"
	"social-network-api/internal/db/models"
	"social-network-api/pkg/dbutil"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	DB *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{DB: db}
}

func (r *Repo) CreatePost(ctx context.Context, post *models.Post) error {
	query := `
	INSERT INTO posts (user_id, body)
	VALUES ($1, $2)
	RETURNING id, created_at
	`

	postArgs := []any{post.User.Id, post.Body}

	// Create post row
	err := r.DB.
		QueryRow(ctx, query, postArgs...).
		Scan(&post.Id, &post.Created_at)

	if err != nil {
		return err
	}

	// If there are images, create image rows
	if len(post.Images) > 0 {
		query = `
		INSERT INTO post_images (post_id, url)
		VALUES %s
		`

		argsPerRow := 2
		imageArgs := make([]interface{}, 0, argsPerRow*len(post.Images))

		for _, image := range post.Images {
			imageArgs = append(imageArgs, post.Id, image.Url)
		}

		batchSQLString := dbutil.GetBulkInsertSQLString(query, argsPerRow, len(post.Images))

		_, err = r.DB.Exec(ctx, batchSQLString, imageArgs...)
	}

	return err
}

func (r *Repo) DeletePost(ctx context.Context, postId int64, userId int64) error {
	query := `
	DELETE FROM posts
	WHERE id = $1 AND user_id = $2
	`

	args := []any{postId, userId}

	ct, err := r.DB.Exec(ctx, query, args...)

	if ct.RowsAffected() != 1 {
		return models.ErrRecordNotFound
	}

	return err
}

func (r *Repo) Like(ctx context.Context, postId int64, userId int64) error {
	query := `
	INSERT INTO post_like (post_id, user_id)
	VALUES ($1, $2)
	`

	args := []any{postId, userId}

	ct, err := r.DB.Exec(ctx, query, args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrRecordNotFound
	}

	if ct.RowsAffected() == 0 {
		return models.ErrAlreadyLiked
	}

	return err
}

func (r *Repo) RemoveLike(ctx context.Context, postId int64, userId int64) error {
	query := `
	DELETE FROM post_like
	WHERE post_id = $1 AND user_id = $2
	`

	args := []any{postId, userId}

	ct, err := r.DB.Exec(ctx, query, args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrRecordNotFound
	}

	if ct.RowsAffected() == 0 {
		return models.ErrNotLiked
	}

	return err
}

func (r *Repo) Comment(ctx context.Context, postId int64, userId int64, body string) error {
	query := `
	INSERT INTO comments (post_id, user_id, body)
	VALUES ($1, $2, $3)
	`

	args := []any{postId, userId, body}

	ct, err := r.DB.Exec(ctx, query, args...)

	if ct.RowsAffected() != 1 {
		return models.ErrRecordNotFound
	}

	return err
}

func (r *Repo) RemoveComment(ctx context.Context, postId int64, commentId int64, userId int64) error {
	query := `
	DELETE FROM comments
	WHERE id = $1 AND post_id = $1 AND user_id = $2
	`

	args := []any{commentId, postId, userId}

	ct, err := r.DB.Exec(ctx, query, args...)

	if ct.RowsAffected() != 1 {
		return models.ErrRecordNotFound
	}

	return err
}

func (r *Repo) GetFeed(ctx context.Context, userId int64, limit int64, offset int64) ([]*models.Post, error) {
	// sql query feed of own posts and people you follow
	query := `
	SELECT p.id, p.user_id, p.body, p.created_at, u.username, coalesce(u.avatar, ''), pi.url
	FROM posts p
	INNER JOIN users u ON u.id = p.user_id
	LEFT JOIN post_images pi ON pi.post_id = p.id
	WHERE p.user_id = $1 OR p.user_id IN (
		SELECT user_id
		FROM followers
		WHERE follower_id = $1
	)
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3
	`

	args := []any{userId, limit, offset}

	rows, err := r.DB.Query(ctx, query, args...)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		}

		return nil, err
	}

	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		var post models.Post

		err := rows.Scan(
			&post.Id,
			&post.User.Id,
			&post.Body,
			&post.Created_at,
			&post.User.Username,
			&post.User.Avatar,
			&post.Images,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return posts, nil
}
