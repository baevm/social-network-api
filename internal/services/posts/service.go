package posts

import (
	"context"
	"mime/multipart"
	"social-network-api/cfg"
	"social-network-api/internal/db/models"
	"social-network-api/internal/repository/media"
	"social-network-api/internal/repository/posts"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	CreatePost(ctx context.Context, files []*multipart.FileHeader, body string, userId int64) error
	DeletePost(ctx context.Context, postId int64, userId int64) error

	Like(ctx context.Context, postId int64, userId int64) error
	RemoveLike(ctx context.Context, postId int64, userId int64) error

	Comment(ctx context.Context, postId int64, userId int64, body string) error
	RemoveComment(ctx context.Context, postId int64, commentId int64, userId int64) error

	GetFeed(ctx context.Context, userId int64, page int64, limit int64) ([]*models.Post, error)
}

type service struct {
	postsRepo *posts.Repo
	mediaRepo *media.Repo
}

func NewService(db *pgxpool.Pool) Service {
	return &service{
		postsRepo: posts.NewRepo(db),
		mediaRepo: media.NewRepo(cfg.Get().Cloud.Name, cfg.Get().Cloud.Key, cfg.Get().Cloud.Secret),
	}
}

func (s *service) CreatePost(ctx context.Context, files []*multipart.FileHeader, body string, userId int64) error {
	// Upload files to cloud
	// And save their urls
	media := make([]models.Media, len(files))
	for i, file := range files {
		file, _ := file.Open()
		defer file.Close()

		res, err := s.mediaRepo.Upload(ctx, file, "posts")
		if err != nil {
			return err
		}

		media[i].Url = res.PublicLink
	}

	// Create post
	post := &models.Post{
		Body:   body,
		Images: media,
		User: &models.User{
			Id: userId,
		},
	}

	err := s.postsRepo.CreatePost(ctx, post)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeletePost(ctx context.Context, postId int64, userId int64) error {
	return s.postsRepo.DeletePost(ctx, postId, userId)
}

func (s *service) Like(ctx context.Context, postId int64, userId int64) error {
	return s.postsRepo.Like(ctx, postId, userId)
}

func (s *service) RemoveLike(ctx context.Context, postId int64, userId int64) error {
	return s.postsRepo.RemoveLike(ctx, postId, userId)
}

func (s *service) Comment(ctx context.Context, postId int64, userId int64, body string) error {
	return s.postsRepo.Comment(ctx, postId, userId, body)
}

func (s *service) RemoveComment(ctx context.Context, postId int64, commentId int64, userId int64) error {
	return s.postsRepo.RemoveComment(ctx, postId, commentId, userId)
}

func (s *service) GetFeed(ctx context.Context, userId int64, page int64, limit int64) ([]*models.Post, error) {
	offset := (page - 1) * limit

	return s.postsRepo.GetFeed(ctx, userId, limit, offset)
}
