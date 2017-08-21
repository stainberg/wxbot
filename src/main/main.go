package main

import (
	"github.com/stainberg/koalart"
	"mirbase"
	_ "routers"
	"wx"
)

func main() {
	mirbase.InitClient()
	go func() {
		wx.WxClient.Start()
	}()
	koala.Run("8888")
}
