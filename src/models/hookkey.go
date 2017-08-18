package models

import (
	"github.com/stainberg/koalart"
	"net/http"
	"io"
	"mirbase"
)

type HookKeyController struct {
	koala.Controller
}

func (k *HookKeyController) URLMapping() {
	k.Mapping(koala.GET, k.Get)
	k.Mapping(koala.POST, k.Post)
}

func (c *HookKeyController) Get()  {
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	token := c.Ctx.Request.Header.Get("token")
	if token == "YenStainberg" {
		t := mirbase.GetToken()
		io.WriteString(c.Ctx.Writer, t)
	} else {
		io.WriteString(c.Ctx.Writer, `token invalid`)
	}
}

func (c *HookKeyController) Post()  {
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	token := c.Ctx.Request.Header.Get("token")
	if token == "YenStainberg" {
		t := mirbase.NewToken()
		io.WriteString(c.Ctx.Writer, t)
	} else {
		io.WriteString(c.Ctx.Writer, `token invalid`)
	}
}
