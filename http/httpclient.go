package http

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	nurl "net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const ua = "User-Agent:Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"

type HttpClient struct {
	client      http.Client
	readTimeout int
}

func NewHttpClient() *HttpClient {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	jar, _ := cookiejar.New(nil)
	client := http.Client{Transport: tr, Jar: jar}
	return &HttpClient{client, 0}
}

func (i *HttpClient) GetResp(url string) (resp *http.Response, err error) {
	rand.Seed(time.Now().Unix())
	r := rand.Int()
	if strings.Contains(url, "?") {
		url += "&"
	} else {
		url += "?"
	}
	url += "_=" + strconv.Itoa(r)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err == nil {
		if err == nil {
			resp, err = i.Do(req)
		}
	}
	return
}

func (i *HttpClient) GetBytes(url string) (bytes []byte, err error) {
	resp, err := i.GetResp(url)
	if err == nil {
		defer resp.Body.Close()
		bytes, err = ioutil.ReadAll(resp.Body)
	}
	return
}

func (i *HttpClient) GetString(url string) (str string, err error) {
	bytes, err := i.GetBytes(url)
	if err == nil {
		str = string(bytes)
	}
	return
}

func (i *HttpClient) PostResp(url string, data []byte) (resp *http.Response, err error) {
	rand.Seed(time.Now().Unix())
	r := rand.Int()
	if strings.Contains(url, "?") {
		url += "&"
	} else {
		url += "?"
	}
	url += "_=" + strconv.Itoa(r)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err == nil {
		if err == nil {
			resp, err = i.Do(req)
		}
	}
	return
}

func (i HttpClient) PostBytes(url string, data []byte) (bytes []byte, err error) {
	resp, err := i.PostResp(url, data)
	if err == nil {
		defer resp.Body.Close()
		bytes, err = ioutil.ReadAll(resp.Body)
	}
	return
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

func (i *HttpClient) Do(req *http.Request) (resp *http.Response, err error) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", ua)
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
			i.client.Transport.(*http.Transport).Proxy = func(*http.Request) (*nurl.URL, error) {
				return u, nil
			}
		} else {
			i.client.Transport.(*http.Transport).Proxy = nil
		}
	}
}

func (i *HttpClient) SetReadBodyTimeout(timeout int) {
	if timeout <= 0 {
		timeout = 0
	}
	i.readTimeout = timeout
}

func (i *HttpClient) ReadBodyWithTimeout(resp *http.Response) (data []byte, err error) {
	if resp == nil {
		return nil, errors.New("resp is nil.")
	}
	ch := make(chan bool, 0)
	buf := make([]byte, bytes.MinRead)
	var t int
	timer := time.NewTimer(time.Second * time.Duration(i.readTimeout))
	go func() {
		t, err = resp.Body.Read(buf)
		data = buf[:t]
		ch <- true
	}()
	if i.readTimeout > 0 {
		select {
		case <-ch:
		case <-timer.C:
			err = errors.New("readbody timeout.")
		}
	} else {
		<-ch
	}
	timer.Stop()
	return
}

func init() {
	go func() {
		for {
			time.Sleep(time.Minute * 2)
			debug.FreeOSMemory()
		}
	}()
}
