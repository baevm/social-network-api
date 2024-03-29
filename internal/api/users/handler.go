package users

import (
	"social-network-api/internal/rabbitmq"
	"social-network-api/internal/redis"
	"social-network-api/internal/services/users"
	"social-network-api/pkg/payload"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Handler interface {
	GetMe() gin.HandlerFunc
	GetUser() gin.HandlerFunc
}

type handler struct {
	logger      *zap.SugaredLogger
	cache       *redis.Client
	payload     *payload.Payload
	userService users.Service
}

func New(logger *zap.SugaredLogger, db *pgxpool.Pool, cache *redis.Client, queue rabbitmq.QueueProducer) Handler {
	return &handler{
		logger:      logger,
		cache:       cache,
		payload:     payload.New(logger),
		userService: users.NewService(db, queue),
	}
}
