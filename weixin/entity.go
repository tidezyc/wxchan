package weixin

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

type Contact struct {
	UserName string
	NickName string
}
