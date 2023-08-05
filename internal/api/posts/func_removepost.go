package posts

import (
	"context"
	"errors"
	"social-network-api/internal/db/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DeletePost godoc
// @Summary      Delete post
// @Description  delete post
// @Tags         posts
// @Produce      json
// @Param 		 Cookie header string true "auth_token" default(Bearer token)
// @Param        id path string true "post id"
// @Success      201  {object} payload.HTTPSuccess
// @Failure      401,500,404  {object}  payload.HTTPError
// @Router       /posts/{id} [delete]
func (h *handler) DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64(gin.AuthUserKey)

		strid := c.Param("id")
		postId, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			h.payload.InternalServerError(c, err)
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		err = h.postsService.DeletePost(ctx, postId, userId)

		if err != nil {
			if errors.Is(err, models.ErrRecordNotFound) {
				h.payload.NotFound(c)
				return
			}

			h.payload.InternalServerError(c, err)
			return
		}

		h.payload.WriteJSON(c, 200, "ok")
	}
}
