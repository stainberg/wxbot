package models

import (
	"github.com/stainberg/koalart"
	"net/http"
)

type HookSendController struct {
	koala.Controller
}

func (k *HookSendController) URLMapping() {
	k.Mapping(koala.POST, k.Post)
}

func (c *HookSendController) Post()  {
	println(c.Ctx.Request.Host)
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	//message := c.Ctx.Form.Get("message")
	//token := c.Ctx.Form.Get("token")
	//if message == "" || token == "" {
	//	io.WriteString(c.Ctx.Writer, `illegal name or token`)
	//	return
	//}
	//name := mirbase.FindNameByToken(token)
	//status := wx.WxClient.SendMessage(message, name)
	//if status {
	//	io.WriteString(c.Ctx.Writer, `send ok`)
	//} else {
	//	io.WriteString(c.Ctx.Writer, `send fail`)
	//}
}
