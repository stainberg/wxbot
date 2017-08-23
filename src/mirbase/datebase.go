package mirbase

import (
	"github.com/go-redis/redis"
	"utils"
)

var client *redis.Client

var (
	WX_ID = "WechatId"
	ID_POLL = "Ids"
	LINK = "Link"
)

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
	return url
}

func GetLink(url string) string {
return url
}


