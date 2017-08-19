package models

import (
	"github.com/stainberg/koalart"
	"net/http"
	"io"
	"mirbase"
	"wx"
	"time"
)

type WxLoginController struct {
	koala.Controller
}

func (k *WxLoginController) URLMapping() {
	k.Mapping(koala.POST, k.Post)
}

func (c *WxLoginController) Post()  {
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	token := c.Ctx.Form.Get("token")
	find, _ := mirbase.FindNameByToken(token)
	if find {
		wx.WxClient.Stop()
		for !wx.WxClient.Stopped() {
			time.Sleep(1 * time.Millisecond)
		}
		go wx.WxClient.Start()
		io.WriteString(c.Ctx.Writer, `please open on web broswer http://zuluki.com:8889/qr`)
		return
	}
	io.WriteString(c.Ctx.Writer, `illegal token`)
}
