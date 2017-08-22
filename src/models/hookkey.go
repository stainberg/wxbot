package models

import (
	"github.com/stainberg/koalart"
	"io"
	"mirbase"
	"net/http"
	"utils"
)

type HookKeyController struct {
	koala.Controller
}

func (k *HookKeyController) URLMapping() {
	k.Mapping(koala.GET, k.Get)
}

func (c *HookKeyController) Get() {
	if !utils.CheckToken(c.Ctx.Request.Header.Get("token")) {
		c.Ctx.Writer.WriteHeader(http.StatusForbidden)
		io.WriteString(c.Ctx.Writer, "illegal token")
		return
	}
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	io.WriteString(c.Ctx.Writer, mirbase.GetId())
}
