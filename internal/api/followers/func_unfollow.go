package followers

import (
	"context"
	"errors"
	"social-network-api/internal/db/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Unfollow godoc
// @Summary      Unfollow user
// @Description  Unfollow user with user_id to unfollow provided in path.
// Follower_id equals to user who is logged in.
// User who is requesting is follower. User who is being requested is user to unfollow.
// @Tags         follow
// @Accept       json
// @Produce      json
// @Param        id path string true "id of user to follow"
// @Success      200  {object}  payload.HTTPSuccess
// @Failure      422,403,500  {object}  payload.HTTPError
// @Header       200 {string} auth_token string
// @Router       /follow/{id} [delete]
func (h *handler) Unfollow() gin.HandlerFunc {
	return func(c *gin.Context) {
		followerId := c.GetInt64(gin.AuthUserKey)

		strid := c.Param("id")
		userId, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		err = h.followService.Unfollow(ctx, userId, followerId)

		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				h.payload.NotFound(c)
				return
			case errors.Is(err, models.ErrNotFollowed):
				h.payload.BadRequest(c, err)
				return
			default:
				h.payload.InternalServerError(c, err)
				return
			}
		}

		h.payload.WriteJSON(c, 200, "ok")
	}
}
