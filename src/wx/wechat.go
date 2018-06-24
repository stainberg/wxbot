package wx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"utils"
)

var WxClient *WxWeb

func init() {
	WxClient = new(WxWeb)
	WxClient.stop = true
	WxClient.stopped = true
	WxClient.fileSerRun = false
	WxClient.login_in = false
}

type WxWeb struct {
	uuid         string
	base_uri     string
	redirect_uri string
	Uin          string
	sid          string
	skey         string
	deviceId     string
	SyncKey      map[string]interface{}
	synckey      string
	User         map[string]interface{}
	BaseRequest  map[string]interface{}
	syncHost     string
	http_client  *http.Client
	contact      *Contact
	stop         bool
	stopped      bool
	fileSerRun   bool
	login_in     bool
}

func (self *WxWeb) _unixStr() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func (self *WxWeb) _run(desc string, f func(...interface{}) bool, args ...interface{}) {
	start := time.Now().UnixNano()
	utils.Log("WxChat _run", desc)
	var result bool
	if len(args) > 1 {
		result = f(args)
	} else if len(args) == 1 {
		result = f(args[0])
	} else {
		result = f()
	}
	useTime := fmt.Sprintf("%.5f", (float64(time.Now().UnixNano()-start) / 1000000000))
	if result {
		utils.Log("WxChat _run", "成功,用时" + useTime + "秒")
	} else {
		utils.Log("WxChat _run", "失败")
		self.Stop()

	}
}

func (self *WxWeb) _init() {
	gCookieJar, _ := cookiejar.New(nil)
	httpclient := http.Client{
		CheckRedirect: nil,
		Jar:           gCookieJar,
		Timeout:       60 * time.Second,
	}
	self.http_client = &httpclient
	rand.Seed(time.Now().Unix())
	str := strconv.Itoa(rand.Int())
	self.deviceId = "e" + str[2:17]
}

func (self *WxWeb) getUuid(args ...interface{}) bool {
	urlstr := "https://login.weixin.qq.com/jslogin"
	urlstr += "?appid=wx782c26e4c19acffb&fun=new&lang=zh_CN&_=" + self._unixStr()
	data, _ := self._get(urlstr, false)
	re := regexp.MustCompile(`"([\S]+)"`)
	find := re.FindStringSubmatch(data)
	if len(find) > 1 {
		self.uuid = find[1]
		return true
	} else {
		return false
	}
}

func (self *WxWeb) genQRcode(args ...interface{}) bool {
	urlstr := "https://login.weixin.qq.com/qrcode/" + self.uuid
	path := "qrcode.jpg"
	out, err := os.Create(path)
	resp, err := self._get(urlstr, false)
	_, err = io.Copy(out, bytes.NewReader([]byte(resp)))
	if err != nil {
		return false
	} else {
		if runtime.GOOS == "darwin" {
			exec.Command("open", path).Run()
		} else {
			if !self.fileSerRun {
				go func() {
					self.fileSerRun = true
					http.HandleFunc("/qr", func(w http.ResponseWriter, req *http.Request) {
						http.ServeFile(w, req, "qrcode.jpg")
						return
					})
					http.ListenAndServe(":" + utils.Conf.HttpConf.FileDownloadPort, nil)
				}()
			}
		}
		return true
	}
}

