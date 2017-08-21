package wx

type Contact struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	MemberCount  int          `json:"MemberCount"`
	MemberList   []*Member    `json:"MemberList"`
	Seq          int          `json:"Seq"`
}

type Member struct {
	UserName   string `json:"UserName"`
	Alias      string `json:"Alias"`
	NickName   string `json:"NickName"`
	RemarkName string `json:"RemarkName"`
	Sex        int    `json:"Sex"`
	Signature  string `json:"Signature"`
	StarFriend int    `json:"StarFriend"`
	KeyWord    string `json:"KeyWord"`
}
