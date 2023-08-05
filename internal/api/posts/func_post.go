package posts

import (
	"context"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreatePostRequest struct {
	Body   string                  `form:"body" binding:"required"`
	Images []*multipart.FileHeader `form:"images" binding:"required"`
}

// CreatePost godoc
// @Summary      Create post
// @Description  create post
// @Tags         posts
// @Accept       mpfd
// @Produce      json
// @Param 		 Cookie header string true "auth_token" default(Bearer token)
// @Param        images formData file true "image files"
// @Param        body formData string true "post body text"
// @Success      201  {object} payload.HTTPSuccess
// @Failure      401,500,400  {object}  payload.HTTPError
// @Router       /posts [post]
func (h *handler) CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64(gin.AuthUserKey)

		var req CreatePostRequest

		if err := c.ShouldBind(&req); err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		files := req.Images
		body := req.Body

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		err := h.postsService.CreatePost(ctx, files, body, userId)

		if err != nil {
			h.payload.BadRequest(c, err)
			return
		}

		h.payload.WriteJSON(c, http.StatusCreated, "ok")
	}

}
