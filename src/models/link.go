package models

import (
	"github.com/stainberg/koalart"
	"io"
	"mirbase"
	"net/http"
	"utils"
)

type LinkController struct {
	koala.Controller
}

func (k *LinkController) URLMapping() {
	k.Mapping(koala.GET, k.Get)
}

func (c *LinkController) Get() {
	if !utils.CheckToken(c.Ctx.Request.Header.Get("token")) {
		c.Ctx.Writer.WriteHeader(http.StatusForbidden)
		io.WriteString(c.Ctx.Writer, "illegal token")
		return
	}
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	url := c.Ctx.Query.Get("url")
	link := mirbase.SaveShortLink(url)
	io.WriteString(c.Ctx.Writer, "http://" + c.Ctx.Request.Host + "/" + link)
}