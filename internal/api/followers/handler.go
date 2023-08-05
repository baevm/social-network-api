package followers

import (
	"social-network-api/internal/redis"
	"social-network-api/internal/services/followers"
	"social-network-api/pkg/payload"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Handler interface {
	Follow() gin.HandlerFunc
	Unfollow() gin.HandlerFunc
}

type handler struct {
	logger        *zap.SugaredLogger
	cache         *redis.Client
	payload       *payload.Payload
	followService followers.Service
}

func New(logger *zap.SugaredLogger, db *pgxpool.Pool, cache *redis.Client) Handler {
	return &handler{
		logger:        logger,
		cache:         cache,
		payload:       payload.New(logger),
		followService: followers.NewService(db),
	}
}


