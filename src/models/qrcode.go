package models

import (
	"github.com/stainberg/koalart"
	"net/http"
	"os"
	"io/ioutil"
)

type QrCodeController struct {
	koala.Controller
}

func (k *QrCodeController) URLMapping() {
	k.Mapping(koala.GET, k.Get)
}

func (c *QrCodeController) Get() {
	defer c.Ctx.Request.Body.Close()
	//if !utils.CheckToken(c.Ctx.Request.Header.Get("token")) {
	//	c.Ctx.Writer.WriteHeader(http.StatusForbidden)
	//	io.WriteString(c.Ctx.Writer, "illegal token")
	//	return
	//}
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	file, err := os.Open("/home/rbo/qrcode.jpg")
	if err != nil {
		c.Ctx.Writer.WriteHeader(http.StatusNotFound)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.Ctx.Writer.WriteHeader(http.StatusNotFound)
	}
	c.Ctx.Writer.Header().Set("Content-Type", "image/jpeg")
	c.Ctx.Writer.Write(data)
	//io.WriteString(c.Ctx.Writer, `please open on web broswer http://zuluki.com:8889/qr`)
	return
}
