package models

import (
	"fmt"
	"github.com/stainberg/koalart"
	"io"
	"mirbase"
	"net/http"
	"utils"
)

type HookBindController struct {
	koala.Controller
}

func (k *HookBindController) URLMapping() {
	k.Mapping(koala.GET, k.Get)
	k.Mapping(koala.POST, k.Post)
}

func (c *HookBindController) Get() {
	if !utils.CheckToken(c.Ctx.Request.Header.Get("token")) {
		c.Ctx.Writer.WriteHeader(http.StatusForbidden)
		io.WriteString(c.Ctx.Writer, "illegal token")
		return
	}
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	token := c.Ctx.Query.Get("id")
	_, name := mirbase.FindNameById(token)
	io.WriteString(c.Ctx.Writer, name)
}

func (c *HookBindController) Post() {
	if !utils.CheckToken(c.Ctx.Request.Header.Get("token")) {
		c.Ctx.Writer.WriteHeader(http.StatusForbidden)
		io.WriteString(c.Ctx.Writer, "illegal token")
		return
	}
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	name := c.Ctx.Form.Get("name")
	id := c.Ctx.Form.Get("id")
	if name == "" || id == "" {
		io.WriteString(c.Ctx.Writer, `illegal id or token`)
		return
	}
	b, msg := mirbase.BindIdToName(id, name)
	if !b {
		io.WriteString(c.Ctx.Writer, msg)
		return
	}
	url := fmt.Sprintf(`%s/v1/wechat/hook/%s/send`, c.Ctx.Request.Host, id)
	io.WriteString(c.Ctx.Writer, url)
}
