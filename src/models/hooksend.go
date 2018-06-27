package models

import (
	"github.com/stainberg/koalart"
	"io"
	"mirbase"
	"net/http"
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
	defer c.Ctx.Request.Body.Close()
	if !utils.CheckToken(c.Ctx.Request.Header.Get("token")) {
		c.Ctx.Writer.WriteHeader(http.StatusForbidden)
		io.WriteString(c.Ctx.Writer, "illegal token")
		return
	}
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	message := c.Ctx.Form.Get("message")
	if message == "" {
		io.WriteString(c.Ctx.Writer, `can't send nil message`)
		return
	}
	members := mirbase.GetAllMembers()
	for _, member := range members {
		status, _ := wx.WxClient.SendMessage(message, member)
		if !status {
			io.WriteString(c.Ctx.Writer, `send failed`)
			return
		}
	}
	io.WriteString(c.Ctx.Writer, `send ok`)
}
