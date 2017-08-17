package models

import (
	"github.com/stainberg/koalart"
	"net/http"
	"io"
	"wx"
)

type WechatController struct {
	koala.Controller
}

func (k *WechatController) URLMapping() {
	k.Mapping(koala.GET, k.Get)
	k.Mapping(koala.POST, k.Post)
}

func (c *WechatController) Get()  {
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	io.WriteString(c.Ctx.Writer, `Get WechatController`)

}

func (c *WechatController) Post()  {
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	message := c.Ctx.Form.Get("message")
	target := c.Ctx.Form.Get("target")
	status := wx.WxClient.SendMessage(message, target)
	if status {
		io.WriteString(c.Ctx.Writer, `sendmessage success`)
	} else {
		io.WriteString(c.Ctx.Writer, `sendmessage failure`)
	}
}
