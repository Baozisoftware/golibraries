package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	nurl "net/url"
	"strconv"
	"strings"
	"time"
)

const ua = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"

type HttpClient struct {
	client   http.Client
	_ua      string
	_referer string
}

func NewHttpClient() *HttpClient {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	jar, _ := cookiejar.New(nil)
	client := http.Client{Transport: tr, Jar: jar}
	return &HttpClient{client, "", ""}
}

func (i *HttpClient) GetResp(url string) (resp *http.Response, err error) {
	req, err := NewGetRequest(url)
	if err == nil {
		resp, err = i.Do(req)
	}
	return
}

func (i *HttpClient) GetBytes(url string) (bytes []byte, err error) {
	resp, err := i.GetResp(url)
	return ReadRespBytes(resp)
}

func (i *HttpClient) GetString(url string) (str string, err error) {
	bytes, err := i.GetBytes(url)
	if err == nil {
		str = string(bytes)
	}
	return
}

func (i *HttpClient) PostResp(url string, data []byte) (resp *http.Response, err error) {
	req, err := NewPostRequest(url, bytes.NewReader(data))
	if err == nil {
		resp, err = i.Do(req)
	}
	return
}

func (i HttpClient) PostBytes(url string, data []byte) (bytes []byte, err error) {
	resp, err := i.PostResp(url, data)
	return ReadRespBytes(resp)
}

func (i *HttpClient) PostString(url, data string) (str string, err error) {
	bytes, err := i.PostBytes(url, []byte(data))
	if err == nil {
		str = string(bytes)
	}
	return
}

func (i *HttpClient) GetCookies(url string) (cookies map[string]string, err error) {
	u, err := nurl.Parse(url)
	if err == nil {
		tc := i.client.Jar.Cookies(u)
		cookies = make(map[string]string)
		for _, v := range tc {
			cookies[v.Name] = v.Value
		}
	}
	return
}

func (i *HttpClient) GetCookie(url, name string) string {
	cookies, err := i.GetCookies(url)
	if err == nil {
		if v, ok := cookies[name]; ok {
			return v
		}
	}
	return ""
}

func (i *HttpClient) SetCookies(url, cookies string) bool {
	u, err := nurl.Parse(url)
	if err != nil {
		return false
	}
	t := make([]*http.Cookie, 0)
	sep := ";"
	if strings.Contains(cookies, "; ") {
		sep = "; "
	}
	cks := strings.Split(cookies, sep)
	for _, c := range cks {
		v := strings.Split(c, "=")
		if len(v) < 2 {
			continue
		}
		t = append(t, &http.Cookie{Name: v[0], Value: strings.Join(v[1:], ""), Expires: time.Now().AddDate(1, 0, 0), Path: "/"})
	}
	i.client.Jar.SetCookies(u, t)
	return true
}

type cookiesJson struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (i *HttpClient) SetCookiesByJson(url, data string) bool {
	bytes := []byte(data)
	list := make([]cookiesJson, 0)
	if json.Unmarshal(bytes, &list) != nil {
		return false
	}
	str := ""
	for _, v := range list {
		str += fmt.Sprintf("%s=%s; ", v.Name, v.Value)
	}
	return i.SetCookies(url, str)
}

func (i *HttpClient) SetCookie(url, name, value string) bool {
	return i.SetCookies(url, fmt.Sprintf("%s=%s", name, value))
}

func (i *HttpClient) ClearCookie() {
	jar, _ := cookiejar.New(nil)
	i.client.Jar = jar
}

func (i *HttpClient) Do(req *http.Request) (resp *http.Response, err error) {
	if req.Header.Get("User-Agent") == "" {
		if i._ua == "" {
			req.Header.Set("User-Agent", ua)
		} else {
			req.Header.Set("User-Agent", i._ua)
		}
		if req.Header.Get("Referer") == "" && i._referer != "" {
			req.Header.Set("Referer", i._referer)
		}
	}
	if req.Method == http.MethodPost && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
	resp, err = i.client.Do(req)
	return
}

func (i *HttpClient) SetTimeout(timeout int) {
	if timeout <= 0 {
		timeout = 0
	}
	i.client.Timeout = time.Second * time.Duration(timeout)
}

func (i *HttpClient) SetResponseHeaderTimeout(timeout int) {
	if timeout <= 0 {
		timeout = 0
	}
	i.client.Transport.(*http.Transport).ResponseHeaderTimeout = time.Second * time.Duration(timeout)
}

func (i *HttpClient) SetProxy(url string) {
	if url == "" {
		i.client.Transport.(*http.Transport).Proxy = nil
	} else {
		u, err := nurl.Parse(url)
		if err == nil {
			i.client.Transport.(*http.Transport).Proxy = http.ProxyURL(u)
		} else {
			i.client.Transport.(*http.Transport).Proxy = nil
		}
	}
}

func (i *HttpClient) SetBodyTimeout(timeout int) {
	if timeout > 0 {
		i.client.Transport.(*http.Transport).DialContext = func(ctx context.Context, netw, addr string) (net.Conn, error) {
			tot := time.Second * time.Duration(timeout)
			conn, err := net.DialTimeout(netw, addr, tot)
			if err != nil {
				return nil, err
			}
			return newTimeoutConn(conn, tot), nil
		}
	} else {
		i.client.Transport.(*http.Transport).DialContext = func(ctx context.Context, netw, addr string) (net.Conn, error) {
			return net.Dial(netw, addr)
		}
	}
}

func (i *HttpClient) GetCookiesString(url string) string {
	cookies, err := i.GetCookies(url)
	if err == nil {
		ret := ""
		for k, v := range cookies {
			ret += fmt.Sprintf("%s=%s; ", k, v)
		}
		return ret
	}
	return ""
}

func NewGetRequest(url string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, url, nil)
}

func NewPostRequest(url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(http.MethodPost, url, body)
}

func AppendUrlRandom(url string) string {
	rand.Seed(time.Now().Unix())
	r := rand.Int()
	if strings.Contains(url, "?") {
		url += "&"
	} else {
		url += "?"
	}
	url += "_=" + strconv.Itoa(r)
	return url
}

func ReadRespBytes(resp *http.Response) (bytes []byte, err error) {
	if resp != nil {
		defer resp.Body.Close()
		bytes, err = ioutil.ReadAll(resp.Body)
	}
	return
}

func GetCookieValueByCookesHeader(cookies, name string) string {
	sep := ";"
	if strings.Contains(cookies, "; ") {
		sep = "; "
	}
	cks := strings.Split(cookies, sep)
	for _, c := range cks {
		v := strings.Split(c, "=")
		if len(v) < 2 {
			continue
		}
		if v[0] == name {
			return strings.Join(v[1:], "")
		}
	}
	return ""
}

func (i *HttpClient) SetUserAgent(ua string) {
	i._ua = ua
}

func (i *HttpClient) SetReferer(referer string) {
	i._referer = referer
}
