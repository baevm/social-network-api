package posts

import (
	"context"
	"errors"
	"social-network-api/internal/db/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RemoveLike godoc
// @Summary      Remove like
// @Description  Remove like from post
// @Tags         posts
// @Produce      json
// @Param 		 Cookie header string true "auth_token" default(Bearer token)
// @Param        id path string true "post id"
// @Success      201  {object} payload.HTTPSuccess
// @Failure      401,500,404  {object}  payload.HTTPError
// @Router       /posts/{id}/like [delete]
func (h *handler) RemoveLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64(gin.AuthUserKey)

		strid := c.Param("id")
		postId, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		err = h.postsService.RemoveLike(ctx, postId, userId)

		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				h.payload.NotFound(c)
				return
			case errors.Is(err, models.ErrNotLiked):
				h.payload.BadRequest(c, err)
				return
			default:
				h.payload.InternalServerError(c, err)
				return
			}
		}

		h.payload.WriteJSON(c, 201, "ok")
	}
}
