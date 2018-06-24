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
	Token string
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
	Conf.HttpConf.Token = getValue("token", c, "97ae848899631d3c6a656c855090a6e5bca4dc65")
	Conf.WechatConf.LogoutCallbackUrl = getValue("wxlogoutcallback", c, "http://localhost:8888/callback")
	Conf.RedisConf.Host = getValue("host[redis]", c, "localhost")
	Conf.RedisConf.Port = getValue("port[redis]", c, "6379")
	Conf.RedisConf.Password = getValue("password", c, "reborn")
	index, err := strconv.Atoi(getValue("db", c, "0"))
	if err != nil {
		index = 0
		println(err.Error())
	}
	Conf.RedisConf.Db = index

	println("RestAPI Port:" + Conf.HttpConf.RestAPIPort)
	println("FileDownload Port:" + Conf.HttpConf.FileDownloadPort)
	println("Token:" + Conf.HttpConf.Token)
	println("LogoutCallbackUrl:" + Conf.WechatConf.LogoutCallbackUrl)
	println("Redis Host:" + Conf.RedisConf.Host)
	println("Redis Post:" + Conf.RedisConf.Port)
	println("Redis Password:" + Conf.RedisConf.Password)
	println("Redis DB:" + strconv.Itoa(Conf.RedisConf.Db))
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