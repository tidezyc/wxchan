package weixin

import (
	"fmt"
	"net/url"

	"github.com/tidezyc/httpclient"
)

type WeixinClient struct {
	httpClient *httpclient.Client
	loginInfo  *LoginInfo
}

func NewweixinClient() *WeixinClient {
	return &WeixinClient{
		httpClient: NewHttpClient(true),
	}
}

func (c *WeixinClient) GetContacts() ([]*Contact, error) {
	params := url.Values{}
	params.Add("r", fmt.Sprintf("%d", getMs()))
	params.Add("lang", "en_US")
	params.Add("pass_ticket", c.loginInfo.PassTicket)
	params.Add("seq", "0")
	params.Add("skey", c.loginInfo.Skey)

	var result struct {
		BaseResponse struct {
			Ret    int
			ErrMsg string
		}
		MemberList []*Contact
	}
	err := c.httpClient.GetAsJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?"+params.Encode(), &result)
	if err != nil {
		return nil, err
	}
	if result.BaseResponse.Ret != 0 {
		return nil, fmt.Errorf("get result:%d,errmsg:%s", result.BaseResponse.Ret, result.BaseResponse.ErrMsg)
	}
	return result.MemberList, nil
}
