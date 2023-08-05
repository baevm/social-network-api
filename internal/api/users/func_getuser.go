package users

import (
	"context"
	"errors"
	"social-network-api/internal/db/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserResponse struct {
	Id        int64       `json:"id"`
	Firstname string      `json:"first_name"`
	Lastname  string      `json:"last_name"`
	Username  string      `json:"username"`
	Avatar    string      `json:"avatar"`
	Birthdate pgtype.Date `json:"birthdate" swaggertype:"string" format:"date" example:"2006-01-02"`
}

// GetUser godoc
// @Summary      Get user profile
// @Description  get user profile
// @Tags         users
// @Produce      json
// @Success      201  {object} UserResponse
// @Failure      401,500,404  {object}  payload.HTTPError
// @Router       /user/{username} [get]
func (h *handler) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		if len(username) == 0 {
			h.payload.BadRequest(c, errors.New("username is required"))
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		user, err := h.userService.FindByUsername(ctx, username)

		if err != nil {
			if errors.Is(err, models.ErrRecordNotFound) {
				h.payload.NotFound(c)
				return
			}

			h.payload.InternalServerError(c, err)
			return
		}

		userResponse := UserResponse{
			Id:        user.Id,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Username:  user.Username,
			Avatar:    user.Avatar,
			Birthdate: *user.Birthdate,
		}

		h.payload.WriteJSON(c, 200, userResponse)
	}
}
