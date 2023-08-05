package auth

import (
	"social-network-api/internal/redis"
	"social-network-api/internal/services/users"
	"social-network-api/pkg/payload"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Handler interface {
	Signup() gin.HandlerFunc
	Login() gin.HandlerFunc
	Logout() gin.HandlerFunc
}

type handler struct {
	logger      *zap.SugaredLogger
	userService users.Service
	cache       *redis.Client
	payload     *payload.Payload
}

func New(logger *zap.SugaredLogger, db *pgxpool.Pool, cache *redis.Client) Handler {
	return &handler{
		logger:      logger,
		cache:       cache,
		userService: users.NewService(db),
		payload:     payload.New(logger),
	}
}
