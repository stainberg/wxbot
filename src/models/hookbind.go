package models

import (
	"github.com/stainberg/koalart"
	"io"
	"mirbase"
	"net/http"
	"utils"
	"wx"
)

type HookBindController struct {
	koala.Controller
}

func (k *HookBindController) URLMapping() {
	k.Mapping(koala.GET, k.Get)
	k.Mapping(koala.POST, k.Post)
}

func (c *HookBindController) Post() {
	if !utils.CheckToken(c.Ctx.Request.Header.Get("token")) {
		c.Ctx.Writer.WriteHeader(http.StatusForbidden)
		io.WriteString(c.Ctx.Writer, "illegal token")
		return
	}
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	name := c.Ctx.Form.Get("name")
	if name == "" {
		io.WriteString(c.Ctx.Writer, `illegal id or token`)
		return
	}
	if wx.WxClient != nil {
		if wx.WxClient.IsLogin() {
			b, msg := mirbase.BindIdToName(wx.WxClient.Uin, name)
			if !b {
				io.WriteString(c.Ctx.Writer, msg)
				return
			}
			io.WriteString(c.Ctx.Writer, `bind ok`)
			return
		}
	}
	io.WriteString(c.Ctx.Writer, `bind failed`)
}
