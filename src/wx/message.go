package wx

type Message struct {
	AddMsgCount string `json:"AddMsgCount"`
	AddMsgList  []Msg  `json:"AddMsgList"`
}

type Msg struct {
	AppMsgType           int    `json:"AppMsgType"`
	Content              string `json:"Content"`
	FromUserName         string `json:"FromUserName"`
	ToUserName           string `json:"ToUserName"`
	StatusNotifyUserName string `json:"StatusNotifyUserName"`
	MsgType              string `json:"MsgType"`
	/**
	Text 1
	Image 3
	*/
	MsgId    string `json:"MsgId"`
	NewMsgId string `json:"NewMsgId"`
}

type BaseResponse struct {
	ErrMsg string `json:"ErrMsg"`
	Ret    int    `json:"Ret"`
}
