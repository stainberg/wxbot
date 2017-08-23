package models

import (
	"github.com/stainberg/koalart"
	"io"
	"net/http"
	"utils"
	"strings"
)

type UrlController struct {
	koala.Controller
}

func (k *UrlController) URLMapping() {
	k.Mapping(koala.GET, k.Get)
}

func (c *UrlController) Get() {
	if !utils.CheckToken(c.Ctx.Request.Header.Get("token")) {
		c.Ctx.Writer.WriteHeader(http.StatusForbidden)
		io.WriteString(c.Ctx.Writer, "illegal token")
		return
	}
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	s := strings.Split(c.Ctx.Request.URL.Path, "/")
	id := s[len(s) - 1]


	io.WriteString(c.Ctx.Writer, "UrlController id = " + id)
}