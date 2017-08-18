package routers

import (
	"github.com/stainberg/koalart"
	"models"
)

func init() {
	ns := koala.NewNamespace("v1",
		koala.NSNamespace("wechat",
			koala.NSNamespace("hook",
				koala.NSNamespace("register",
					koala.NSController(new(models.HookRegisterController)),
				),
				koala.NSNamespace(":id",
					koala.NSNamespace("send",
						koala.NSController(new(models.HookSendController)),
					),
				),
			),
			koala.NSController(new(models.WechatController)),
		),
		koala.NSNamespace("info",
			koala.NSController(new(models.CoinInfoController)),
		),
	)
	koala.RegisterNamespace(ns)
}
