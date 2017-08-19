package mirbase

import (
	"github.com/go-redis/redis"
	"fmt"
	"utils"
)

var client *redis.Client

var (
	HISTORY_INFO = "HistoryInfo"
	WX_TOKEN = "WechatToken"
	TOKEN_POLL = "TokenPoll"
)

func InitClient() {
	client = redis.NewClient(&redis.Options {
		Addr:     "localhost:8868",
		Password: "reborn",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func SaveInfo(info string) int64 {
	count, _ := client.RPush(HISTORY_INFO, info).Result()
	return count
}

func FetchHistoryInfo(count int64) ([]string, error) {
	c, err := client.LLen(HISTORY_INFO).Result()
	if err != nil {
		return []string{}, err
	}
	if count > c {
		count = c
	}
	rs, err := client.LRange(HISTORY_INFO, c - count, c).Result()
	if err != nil {
		return []string{}, err
	}
	return rs, nil
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
