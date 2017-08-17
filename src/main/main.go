package main

import (
	"wx"
	"github.com/stainberg/koalart"
	_ "routers"
)

func main() {
	go func() {
		wx.WxClient = new(wx.WxWeb)
		wx.WxClient.Start()
	}()
	koala.Run("8888")
}
