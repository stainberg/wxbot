package models

import (
	"github.com/stainberg/koalart"
	"io"
	"mirbase"
	"net/http"
	"strings"
	"wx"
	"utils"
)

type HookSendController struct {
	koala.Controller
}

func (k *HookSendController) URLMapping() {
	k.Mapping(koala.POST, k.Post)
}

func (c *HookSendController) Post() {
	if !utils.CheckToken(c.Ctx.Request.Header.Get("token")) {
		c.Ctx.Writer.WriteHeader(http.StatusForbidden)
		io.WriteString(c.Ctx.Writer, "illegal token")
		return
	}
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	message := c.Ctx.Form.Get("message")
	s := strings.Split(c.Ctx.Request.URL.Path, "/")
	id := s[len(s)-2]
	b, name := mirbase.FindNameById(id)
	if !b {
		io.WriteString(c.Ctx.Writer, `name or id not bind`)
		return
	}
	if message == "" {
		io.WriteString(c.Ctx.Writer, `can't send nil message`)
		return
	}
	status, resp := wx.WxClient.SendMessage(message, name)
	if status {
		io.WriteString(c.Ctx.Writer, `send ok msg = ` + resp)
	} else {
		io.WriteString(c.Ctx.Writer, `send fail msg = ` + resp)
	}
}
