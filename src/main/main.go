package main

import (
	"github.com/stainberg/koalart"
	_ "routers"
	"mirbase"
	"wx"
)

func main() {
	mirbase.InitClient()
	go func() {
		wx.WxClient = new(wx.WxWeb)
		wx.WxClient.Start()
	}()
	koala.Run("8888")
}
