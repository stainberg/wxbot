package routers

import (
	"github.com/stainberg/koalart"
	"models"
)

func init() {
	ns := koala.NewNamespace("wechat",
		koala.NSController(new(models.WechatController)))
	koala.RegisterNamespace(ns)
}
