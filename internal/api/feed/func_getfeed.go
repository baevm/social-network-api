package feed

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	DEFAULT_PAGE            = "1"
	DEFAULT_LIMIT           = "20"
	DEFAULT_MIN_LIMIT int64 = 10
	DEFAULT_MAX_LIMIT int64 = 50
)

// GetFeed godoc
// @Summary      Get feed of authorized user
// @Description  get feed of authorized user (own posts and people you follow) with pagination
// @Tags         feed
// @Accept       json
// @Produce      json
// @Param 		 Cookie header string true "auth_token" default(Bearer token)
// @Param        id query string false "page"
// @Param        id query string false "limit"
// @Success      201  {object} []models.Post
// @Failure      401,500,400  {object}  payload.HTTPError
// @Router       /feed [get]
func (h *handler) GetFeed() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64(gin.AuthUserKey)

		pagestr := c.DefaultQuery("page", DEFAULT_PAGE)
		limitstr := c.DefaultQuery("limit", DEFAULT_LIMIT)

		page, err := strconv.ParseInt(pagestr, 10, 64)

		if err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		if page < 0 {
			h.payload.BadRequest(c, errors.New("page must be greater than 0"))
			return
		}

		limit, err := strconv.ParseInt(limitstr, 10, 64)

		if err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		if limit < DEFAULT_MIN_LIMIT || limit > DEFAULT_MAX_LIMIT {
			h.payload.BadRequest(c, errors.New("limit must be between 10 and 50"))
			return
		}

		posts, err := h.postsService.GetFeed(c.Request.Context(), userId, page, limit)

		if err != nil {
			h.payload.InternalServerError(c, err)
			return
		}

		h.payload.WriteJSON(c, 200, posts)
	}
}
