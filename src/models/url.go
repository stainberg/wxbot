package models

import (
	"github.com/stainberg/koalart"
	"net/http"
	"strings"
	"mirbase"
)

type UrlController struct {
	koala.Controller
}

func (k *UrlController) URLMapping() {
	k.Mapping(koala.GET, k.Get)
}

func (c *UrlController) Get() {
	s := strings.Split(c.Ctx.Request.URL.Path, "/")
	id := s[len(s) - 1]
	url := mirbase.GetLink(id)
	if len(url) > 0 {
		http.Redirect(c.Ctx.Writer, c.Ctx.Request, url, http.StatusMovedPermanently)
	} else {
		c.Ctx.Writer.WriteHeader(http.StatusNotFound)
	}

}