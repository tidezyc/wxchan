package weixin

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/tuotoo/qrcode"
)

func (c *WeixinClient) Login() error {
	uuid, err := c.getQRuuid()
	if err != nil {
		return fmt.Errorf("get uuid err:%s", err)
	}
	if uuid == "" {
		return errors.New("empty uuid")
	}
	log.Printf("get uuid:%s", uuid)
	err = c.getQR(uuid)
	if err != nil {
		return fmt.Errorf("get qrcode err:%s", err)
	}
	loginInfo, err := c.checkLogin(uuid)
	if err != nil {
		return fmt.Errorf("check login err:%s", err)
	}
	devideId := "e"
	for i := 0; i < 15; i++ {
		devideId += strconv.Itoa(rand.Intn(10))
	}
	loginInfo.DeviceID = devideId
	log.Printf("check login login info:%s,err:%s", loginInfo, err)
	c.loginInfo = loginInfo
	err = c.webInit()
	if err != nil {
		return fmt.Errorf("web init err:%s", err)
	}
	err = c.showMobileLogin()
	if err != nil {
		return fmt.Errorf("show mobile login err:%s", err)
	}
	return nil
}

func (c *WeixinClient) getQRuuid() (string, error) {
	params := url.Values{}
	params.Add("appid", "wx782c26e4c19acffb")
	params.Add("fun", "new")
	params.Add("lang", "en_US")
	params.Add("redirect_uri", "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage")
	params.Add("_", fmt.Sprintf("%d", getMs()))
	str, err := c.httpClient.GetAsString("https://login.wx.qq.com/jslogin?" + params.Encode())
	if err != nil {
		return "", err
	}
	log.Printf("jslogin get response:%s", str)
	reg := regexp.MustCompile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "(\S+?)";`)
	ss := reg.FindStringSubmatch(str)
	if len(ss) != 3 {
		return "", fmt.Errorf("invaild response:%s", str)
	}
	if ss[1] != "200" {
		return "", errors.New("get result code:" + ss[0])
	}
	return ss[2], nil
}

func (c *WeixinClient) getQR(uuid string) error {
	req, err := http.NewRequest("GET", "https://login.weixin.qq.com/qrcode/"+uuid, nil)
	if err != nil {
		return err
	}
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	qrmatrix, err := qrcode.Decode(rsp.Body)
	if err != nil {
		return fmt.Errorf("decode qrcode err:%s", err)
	}
	return PrintQrcode(qrmatrix)
}

func (c *WeixinClient) checkLogin(uuid string) (*LoginInfo, error) {
	tip := 1
	for {
		params := url.Values{}
		params.Add("loginicon", "true")
		params.Add("uuid", uuid)
		params.Add("tip", fmt.Sprintf("%d", tip))
		params.Add("_", fmt.Sprintf("%d", getMs()))
		params.Add("r", getR())
		str, err := c.httpClient.GetAsString("https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login?" + params.Encode())
		if err != nil {
			return nil, err
		}
		reg := regexp.MustCompile(`window.code=(\d+)`)
		ss := reg.FindStringSubmatch(str)
		if len(ss) != 2 {
			return nil, fmt.Errorf("invaild response:%s", str)
		}
		status, err := strconv.Atoi(ss[1])
		if err != nil {
			return nil, err
		}
		switch status {
		case 200:
			info, err := c.getLoginInfo(str)
			if err != nil {
				return nil, fmt.Errorf("get login info err:%s", err)
			}
			return info, nil
		case 201:
			tip = 0
			log.Println("confirm on phone")
			continue
		case 408:
			continue
		default:
			return nil, fmt.Errorf("check login get invaild status:%d", status)
		}
	}
	return nil, errors.New("check login loop exit")

}

func (c *WeixinClient) getLoginInfo(str string) (*LoginInfo, error) {
	reg := regexp.MustCompile(`window.redirect_uri="(\S+)";`)
	ss := reg.FindStringSubmatch(str)
	if len(ss) != 2 {
		return nil, fmt.Errorf("invaild response:%s", str)
	}
	redirectUri := ss[1]
	log.Printf("call redirect uri:%s", redirectUri)
	resp, err := c.httpClient.Get(redirectUri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var loginInfo LoginInfo
	d := xml.NewDecoder(resp.Body)
	err = d.Decode(&loginInfo)
	if err != nil {
		return nil, err
	}
	if loginInfo.PassTicket == "" || loginInfo.Skey == "" || loginInfo.Wxsid == "" || loginInfo.Wxuin == "" {
		return nil, fmt.Errorf("get invaild login info:%s", loginInfo)
	}
	return &loginInfo, nil
}

func (c *WeixinClient) webInit() error {
	params := url.Values{}
	params.Add("r", getR())
	params.Add("lang", "en_US")
	params.Add("pass_ticket", c.loginInfo.PassTicket)

	var request = &struct {
		BaseRequest *BaseRequest
	}{c.loginInfo.GetBaseRequest()}

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	var result struct {
		BaseResponse struct {
			Ret    int
			ErrMsg string
		}
		User struct {
			UserName string
		}
	}
	err = c.httpClient.PostAsJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxinit?"+params.Encode(), "application/json;charset=UTF-8", bytes.NewReader(data), &result)
	if err != nil {
		return err
	}
	if result.BaseResponse.Ret != 0 {
		return fmt.Errorf("get result:%d,errmsg:%s", result.BaseResponse.Ret, result.BaseResponse.ErrMsg)
	}
	c.loginInfo.UserName = result.User.UserName
	return nil
}

func (c *WeixinClient) showMobileLogin() error {
	params := url.Values{}
	params.Add("lang", "en_US")
	params.Add("pass_ticket", c.loginInfo.PassTicket)

	var request = &struct {
		BaseRequest  *BaseRequest
		Code         int
		FromUserName string
		ToUserName   string
		ClientMsgId  int64
	}{
		BaseRequest:  c.loginInfo.GetBaseRequest(),
		Code:         3,
		FromUserName: c.loginInfo.UserName,
		ToUserName:   c.loginInfo.UserName,
		ClientMsgId:  getMs(),
	}

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	var result struct {
		BaseResponse struct {
			Ret    int
			ErrMsg string
		}
	}
	err = c.httpClient.PostAsJson("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify?"+params.Encode(), "application/json;charset=UTF-8", bytes.NewReader(data), &result)
	if err != nil {
		return err
	}
	if result.BaseResponse.Ret != 0 {
		return fmt.Errorf("request result:%d,errmsg:%s", result.BaseResponse.Ret, result.BaseResponse.ErrMsg)
	}
	return nil
}
