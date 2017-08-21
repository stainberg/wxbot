package mirbase

import (
	"fmt"
	"github.com/go-redis/redis"
	"utils"
)

var client *redis.Client

var (
	WX_TOKEN   = "WechatToken"
	TOKEN_POLL = "TokenPoll"
)

func InitClient() {
	client = redis.NewClient(&redis.Options {
		Addr:     utils.Conf.RedisConf.Host + ":" + utils.Conf.RedisConf.Port,
		Password: utils.Conf.RedisConf.Password,
		DB:       utils.Conf.RedisConf.Db,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func NewToken() string {
	token := utils.SecurityMD5(utils.GenerateId())
	client.SAdd(TOKEN_POLL, token).Result()
	return token
}

func GetToken() string {
	token, err := client.SPop(TOKEN_POLL).Result()
	if err != nil {
		return NewToken()
	}
	return token
}

func BindTokenToName(token, name string) (bool, string) {
	b, _ := client.HExists(WX_TOKEN, token).Result()
	if b {
		_, err := client.HSet(WX_TOKEN, token, name).Result()
		if err != nil {
			return false, err.Error()
		}
		return true, "ok"
	}
	in, err := client.SIsMember(TOKEN_POLL, token).Result()
	if err != nil {
		return false, err.Error()
	}
	if !in {
		return false, "invalid token"
	}
	_, err = client.HSet(WX_TOKEN, token, name).Result()
	if err != nil {
		return false, err.Error()
	}
	client.SRem(TOKEN_POLL, token).Result()
	return true, "ok"
}

func FindNameByToken(token string) (bool, string) {
	b, err := client.HExists(WX_TOKEN, token).Result()
	if err != nil {
		return false, err.Error()
	}
	if !b {
		return false, "did not bind token to any name"
	}
	name, err := client.HGet(WX_TOKEN, token).Result()
	if err != nil {
		return false, err.Error()
	}
	return true, name
}
