package models

import (
	"github.com/stainberg/koalart"
	"io"
	"mirbase"
	"net/http"
	"strings"
	"wx"
)

type HookSendController struct {
	koala.Controller
}

func (k *HookSendController) URLMapping() {
	k.Mapping(koala.POST, k.Post)
}

func (c *HookSendController) Post() {
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	message := c.Ctx.Form.Get("message")
	s := strings.Split(c.Ctx.Request.URL.Path, "/")
	token := s[len(s)-2]
	b, name := mirbase.FindNameByToken(token)
	if !b {
		io.WriteString(c.Ctx.Writer, name)
		return
	}
	if message == "" {
		io.WriteString(c.Ctx.Writer, `illegal name or token`)
		return
	}
	status := wx.WxClient.SendMessage(message, name)
	if status {
		io.WriteString(c.Ctx.Writer, `send ok`)
	} else {
		io.WriteString(c.Ctx.Writer, `send fail`)
	}
}
