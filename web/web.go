package web

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Web struct{}

func (w Web) ShouldBind(c *gin.Context, obj any) error {
	var b binding.Binding = binding.JSON
	switch c.Request.Method {
	case http.MethodGet, http.MethodHead:
		b = binding.Form
	default:
		switch c.ContentType() {

		case binding.MIMEXML, binding.MIMEXML2:
			b = binding.XML
		case binding.MIMEPROTOBUF:
			b = binding.ProtoBuf
		case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
			b = binding.MsgPack
		case binding.MIMEYAML, binding.MIMEYAML2:
			b = binding.YAML
		case binding.MIMETOML:
			b = binding.TOML
		case binding.MIMEMultipartPOSTForm:
			b = binding.FormMultipart
		case binding.MIMEBSON:
			b = binding.BSON
		case binding.MIMEPOSTForm:
			b = binding.Form
		default: // case binding.MIMEJSON:
			b = binding.JSON
		}
	}
	return c.ShouldBindWith(obj, b)
}
func (w Web) NegotiateData(c *gin.Context, code int, data any) {
	switch c.NegotiateFormat(
		binding.MIMEJSON,
		binding.MIMEXML,
		binding.MIMEYAML, binding.MIMEYAML2,
		binding.MIMETOML,
		binding.MIMEPROTOBUF,
		binding.MIMEBSON,
	) {
	case binding.MIMEXML:
		c.XML(code, data)
	case binding.MIMEYAML, binding.MIMEYAML2:
		c.YAML(code, data)
	case binding.MIMETOML:
		c.TOML(code, data)
	case binding.MIMEPROTOBUF:
		c.ProtoBuf(code, data)
	case binding.MIMEBSON:
		c.BSON(code, data)
	default: // case binding.MIMEJSON:
		beauty := c.GetHeader("beauty")
		switch beauty {
		case "1", "True", "TRUE", "true":
			c.IndentedJSON(code, data)
		case "0", "False", "FALSE", "false":
			c.JSON(code, data)
		default:
			ua := strings.ToLower(c.GetHeader("User-Agent"))
			if strings.HasPrefix(ua, "curl") || strings.HasPrefix(ua, "wget") {
				c.IndentedJSON(code, data)
			} else {
				c.JSON(code, data)
			}
		}
	}
}
