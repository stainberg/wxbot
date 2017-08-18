package models

import (
	"github.com/stainberg/koalart"
	"net/http"
	"io"
	"strconv"
	"encoding/json"
	"mirbase"
)

type CoinInfoController struct {
	koala.Controller
}

func (k *CoinInfoController) URLMapping() {
	k.Mapping(koala.POST, k.Post)
	k.Mapping(koala.GET, k.Get)
}

func (c *CoinInfoController) Get()  {
	handler(c)
}

func (c *CoinInfoController) Post()  {
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	message := c.Ctx.Form.Get("payload")
	msg := new(Payload)
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		io.WriteString(c.Ctx.Writer, `illegal params`)
		return
	}
	m, err := json.Marshal(msg)
	if err != nil {
		io.WriteString(c.Ctx.Writer, `illegal params`)
		return
	}
	mirbase.SaveInfo(string(m))
	io.WriteString(c.Ctx.Writer, `save ok message = ` + string(m))
}

func handler(c *CoinInfoController) {
	var message string
	c.Ctx.Writer.WriteHeader(http.StatusOK)
	count := c.Ctx.Query.Get("count")
	if count == "" {
		message = getNewInfo()
	} else {
		cc, e := strconv.Atoi(count)
		if e != nil {
			message = getNewInfo()
		} else {
			message = getHistoryInfo(int64(cc))
		}
	}

	io.WriteString(c.Ctx.Writer, message)
}

func getNewInfo() string {
	info := new(Info)
	info.Info = make([]string, 0)
	m, err := mirbase.FetchNewInfo()
	if err != nil {

	}
	if m != "" {
		info.Info = append(info.Info, m)
	}
	data, err := json.Marshal(info)
	if err != nil {

	}
	message := string(data)
	return message
}

func getHistoryInfo(c int64) string {
	info := new(Info)
	info.Info = make([]string, 0)
	m, err := mirbase.FetchHistoryInfo(c)
	if err != nil {

	}
	info.Info = m
	data, err := json.Marshal(info)
	if err != nil {

	}
	message := string(data)
	return message
}

type Info struct {
	Info []string `json:"info"`
}

type Payload struct {
	Message string `json:"msg"`
	Time int64 `json:"time"`
}