package utils

import (
	"os"
	"bufio"
	"io"
	"strings"
	"strconv"
)

var Conf *Config

type Config struct {
	HttpConf *Http
	WechatConf *Wechat
	RedisConf *Redis
}

type Http struct {
	RestAPIPort string
	FileDownloadPort string
}

type Wechat struct {
	LogoutCallbackUrl string
}

type Redis struct {
	Host string
	Port string
	Password string
	Db int
}

func LoadConfig(file string) {
	c := make(map[string]string, 0)
	Conf = new(Config)
	Conf.HttpConf = new(Http)
	Conf.WechatConf = new(Wechat)
	Conf.RedisConf = new(Redis)
	f, err := os.Open(file)
	if err == nil && f != nil {
		defer f.Close()
		rd := bufio.NewReader(f)
		for {
			line, err := rd.ReadString('\n')
			if err != nil || io.EOF == err {
				break
			}
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				if line[:1] != "#" {
					s := strings.Split(line, ` `)
					if len(s) == 2 {
						c[s[0]] = s[1]
					}
				}
			}
		}
	}

	Conf.HttpConf.RestAPIPort = getValue("port[http]", c, "8888")
	Conf.HttpConf.FileDownloadPort = getValue("port[file]", c, "8889")
	Conf.WechatConf.LogoutCallbackUrl = getValue("wxlogoutcallback", c, "http://localhost:8888/callback")
	Conf.RedisConf.Host = getValue("host[redis]", c, "127.0.0.1")
	Conf.RedisConf.Port = getValue("port[redis]", c, "8868")
	Conf.RedisConf.Password = getValue("password", c, "reborn")
	index, err := strconv.Atoi(getValue("db", c, "1"))
	if err != nil {
		index = 0
		println(err.Error())
	}
	Conf.RedisConf.Db = index

	println(Conf.HttpConf.RestAPIPort)
	println(Conf.HttpConf.FileDownloadPort)
	println(Conf.WechatConf.LogoutCallbackUrl)
	println(Conf.RedisConf.Host)
	println(Conf.RedisConf.Port)
	println(Conf.RedisConf.Password)
	println(Conf.RedisConf.Db)
}

func getValue(key string, m map[string]string, def string) string {
	str := m[key]
	if str == "" {
		return def
	} else {
		return str
	}

}

func SendLogoutCallback() {

}