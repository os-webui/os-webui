package v1

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/os-webui/os-webui/internal/plugins"
	"github.com/os-webui/os-webui/web"
)

type Plugins struct {
	web.Web
}

func (w *Plugins) bindID(c *gin.Context) (string, bool) {
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
func (w *Plugins) bindPlugin(c *gin.Context) (*plugins.Plugin, bool) {
	id, ok := w.bindID(c)
	if !ok {
		return nil, false
	}
	plugin := plugins.DefaultPluginsManager.Plugin(id)
	if plugin == nil {
		c.String(http.StatusNotFound, `plugin not found`)
		return nil, false
	}
	return plugin, true
}
func (w *Plugins) List(c *gin.Context) {
	items := plugins.DefaultPluginsManager.List(c.Request.Header.Get(`Accept-Language`))
	w.NegotiateData(c, http.StatusOK, items)
}
func (w *Plugins) Get(c *gin.Context) {
	id, ok := w.bindID(c)
	if !ok {
		return
	}
	plugin, info, ok := plugins.DefaultPluginsManager.Get(id, c.Request.Header.Get(`Accept-Language`))
	if !ok {
		c.String(http.StatusNotFound, `plugin not found`)
		return
	}
	items, err := plugin.Features(c.Request.Context(), c.Request.Header.Get(`Accept-Language`))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	features := make([]map[string]any, len(items))
	set := make(map[string]bool)
	for i, v := range items {
		if v.ID == `` || set[v.ID] {
			continue
		}
		if v.Metadata == nil {
			continue
		}
		found, ok := v.Metadata[`name`]
		if !ok {
			continue
		}
		name, ok := found.(string)
		if !ok {
			continue
		}
		set[v.ID] = true

		features[i] = map[string]any{
			`id`:          v.ID,
			`name`:        name,
			`description`: v.Metadata[`description`],
		}
	}

	w.NegotiateData(c, http.StatusOK, map[string]any{
		`info`:     info,
		`features`: features,
	})
}
func (w *Plugins) Features(c *gin.Context) {
	plugin, ok := w.bindPlugin(c)
	if !ok {
		return
	}
	items, err := plugin.Features(c.Request.Context(), c.Request.Header.Get(`Accept-Language`))
	if err == nil {
		w.NegotiateData(c, http.StatusOK, items)
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
func (w *Plugins) LoadConf(c *gin.Context) {
	plugin, ok := w.bindPlugin(c)
	if !ok {
		return
	}
	s, err := plugin.LoadConf(c.Request.Context())
	if err != nil {
		if !os.IsNotExist(err) {
			c.String(http.StatusNotFound, err.Error())
			return
		}
	}
	w.NegotiateData(c, http.StatusOK, s)
}
func (w *Plugins) SaveConf(c *gin.Context) {
	plugin, ok := w.bindPlugin(c)
	if !ok {
		return
	}
	var req struct {
		Data string
	}
	err := w.ShouldBind(c, &req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = plugin.SaveConf(c.Request.Context(), req.Data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}
