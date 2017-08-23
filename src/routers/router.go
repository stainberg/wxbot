package routers

import (
	"github.com/stainberg/koalart"
	"models"
)

func init() {
	ns := koala.NewNamespace("v1",
		koala.NSNamespace("wechat",
			koala.NSNamespace("hook",
				koala.NSNamespace("id",
					koala.NSController(new(models.HookKeyController)),
				),
				koala.NSNamespace("bind",
					koala.NSController(new(models.HookBindController)),
				),
				koala.NSNamespace(":id",
					koala.NSNamespace("send",
						koala.NSController(new(models.HookSendController)),
					),
				),
			),
			koala.NSNamespace("login",
				koala.NSController(new(models.WxLoginController)),
			),
		),
		koala.NSNamespace("link",
			koala.NSController(new(models.LinkController)),
			koala.NSNamespace(":id",
				koala.NSController(new(models.UrlController)),
			),
		),
	)
	koala.RegisterNamespace(ns)
}
