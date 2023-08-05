package auth

import (
	"context"
	"net/http"
	"social-network-api/internal/db/models"
	"time"

	"github.com/gin-gonic/gin"
)

type SignupRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
}

type SignupResponse struct {
	Message string `json:"message"`
	UserId  int64  `json:"user_id"`
}

// Signup godoc
// @Summary      User signup
// @Description  signup user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  SignupRequest  true  "Signup"
// @Success      201  {object} SignupResponse
// @Failure      400,500  {object}  payload.HTTPError
// @Router       /auth/signup [post]
func (h *handler) Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input SignupRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		user := &models.User{
			Email:     input.Email,
			Username:  input.Username,
			Firstname: input.Firstname,
			Lastname:  input.Lastname,
		}
		user.Password.PlainTextPass = input.Password

		if err := h.userService.Create(ctx, user); err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		payload := SignupResponse{
			Message: "User created successfully.",
			UserId:  user.Id,
		}

		h.payload.WriteJSON(c, http.StatusCreated, payload)
	}

}