func (self *WxWeb) _post(urlstr string, params map[string]interface{}, jsonFmt bool) (string, error) {
	var err error
	var resp *http.Response
	if jsonFmt == true {
		jsonPost := utils.JsonEncode(params)
		requestBody := bytes.NewBuffer([]byte(jsonPost))
		request, err := http.NewRequest("POST", urlstr, requestBody)
		if err != nil {
			return "", err
		}
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
		request.Header.Add("Origin", "https://wx.qq.com")
		request.Header.Add("Referer", "https://wx.qq.com/")
		request.Header.Add("User-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
		resp, err = self.http_client.Do(request)
		// resp, err = self.http_client.Post(urlstr, "application/json;charset=utf-8", requestBody)
	} else {
		v := url.Values{}
		for key, value := range params {
			v.Add(key, value.(string))
		}
		resp, err = self.http_client.PostForm(urlstr, v)
	}

	if err != nil {
		utils.Log("WxChat _post", err.Error())
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Log("WxChat _post", err.Error())
		return "", err
	} else {
		defer resp.Body.Close()
	}
	return string(body), nil
}

func (self *WxWeb) _get(urlstr string, jsonFmt bool) (string, error) {
	var err error
	res := ""
	request, _ := http.NewRequest("GET", urlstr, nil)
	request.Header.Add("Referer", "https://wx.qq.com/")
	request.Header.Add("User-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
	resp, err := self.http_client.Do(request)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	return string(body), nil
}

func (self *WxWeb) waitForLogin(tip int) bool {
	_url := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login"
	_url += "?loginicon=true&tip=0&uuid=" + self.uuid + "&_=" + self._unixStr() + "&r=" + strconv.Itoa(int(time.Now().Unix())/1579)
	data, _ := self._get(_url, false)
	re := regexp.MustCompile(`window.code=(\d+);`)
	find := re.FindStringSubmatch(data)
	if len(find) > 1 {
		code := find[1]
		if code == "201" {
			return false
		} else if code == "200" {
			re := regexp.MustCompile(`window.redirect_uri="(\S+?)";`)
			find := re.FindStringSubmatch(data)
			if len(find) > 1 {
				r_uri := find[1] + "&fun=new"
				self.redirect_uri = r_uri
				re = regexp.MustCompile(`/`)
				finded := re.FindAllStringIndex(r_uri, -1)
				self.base_uri = r_uri[:finded[len(finded)-1][0]]
				return true
			}
			return false
		} else if code == "408" {
			utils.Log("WxChat waitForLogin", "[登陆超时]")
		} else {
			utils.Log("WxChat waitForLogin", "[登陆异常]")
		}
	}
	return false
}

func (self *WxWeb) login(args ...interface{}) bool {
	data, _ := self._get(self.redirect_uri, false)
	type Result struct {
		Skey        string `xml:"skey"`
		Wxsid       string `xml:"wxsid"`
		Wxuin       string `xml:"wxuin"`
		Pass_ticket string `xml:"pass_ticket"`
	}
	v := Result{}
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		utils.Log("Wechat", err.Error())
		return false
	}
	self.skey = v.Skey
	self.sid = v.Wxsid
	self.Uin = v.Wxuin
	self.BaseRequest = make(map[string]interface{})
	self.BaseRequest["Uin"], _ = strconv.Atoi(v.Wxuin)
	self.BaseRequest["Sid"] = v.Wxsid
	self.BaseRequest["Skey"] = v.Skey
	self.BaseRequest["DeviceID"] = self.deviceId
	return true
}

func (self *WxWeb) webwxinit(args ...interface{}) bool {
	wxurl := fmt.Sprintf("%s/webwxinit?r=%s", self.base_uri, self._unixStr())
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	res, err := self._post(wxurl, params, true)
	if err != nil {
		return false
	}
	d := utils.JsonDecode(res)
	switch vtype := d.(type) {
	case bool:
		return vtype
	}
	data := d.(map[string]interface{})
	self.User = data["User"].(map[string]interface{})
	self.SyncKey = data["SyncKey"].(map[string]interface{})
	self._setsynckey()

	//interface int和int型不能使用==
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(int)
	self.login_in = true
	return retCode == 0
}

func (self *WxWeb) _setsynckey() {
	keys := []string{}
	for _, keyVal := range self.SyncKey["List"].([]interface{}) {
		key := strconv.Itoa(int(keyVal.(map[string]interface{})["Key"].(int)))
		value := strconv.Itoa(int(keyVal.(map[string]interface{})["Val"].(int)))
		keys = append(keys, key+"_"+value)
	}
	self.synckey = strings.Join(keys, "|")
}

func (self *WxWeb) synccheck() (string, string) {
	urlstr := fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/synccheck", self.syncHost)
	v := url.Values{}
	v.Add("r", self._unixStr())
	v.Add("skey", self.skey)
	v.Add("sid", self.sid)
	v.Add("uin", self.Uin)
	v.Add("deviceid", self.deviceId)
	v.Add("synckey", self.synckey)
	v.Add("_", self._unixStr())
	urlstr = urlstr + "?" + v.Encode()
	data, _ := self._get(urlstr, false)
	paren := `window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`
	re := regexp.MustCompile(paren)
	find := re.FindStringSubmatch(data)
	if len(find) > 2 {
		retcode := find[1]
		selector := find[2]
		return retcode, selector
	} else {
		return "9999", "0"
	}
}

func (self *WxWeb) testsynccheck(args ...interface{}) bool {
	SyncHost := []string{
		"webpush.wx2.qq.com",
		"webpush.wx8.qq.com",
		"webpush.wx.qq.com",
		"webpush.web2.wechat.com",
		"webpush.web.wechat.com",
		"webpush2.wx.qq.com",
		"webpush.wechat.com",
		"webpush1.wechat.com",
		"webpush2.wechat.com",
		"webpush1.wechatapp.com",
		"webpush.wechatapp.com",
	}
	for _, host := range SyncHost {
		self.syncHost = host
		retcode, _ := self.synccheck()
		if retcode == "0" {
			return true
		}
	}
	return false
}

func (self *WxWeb) webwxstatusnotify(args ...interface{}) bool {
	urlstr := fmt.Sprintf("%s/webwxstatusnotify", self.base_uri)
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	params["Code"] = 3
	params["FromUserName"] = self.User["UserName"]
	params["ToUserName"] = self.User["UserName"]
	params["ClientMsgId"] = int(time.Now().Unix())
	res, err := self._post(urlstr, params, true)
	if err != nil {
		return false
	}
	d := utils.JsonDecode(res)
	switch vtype := d.(type) {
	case bool:
		return vtype
	}
	data := d.(map[string]interface{})
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(int)
	return retCode == 0
}

func (self *WxWeb) webwxsync() interface{} {
	urlstr := fmt.Sprintf("%s/webwxsync?sid=%s&skey=%s", self.base_uri, self.sid, self.skey)
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	params["SyncKey"] = self.SyncKey
	params["rr"] = ^int(time.Now().Unix())
	res, err := self._post(urlstr, params, true)
	if err != nil {
		return false
	}
	d := utils.JsonDecode(res)
	switch vtype := d.(type) {
	case bool:
		return vtype
	}
	data := d.(map[string]interface{})
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(int)
	if retCode == 0 {
		self.SyncKey = data["SyncKey"].(map[string]interface{})
		self._setsynckey()
	}
	return data
}

func (self *WxWeb) handleMsg(r interface{}) {
	for _, msg := range r.(map[string]interface{})["AddMsgList"].([]interface{}) {
		msgType := msg.(map[string]interface{})["MsgType"].(int)
		content := msg.(map[string]interface{})["Content"].(string)
		content = strings.Replace(content, "&lt;", "<", -1)
		content = strings.Replace(content, "&gt;", ">", -1)
		content = strings.Replace(content, " ", " ", 1)
		if msgType == 1 {

		} else if msgType == 42 {
			member := new(Member)
			info := msg.(map[string]interface{})["RecommendInfo"].(interface{})
			member.UserName = info.(map[string]interface{})["UserName"].(string)
			member.Alias = info.(map[string]interface{})["Alias"].(string)
			member.NickName = info.(map[string]interface{})["NickName"].(string)
			member.Sex = info.(map[string]interface{})["Sex"].(int)
			member.Signature = info.(map[string]interface{})["Signature"].(string)
			self.contact.MemberList = append(self.contact.MemberList, member)
		}
		ModContactCount := r.(map[string]interface{})["ModContactCount"].(int)
		if ModContactCount > 0 {
			for _, contact := range r.(map[string]interface{})["ModContactList"].([]interface{}) {
				find := false
				for _, member := range self.contact.MemberList {
					if member.UserName == contact.(map[string]interface{})["UserName"].(string) {
						member.NickName = contact.(map[string]interface{})["NickName"].(string)
						println("1 Update Nickname = " + member.NickName + " for user = " + member.UserName)
						find = true
						break
					}
				}
				if find {
					continue
				}
				m := new(Member)
				m.NickName = contact.(map[string]interface{})["NickName"].(string)
				m.UserName = contact.(map[string]interface{})["UserName"].(string)
				self.contact.MemberList = append(self.contact.MemberList, m)
				println("2 Add Nickname = " + m.NickName + " for user = " + m.UserName)
			}
		}
	}
}

func (self *WxWeb) webwxsendmsg(message string, toUseNname string) (bool, string) {
	urlstr := fmt.Sprintf("%s/webwxsendmsg", self.base_uri)
	clientMsgId := strconv.Itoa(int(time.Now().UnixNano()))[:13] + "0" + strconv.Itoa(rand.Int())[3:6]
	params := make(map[string]interface{})
	msg := make(map[string]interface{})
	msg["Type"] = 1
	msg["Content"] = message
	msg["FromUserName"] = self.User["UserName"]
	msg["ToUserName"] = toUseNname
	msg["LocalID"] = clientMsgId
	msg["ClientMsgId"] = clientMsgId
	params["Msg"] = msg
	params["Scene"] = 0
	params["BaseRequest"] = self.BaseRequest
	data, err := self._post(urlstr, params, true)
	if err != nil {
		return false, data
	} else {
		return true, data
	}
}

func (self *WxWeb) webwxgetcontact(args ...interface{}) bool {
	urlstr := fmt.Sprintf("%s/webwxgetcontact?r=%s&seq=%s&skey=%s", self.base_uri, self._unixStr(), 0, self.skey)
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	data, err := self._post(urlstr, params, true)
	if err != nil {
		return false
	} else {
		contact := new(Contact)
		err = json.Unmarshal([]byte(data), &contact)
		self.contact = contact
		return true
	}
}

func (self *WxWeb) getUserIdByNickName(nickname string) *Member {
	if len(self.contact.MemberList) > 0 {
		for _, member := range self.contact.MemberList {
			if member.NickName == nickname || member.RemarkName == nickname {
				return member
			}
		}
	}
	return nil
}

func (self *WxWeb) GetContact() {
	go self.webwxgetcontact()
}

func (self *WxWeb) SendMessage(message string, nickname string) (bool, string) {
	toUseName := self.getUserIdByNickName(nickname)
	if toUseName == nil {
		return self.webwxsendmsg(message, nickname)
	} else {
		return self.webwxsendmsg(message, toUseName.UserName)
	}
}

func (self *WxWeb) Stopped() bool {
	return self.stopped
}

func (self *WxWeb) Stop() {
	self.stop = true
	self.login_in = false
}

func (self *WxWeb) IsLogin() bool {
	return self.login_in
}

func (self *WxWeb) Start() {
	WxClient.stop = true
	for !WxClient.stopped {
		time.Sleep(1 * time.Millisecond)
	}
	WxClient.stop = false
	WxClient.stopped = false
	self._init()
	self._run("[*] 正在获取 uuid ... ", self.getUuid)
	self._run("[*] 正在获取 二维码 ... ", self.genQRcode)
	for !WxClient.stop {
		time.Sleep(3 * time.Second)
		if self.waitForLogin(1) == false {
			continue
		}
		utils.Log("WxChat Start", "[*] 请在手机上点击确认以登录 ... ")
		if self.waitForLogin(0) == false {
			continue
		}
		break
	}
	self._run("[*] 正在登录 ... ", self.login)
	self._run("[*] 微信初始化 ... ", self.webwxinit)
	self._run("[*] 开启状态通知 ... ", self.webwxstatusnotify)
	self._run("[*] 进行同步线路测试 ... ", self.testsynccheck)
	self._run("[*] 获取用户列表 ... ", self.webwxgetcontact)
	for !WxClient.stop {
		retcode, selector := self.synccheck()
		if retcode == "0" {
			if selector == "2" {
				r := self.webwxsync()
				switch r.(type) {
				case bool:
				default:
					self.handleMsg(r)
				}
			} else if selector == "6" || selector == "4" {
				self.webwxsync()
			}
		} else if retcode == "1101" {
			utils.Log("WxChat Start", "[*] 重新登陆 在其他地方登录")
			WxClient.stop = true
			self.login_in = false
			break
		} else if retcode == "1100" {
			utils.Log("WxChat Start", "[*] 重新登陆 登出了微信")
			WxClient.stop = true
			self.login_in = false
			break
		} else if retcode == "1102" {
			utils.Log("WxChat Start", "[*] 重新登陆 Cokice过期")
			WxClient.stop = true
			self.login_in = false
			break
		} else {
			utils.Log("WxChat Start", "[*] retcode = " + retcode)
			self.webwxsync()
		}
	}
	p := map[string]interface{}{}
	p["status"] = "-1"
	self._post(utils.Conf.WechatConf.LogoutCallbackUrl, p, false)
	WxClient.stopped = true
	self.login_in = false
}