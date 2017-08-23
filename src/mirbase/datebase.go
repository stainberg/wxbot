package mirbase

import (
	"github.com/go-redis/redis"
	"utils"
	"math/rand"
	"time"
	"strconv"
)

var client *redis.Client

var (
	WX_ID = "WechatId"
	ID_POLL = "Ids"
	LINK_URL = "LinkUrl"
	URL_LINK = "UrlLink"
	URL_EXPIRED = "UrlExpired"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890+-=_")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func InitClient() {
	client = redis.NewClient(&redis.Options {
		Addr:     utils.Conf.RedisConf.Host + ":" + utils.Conf.RedisConf.Port,
		Password: utils.Conf.RedisConf.Password,
		DB:       utils.Conf.RedisConf.Db,
	})
	_, err := client.Ping().Result()
	if err != nil {
		utils.Log("Database InitClient", err.Error())
	}
}

func NewId() string {
	id := utils.SecurityMD5(utils.GenerateId())
	client.SAdd(ID_POLL, id).Result()
	return id
}

func GetId() string {
	id, err := client.SPop(ID_POLL).Result()
	if err != nil {
		return NewId()
	}
	return id
}

func BindIdToName(id, name string) (bool, string) {
	b, _ := client.HExists(WX_ID, id).Result()
	if b {
		_, err := client.HSet(WX_ID, id, name).Result()
		if err != nil {
			return false, err.Error()
		}
		return true, "ok"
	}
	in, err := client.SIsMember(ID_POLL, id).Result()
	if err != nil {
		return false, err.Error()
	}
	if !in {
		return false, "invalid id"
	}
	_, err = client.HSet(WX_ID, id, name).Result()
	if err != nil {
		return false, err.Error()
	}
	client.SRem(ID_POLL, id).Result()
	return true, "ok"
}

func FindNameById(id string) (bool, string) {
	b, err := client.HExists(WX_ID, id).Result()
	if err != nil {
		return false, err.Error()
	}
	if !b {
		return false, "did not bind id to any name"
	}
	name, err := client.HGet(WX_ID, id).Result()
	if err != nil {
		return false, err.Error()
	}
	return true, name
}

func SaveShortLink(url string) string {
	has, _ := client.HExists(URL_LINK, url).Result()
	if has {
		k, _ := client.HGet(URL_LINK, url).Result()
		return k
	}
	key := getKey()
	client.HSet(LINK_URL, key, url).Result()
	client.HSet(URL_EXPIRED, key, time.Now().Unix()).Result()
	client.HSet(URL_LINK, url, key).Result()
	return key
}

func getKey() string {
	key := randSeq(4)
	result, err := client.HExists(URL_EXPIRED, key).Result()
	if err != nil {
		return getKey()
	}
	if result {
		t, err := client.HGet(URL_EXPIRED, key).Result()
		if err != nil {
			return getKey()
		}
		tm, err := strconv.ParseInt(t, 10, len(t))
		if err != nil {
			return getKey()
		}
		if time.Now().Unix() - tm > 3600 * 24 * 30 {
			return key
		} else {
			return getKey()
		}
	}
	return key
}

func GetLink(key string) string {
	url, _ := client.HGet(LINK_URL, key).Result()
	return url
}


