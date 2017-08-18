package models

import (
	"github.com/stainberg/koalart"
	"net/http"
	"io"
	"mirbase"
)

type HookRegisterController struct {
	koala.Controller
}

func (k *HookRegisterController) URLMapping() {
	k.Mapping(koala.GET, k.Get)
	k.Mapping(koala.POST, k.Post)
}

func (c *HookRegisterController) Get()  {
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	token := c.Ctx.Query.Get("token")
	_, name := mirbase.FindNameByToken(token)
	io.WriteString(c.Ctx.Writer, name)
}

func (c *HookRegisterController) Post()  {
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	name := c.Ctx.Form.Get("name")
	token := c.Ctx.Form.Get("token")
	if name == "" || token == "" {
		io.WriteString(c.Ctx.Writer, `illegal name or token`)
		return
	}
	b, msg := mirbase.BindTokenToName(token, name)
	if !b {
		io.WriteString(c.Ctx.Writer, msg)
		return
	}
	io.WriteString(c.Ctx.Writer, `bind ok`)
}
