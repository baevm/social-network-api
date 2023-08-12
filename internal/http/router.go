package http

import (
	"social-network-api/internal/api/auth"
	"social-network-api/internal/api/feed"
	"social-network-api/internal/api/followers"
	"social-network-api/internal/api/posts"
	"social-network-api/internal/api/users"

	"github.com/gin-gonic/gin"
)

func (s *Server) setHTTPRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(s.ZapLogger(s.logger))
	router.MaxMultipartMemory = 32 << 20 // 32 mb

	authHandler := auth.New(s.logger, s.db, s.cache, s.queue)
	usersHandler := users.New(s.logger, s.db, s.cache, s.queue)
	postsHandler := posts.New(s.logger, s.db, s.cache)
	followHandler := followers.New(s.logger, s.db, s.cache)
	feedHandler := feed.New(s.logger, s.db, s.cache)

	v1 := router.Group("/v1")
	{
		// AUTH
		auth := v1.Group("/auth")
		auth.POST("/signup", authHandler.Signup())
		auth.POST("/login", authHandler.Login())
		auth.POST("/logout", authHandler.Logout())
	}

	{
		// USERS
		users := v1.Group("/user")
		users.GET("/:username", usersHandler.GetUser())

		users.Use(s.AuthSession())
		users.GET("/me", usersHandler.GetMe())
	}

	{
		// POSTS
		posts := v1.Group("/posts")
		posts.Use(s.AuthSession())
		posts.POST("/", postsHandler.CreatePost())
		posts.DELETE("/:id", postsHandler.DeletePost())
		posts.POST("/:id/like", postsHandler.Like())
		posts.DELETE("/:id/like", postsHandler.RemoveLike())
		posts.POST("/:id/comment", postsHandler.Comment())
		posts.DELETE("/:id/comment/:comment_id", postsHandler.RemoveComment())
	}

	{
		// FOLLOWS
		follows := v1.Group("/follow")
		follows.Use(s.AuthSession())
		follows.POST("/:id", followHandler.Follow())
		follows.DELETE("/:id", followHandler.Unfollow())
	}

	{
		// Feed
		feed := v1.Group("/feed")
		feed.Use(s.AuthSession())
		feed.GET("/", feedHandler.GetFeed())
	}

	return router
}
