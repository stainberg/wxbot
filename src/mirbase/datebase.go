package mirbase

import (
	"github.com/go-redis/redis"
	"fmt"
)

var client *redis.Client

var (
	HISTORY_INFO = "HistoryInfo"
	NEW_INFO = "NewInfo"
	WX_TOKEN = "WechatToken"
	TOKEN = "Token"
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
	client.RPush(HISTORY_INFO, info).Result()
	count, _ := client.RPush(NEW_INFO, info).Result()
	return count
}

func FetchNewInfo() (string, error) {
	info, err := client.RPop(NEW_INFO).Result()
	if err != nil {
		return "", err
	}
	return info, nil
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

func BindTokenToName(token, name string) (bool, string) {
	in, err := client.SIsMember(TOKEN, token).Result()
	if err != nil {
		return false, err.Error()
	}
	if !in {
		return false, "invalid token"
	}
	b, e := client.HSet(WX_TOKEN, token, name).Result()
	if e != nil {
		return false, err.Error()
	}
	if !b {
		return false, "bind fail"
	}
	return true, "ok"
}

func FindNameByToken(token string) string {
	in, err := client.SIsMember(TOKEN, token).Result()
	if err != nil {
		return err.Error()
	}
	if !in {
		return "invalid token"
	}
	b, e := client.HExists(WX_TOKEN, token).Result()
	if e != nil {
		return err.Error()
	}
	if !b {
		return "did not bind token to any name"
	}
	name, e := client.HGet(WX_TOKEN, token).Result()
	if e != nil {
		return err.Error()
	}
	return name
}
