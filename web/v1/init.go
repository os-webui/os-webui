package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(api *gin.RouterGroup) {
	v1 := api.Group(`v1`)

	var plugins Plugins
	router := v1.Group(`plugins`)
	router.GET(``, plugins.List)
	router.GET(`:id`, plugins.Get)
	router.GET(`:id/features`, plugins.Features)
	router.GET(`:id/run`, notImplemented)
	router.GET(`:id/attach`, notImplemented)
	router.GET(`:id/history`, notImplemented)

	router = api.Group(`store`)
	router.GET(``, notImplemented)
	router.POST(`:id`, notImplemented)
	router.DELETE(`:id`, notImplemented)
}
func notImplemented(c *gin.Context) {
	c.String(http.StatusNotImplemented, `not implemented
`)
}
