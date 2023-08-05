package posts

import (
	"context"
	"errors"
	"social-network-api/internal/db/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentRequest struct {
	Body string `json:"body" binding:"required,min=1,max=1000"`
}

// Comment godoc
// @Summary      Create comment
// @Description  Create comment on post
// @Tags         posts
// @Produce      json
// @Param 		 Cookie header string true "auth_token" default(Bearer token)
// @Param        id path string true "post id"
// @Param        body  body  CommentRequest  true  "Body"
// @Success      201  {object} payload.HTTPSuccess
// @Failure      401,500,404  {object}  payload.HTTPError
// @Router       /posts/{id}/comment [post]
func (h *handler) Comment() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64(gin.AuthUserKey)

		strid := c.Param("id")
		postId, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		var req CommentRequest

		if err = c.ShouldBindJSON(&req); err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		err = h.postsService.Comment(ctx, postId, userId, req.Body)

		if err != nil {
			if errors.Is(err, models.ErrRecordNotFound) {
				h.payload.NotFound(c)
				return
			} else {
				h.payload.InternalServerError(c, err)
				return
			}
		}

		h.payload.WriteJSON(c, 201, "ok")
	}
}
