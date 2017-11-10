package weixin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"

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

func (c *WeixinClient) Serve() error {
	for {
		selector, err := c.syncCheck()
		if err != nil {
			return fmt.Errorf("sync check err:%s", err)
		}
		switch selector {
		case 0, 4, 7:
			continue
		case 2:
			c.sync()
		default:
			return fmt.Errorf("sync check get invaild selector:%d", selector)
		}
	}
}

func (c *WeixinClient) syncCheck() (int, error) {
	log.Print("sync checking...")
	params := url.Values{}
	params.Add("r", fmt.Sprintf("%d", getMs()))
	params.Add("skey", c.loginInfo.Skey)
	params.Add("sid", c.loginInfo.Wxsid)
	params.Add("uin", c.loginInfo.Wxuin)
	params.Add("deviceid", c.loginInfo.DeviceID)
	params.Add("synckey", c.loginInfo.SyncKey.String())
	params.Add("_", fmt.Sprintf("%d", getMs()))
	str, err := c.httpClient.GetAsString("https://webpush.wx.qq.com/cgi-bin/mmwebwx-bin/synccheck?" + params.Encode())
	if err != nil {
		return 0, err
	}
	log.Printf("sync check response:%s", str)
	reg := regexp.MustCompile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
	ss := reg.FindStringSubmatch(str)
	if len(ss) != 3 {
		return 0, fmt.Errorf("invaild response:%s", str)
	}
	if ss[1] != "0" {
		return 0, fmt.Errorf("get ret:%s", ss[1])
	}
	selector, err := strconv.Atoi(ss[2])
	if err != nil {
		return 0, fmt.Errorf("invaild selector value:%s", ss[2])
	}
	return selector, nil
}

func (c *WeixinClient) sync() error {
	log.Print("syncing....")
	params := url.Values{}
	params.Add("skey", c.loginInfo.Skey)
	params.Add("sid", c.loginInfo.Wxsid)
	params.Add("pass_ticket", c.loginInfo.PassTicket)

	req := struct {
		BaseRequest *BaseRequest
		SyncKey     *SyncKey
		RR          int `json:"rr"`
	}{
		BaseRequest: c.loginInfo.GetBaseRequest(),
		SyncKey:     c.loginInfo.SyncKey,
		RR:          ^int(time.Now().Unix()) + 1,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	var result struct {
		BaseResponse *BaseResponse
		SyncKey      *SyncKey
	}
	err = c.httpClient.PostAsJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsync?"+params.Encode(), "application/json;charset=UTF-8", bytes.NewReader(data), &result)
	if err != nil {
		return err
	}
	if result.BaseResponse.Ret != 0 {
		return errors.New(result.BaseResponse.ErrMsg)
	}
	c.loginInfo.SyncKey = result.SyncKey
	return nil
}

func (c *WeixinClient) GetContacts() ([]*Contact, error) {
	params := url.Values{}
	params.Add("r", fmt.Sprintf("%d", getMs()))
	params.Add("lang", "en_US")
	params.Add("pass_ticket", c.loginInfo.PassTicket)
	params.Add("seq", "0")
	params.Add("skey", c.loginInfo.Skey)

	var result struct {
		BaseResponse *BaseResponse
		MemberList   []*Contact
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
