package posts

import (
	"social-network-api/internal/redis"
	"social-network-api/internal/services/posts"
	"social-network-api/pkg/payload"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Handler interface {
	CreatePost() gin.HandlerFunc
	DeletePost() gin.HandlerFunc

	Like() gin.HandlerFunc
	RemoveLike() gin.HandlerFunc

	Comment() gin.HandlerFunc
	RemoveComment() gin.HandlerFunc
}

type handler struct {
	logger       *zap.SugaredLogger
	cache        *redis.Client
	payload      *payload.Payload
	postsService posts.Service
}

func New(logger *zap.SugaredLogger, db *pgxpool.Pool, cache *redis.Client) Handler {
	return &handler{
		logger:       logger,
		cache:        cache,
		payload:      payload.New(logger),
		postsService: posts.NewService(db),
	}
}
