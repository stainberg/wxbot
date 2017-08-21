package models

import (
	"fmt"
	"github.com/stainberg/koalart"
	"io"
	"mirbase"
	"net/http"
)

type HookBindController struct {
	koala.Controller
}

func (k *HookBindController) URLMapping() {
	k.Mapping(koala.GET, k.Get)
	k.Mapping(koala.POST, k.Post)
}

func (c *HookBindController) Get() {
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	token := c.Ctx.Query.Get("token")
	_, name := mirbase.FindNameByToken(token)
	io.WriteString(c.Ctx.Writer, name)
}

func (c *HookBindController) Post() {
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
	url := fmt.Sprintf(`%s/v1/wechat/hook/%s/send`, c.Ctx.Request.Host, token)
	io.WriteString(c.Ctx.Writer, url)
}
