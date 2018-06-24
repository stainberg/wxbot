package mirbase

import (
	"github.com/go-redis/redis"
	"utils"
	"wx"
)

var client *redis.Client

var (
	WX_ID = "WechatId"
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

func BindIdToName(id, name string) (bool, string) {
	_, err := client.SAdd(id, name).Result()
	if err != nil {
		return false, err.Error()
	}
	return true, "ok"
}

func GetAllMembers() []string {
	if wx.WxClient != nil {
		if wx.WxClient.IsLogin() {
			members, err := client.SMembers(wx.WxClient.Uin).Result()
			if err != nil {
				return make([]string, 0)
			}
			return members
		}
	}
	return make([]string, 0)
}