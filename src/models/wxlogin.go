package models

import (
	"github.com/stainberg/koalart"
	"io"
	"mirbase"
	"net/http"
	"time"
	"wx"
	"utils"
)

type WxLoginController struct {
	koala.Controller
}

func (k *WxLoginController) URLMapping() {
	k.Mapping(koala.POST, k.Post)
}

func (c *WxLoginController) Post() {
	if !utils.CheckToken(c.Ctx.Request.Header.Get("token")) {
		c.Ctx.Writer.WriteHeader(http.StatusForbidden)
		io.WriteString(c.Ctx.Writer, "illegal token")
		return
	}
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	id := c.Ctx.Form.Get("id")
	find, _ := mirbase.FindNameById(id)
	if find {
		wx.WxClient.Stop()
		for !wx.WxClient.Stopped() {
			time.Sleep(1 * time.Millisecond)
		}
		go wx.WxClient.Start()
		io.WriteString(c.Ctx.Writer, `please open on web broswer http://zuluki.com:8889/qr`)
		return
	} else {
		io.WriteString(c.Ctx.Writer, `illegal id`)
	}
}
