package models

import (
	"github.com/stainberg/koalart"
	"io"
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
	defer c.Ctx.Request.Body.Close()
	if !utils.CheckToken(c.Ctx.Request.Header.Get("token")) {
		c.Ctx.Writer.WriteHeader(http.StatusForbidden)
		io.WriteString(c.Ctx.Writer, "illegal token")
		return
	}
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	wx.WxClient.Stop()
	for !wx.WxClient.Stopped() {
		time.Sleep(1 * time.Millisecond)
	}
	go wx.WxClient.Start()
	io.WriteString(c.Ctx.Writer, `please open https://url.link/v1/wechat/qrcode`)
	return
}
