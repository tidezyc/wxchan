package weixin

import (
	"fmt"
	"strings"
)

type LoginInfo struct {
	Ret         int    `xml:"ret"`
	Message     string `xml:"message"`
	Skey        string `xml:"skey"`
	Wxsid       string `xml:"wxsid"`
	Wxuin       string `xml:"wxuin"`
	PassTicket  string `xml:"pass_ticket"`
	Isgrayscale int    `xml:"isgrayscale"`
	UserName    string
	DeviceID    string
	SyncKey     *SyncKey
}

func (l *LoginInfo) GetBaseRequest() *BaseRequest {
	return &BaseRequest{
		Uin:      l.Wxuin,
		Sid:      l.Wxsid,
		Skey:     l.Skey,
		DeviceID: l.DeviceID,
	}
}

type BaseRequest struct {
	Uin      string
	Sid      string
	Skey     string
	DeviceID string
}

type BaseResponse struct {
	Ret    int
	ErrMsg string
}

type Contact struct {
	UserName string
	NickName string
}

type SyncKey struct {
	Count int
	List  []struct {
		Key int
		Val int
	}
}

func (s *SyncKey) String() string {
	keys := []string{}
	for _, v := range s.List {
		keys = append(keys, fmt.Sprintf("%d_%d", v.Key, v.Val))
	}
	return strings.Join(keys, "|")
}
