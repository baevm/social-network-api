package users

import (
	"context"
	"errors"
	"social-network-api/internal/db/models"
	"time"

	"github.com/gin-gonic/gin"
)

// GetMe godoc
// @Summary      Get authorized user
// @Description  used to get authorized user
// @Tags         users
// @Produce      json
// @Param 		 Cookie header string true "auth_token" default(Bearer token)
// @Success      201  {object} models.User
// @Failure      401,500,404  {object}  payload.HTTPError
// @Router       /user/me [get]
func (h *handler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64(gin.AuthUserKey)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		user, err := h.userService.FindById(ctx, userId)

		if err != nil {
			if errors.Is(err, models.ErrRecordNotFound) {
				h.payload.NotFound(c)
				return
			}

			h.payload.InternalServerError(c, err)
			return
		}

		h.payload.WriteJSON(c, 200, user)
	}
}
