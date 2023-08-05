package auth

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Logout godoc
// @Summary      User logout
// @Description  logout user and delete token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  payload.HTTPSuccess
// @Failure      401,500  {object}  payload.HTTPError
// @Router       /auth/logout [post]
func (h *handler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token")

		if err != nil {
			h.payload.Unauthorized(c)
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		if err := h.cache.Del(ctx, token).Err(); err != nil {
			h.payload.InternalServerError(c, err)
			return
		}

		h.payload.WriteJSON(c, 200, "ok")
	}
}
