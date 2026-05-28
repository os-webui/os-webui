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

func (p *Plugins) bindID(c *gin.Context) (string, bool) {
	var uri struct {
		ID string `uri:"id" binding:"required"`
	}
	err := c.BindUri(&uri)
	if err != nil {
		return ``, false
	}
	if !plugins.MatchID(uri.ID) {
		c.String(http.StatusBadRequest, `plugin id invalid`)
		return ``, false
	}
	return uri.ID, true
}
func (p *Plugins) List(c *gin.Context) {
	items := plugins.DefaultPluginsManager.List(c.Request.Header.Get(`Accept-Language`))
	p.NegotiateData(c, http.StatusOK, items)
}
func (p *Plugins) Get(c *gin.Context) {
	id, ok := p.bindID(c)
	if !ok {
		return
	}
	info, ok := plugins.DefaultPluginsManager.Get(id, c.Request.Header.Get(`Accept-Language`))
	if ok {
		p.NegotiateData(c, http.StatusOK, info)
	} else {
		c.String(http.StatusNotFound, `plugin not found`)
	}
}
func (p *Plugins) Features(c *gin.Context) {
	id, ok := p.bindID(c)
	if !ok {
		return
	}
	plugin := plugins.DefaultPluginsManager.Plugin(id)
	if plugin == nil {
		c.String(http.StatusNotFound, `plugin not found`)
		return
	}
	items, err := plugin.Features(c.Request.Context(), c.Request.Header.Get(`Accept-Language`))
	if err == nil {
		p.NegotiateData(c, http.StatusOK, items)
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
