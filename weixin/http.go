package weixin

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/tidezyc/httpclient"
)

func NewHttpClient(sessionMode bool) *httpclient.Client {
	jar, _ := cookiejar.New(nil)
	return httpclient.NewClient(&http.Client{
		Timeout:       time.Second * 30,
		CheckRedirect: checkRedirect,
		Transport:     &weixinRoundTrip{},
		Jar:           jar,
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
