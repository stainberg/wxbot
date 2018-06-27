package routers

import (
	"github.com/stainberg/koalart"
	"models"
)

func init() {
	ns := koala.NewNamespace("v1",
		koala.NSNamespace("wechat",
			koala.NSNamespace("hook",
				koala.NSNamespace("bind",
					koala.NSController(new(models.HookBindController)),
				),
				koala.NSNamespace("send",
					koala.NSController(new(models.HookSendController)),
				),
			),
			koala.NSNamespace("login",
				koala.NSController(new(models.WxLoginController)),
			),
			koala.NSNamespace("qrcode",
				koala.NSController(new(models.QrCodeController)),
			),
		),
	)
	koala.RegisterNamespace(ns)
}
