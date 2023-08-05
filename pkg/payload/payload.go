package payload

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Payload struct {
	logger *zap.SugaredLogger
}

type ValidationError struct {
	Key     string
	Message string
}

type HTTPError struct {
	Code  int16  `json:"code" example:"400"`
	Error string `json:"error" example:"status bad request"`
}

type HTTPSuccess struct {
	Data interface{} `json:"data"`
}

func New(logger *zap.SugaredLogger) *Payload {
	return &Payload{
		logger: logger,
	}
}

func (p *Payload) ReadJSON(c *gin.Context, payload interface{}) []ValidationError {
	if err := c.ShouldBindJSON(&payload); err != nil {
		var ve validator.ValidationErrors

		if errors.As(err, &ve) {
			out := make([]ValidationError, len(ve))

			for i, fe := range ve {
				out[i] = ValidationError{Key: fe.Field(), Message: msgForTag(fe)}
			}

			return out
		}
	}

	return nil
}

func (p *Payload) WriteJSON(c *gin.Context, status int, payload interface{}) {
	c.Header("Content-Type", "application/json")

	data := HTTPSuccess{
		Data: payload,
	}

	c.JSON(status, data)
}

// Bad Request 400
func (p *Payload) BadRequest(c *gin.Context, err error) {
	httpErr := HTTPError{
		Error: err.Error(),
		Code:  http.StatusBadRequest,
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, httpErr)
}

// Validation Error 422
func (p *Payload) ValidationError(c *gin.Context, errors []ValidationError) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"error": errors,
		"code":  http.StatusUnprocessableEntity,
	})
}

// Internal Server Error 500
func (p *Payload) InternalServerError(c *gin.Context, err error) {
	httpErr := HTTPError{
		Error: "The server encountered a problem and could not process your request",
		Code:  http.StatusInternalServerError,
	}

	p.logger.Errorln(err, map[string]interface{}{
		"req_method": c.Request.Method,
		"req_url":    c.Request.URL,
	})

	c.AbortWithStatusJSON(http.StatusInternalServerError, httpErr)
}

// Not Found 404
func (p *Payload) NotFound(c *gin.Context) {
	err := HTTPError{
		Error: "The requested resource was not found",
		Code:  http.StatusNotFound,
	}

	c.AbortWithStatusJSON(http.StatusNotFound, err)
}

// Unauthorized 401
func (p *Payload) Unauthorized(c *gin.Context) {
	err := HTTPError{
		Error: "You are not authorized to access this resource",
		Code:  http.StatusUnauthorized,
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, err)
}

// Forbidden 403
func (p *Payload) InvalidCredentials(c *gin.Context) {
	err := HTTPError{
		Error: "Invalid credentials",
		Code:  http.StatusForbidden,
	}

	c.AbortWithStatusJSON(http.StatusForbidden, err)
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		return "This field must be at least " + fe.Param() + " characters long"
	case "max":
		return "This field must be at most " + fe.Param() + " characters long"
	}
	return fe.Error() // default error
}
