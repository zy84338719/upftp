package response

import (
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func JSON(c *app.RequestContext, statusCode int, data interface{}) {
	c.Header("Content-Type", "application/json")
	c.SetStatusCode(statusCode)
	json.NewEncoder(c.Response.BodyWriter()).Encode(data)
}

func Success(c *app.RequestContext, data interface{}) {
	JSON(c, consts.StatusOK, data)
}

func Error(c *app.RequestContext, statusCode int, message string) {
	JSON(c, statusCode, map[string]string{"error": message})
}

func BadRequest(c *app.RequestContext, message string) {
	Error(c, consts.StatusBadRequest, message)
}

func NotFound(c *app.RequestContext, message string) {
	Error(c, consts.StatusNotFound, message)
}

func InternalError(c *app.RequestContext, message string) {
	Error(c, consts.StatusInternalServerError, message)
}

func Unauthorized(c *app.RequestContext, message string) {
	Error(c, consts.StatusUnauthorized, message)
}

func Forbidden(c *app.RequestContext, message string) {
	Error(c, consts.StatusForbidden, message)
}
