package posts

import (
	"context"
	"errors"
	"social-network-api/internal/db/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RemoveComment godoc
// @Summary      Remove comment
// @Description  Remove comment from post
// @Tags         posts
// @Produce      json
// @Param 		 Cookie header string true "auth_token" default(Bearer token)
// @Param        id path string true "post id"
// @Success      201  {object} payload.HTTPSuccess
// @Failure      401,500,404  {object}  payload.HTTPError
// @Router       /posts/{id}/comment/{comment_id} [delete]
func (h *handler) RemoveComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64(gin.AuthUserKey)

		strPostId := c.Param("id")
		postId, err := strconv.ParseInt(strPostId, 10, 64)

		if err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		strCommentId := c.Param("comment_id")
		commentId, err := strconv.ParseInt(strCommentId, 10, 64)

		if err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		err = h.postsService.RemoveComment(ctx, postId, commentId, userId)

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
