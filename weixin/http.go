package weixin

import (
	"net/http"
	"net/url"
	"time"

	"github.com/tidezyc/httpclient"
)

func NewHttpClient(sessionMode bool) *httpclient.Client {
	return httpclient.NewClient(&http.Client{
		Timeout:       time.Second * 30,
		CheckRedirect: checkRedirect,
		Transport:     &weixinRoundTrip{},
		Jar:           &weixinCookieJar{make(map[string]*http.Cookie)},
	})
}

func checkRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

type weixinRoundTrip struct {
}

func (this *weixinRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36")
	req.Header.Set("Referer", "https://wx.qq.com/?&lang=en_US")
	return http.DefaultTransport.RoundTrip(req)
}

type weixinCookieJar struct {
	cookies map[string]*http.Cookie
}

func (this *weixinCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		this.cookies[cookie.Name] = cookie
	}
}

func (this *weixinCookieJar) Cookies(u *url.URL) []*http.Cookie {
	cookies := []*http.Cookie{}
	for _, cookie := range this.cookies {
		cookies = append(cookies, cookie)
	}
	return cookies
}
