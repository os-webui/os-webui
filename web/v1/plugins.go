package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/os-webui/os-webui/internal/plugins"
	"github.com/os-webui/os-webui/web"
)

type Plugins struct {
	web.Web
}

func (p *Plugins) List(c *gin.Context) {
	items := plugins.DefaultPluginsManager.List(c.Request.Header.Get(`Accept-Language`))
	p.NegotiateData(c, http.StatusOK, items)
}
