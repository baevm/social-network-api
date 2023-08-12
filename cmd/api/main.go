package main

import (
	"social-network-api/cfg"
	"social-network-api/internal/db"
	"social-network-api/internal/http"
	"social-network-api/internal/mail"
	"social-network-api/internal/rabbitmq"
	"social-network-api/internal/redis"
	"social-network-api/pkg/logger"

	"go.uber.org/zap"
)

// @title           Social network API
// @version         1.0
// @description     Twitter like api made with golang.

// @host      localhost:5000
// @BasePath  /v1
func main() {
	logger := logger.New()
	err := cfg.Load(".")

	if err != nil {
		logger.Fatalf("Error reading config: %s", err)
	}

	db, err := db.New()

	if err != nil {
		logger.Fatalf("Error starting db: %s", err)
	}

	defer db.Close()

	cache := redis.New(cfg.Get().Redis.Host, cfg.Get().Redis.Port, cfg.Get().Redis.Pass)

	defer cache.Close()

	queue, err := rabbitmq.NewProducer(cfg.Get().RabbitMQ.URL)

	if err != nil {
		logger.Fatalf("Error starting queue: %s", err)
	}

	mailer := mail.NewEmailSender(cfg.Get().Email.Name, cfg.Get().Email.Address, cfg.Get().Email.Password)

	go startRabbitConsumer(cfg.Get().RabbitMQ.URL, logger, mailer)

	httpServer := http.New(logger, db, cache, queue)
	httpServer.Run()
}

func startRabbitConsumer(url string, logger *zap.SugaredLogger, mailer mail.EmailSender) {
	c, err := rabbitmq.NewConsumer(url, logger, mailer)

	if err != nil {
		logger.Errorln("Failed to start RabbitMQ consumer: ", err)
	}

	c.Start()
}
