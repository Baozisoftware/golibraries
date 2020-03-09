package http

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/Baozisoftware/golibraries/utils"
	"html"
	"io/ioutil"
	"net/http"
	nurl "net/url"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type cfBypass struct {
	*HttpClient
	url          string
	__cfduid     string
	cf_clearance string
	ua           string
}

func NewCFBypass(url, proxy string) *cfBypass {
	obj := new(cfBypass)
	obj.HttpClient = NewHttpClient()
	parsedUrl, err := nurl.Parse(url)
	if err == nil {
		parsedUrl.Scheme = "http"
		obj.url = parsedUrl.String()
		if proxy != "" {
			obj.SetProxy(proxy)
		}
		for i := 0; i < 5; i++ {
			obj.genUA()
			if obj.Bypass() {
				return obj
			}
		}
	}
	return nil
}

func (i *cfBypass) Bypass() bool {
	req, err := NewGetRequest(i.url)
	if err == nil {
		resp, err := i.do(req, "")
		if err == nil && resp.StatusCode == 503 && resp.Header.Get("Server") == "cloudflare" {
			data, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				body := string(data)
				parsedUrl, err := nurl.Parse(i.url)
				if err == nil {
					domain := parsedUrl.Hostname()
					reg, _ := regexp.Compile(`<form id="challenge-form" action="(\S+)" method="POST" enctype="\S+">`)
					arr := reg.FindStringSubmatch(body)
					if len(arr) > 1 {
						action := arr[1]
						action = html.UnescapeString(action)
						reg, _ = regexp.Compile(`input type="hidden" name="r" value="(\S+)"/?>`)
						arr = reg.FindStringSubmatch(body)
						if len(arr) > 1 {
							r := arr[1]
							reg, _ = regexp.Compile(`input type="hidden" name="jschl_vc" value="(\S+)"/?>`)
							arr = reg.FindStringSubmatch(body)
							if len(arr) > 1 {
								jschl_vc := arr[1]
								reg, _ = regexp.Compile(`input type="hidden" name="pass" value="(\S+)"/?>`)
								arr = reg.FindStringSubmatch(body)
								if len(arr) > 1 {
									pass := arr[1]
									reg, _ = regexp.Compile(`k = '(\S+)';`)
									arr = reg.FindStringSubmatch(body)
									if len(arr) > 1 {
										k := arr[1]
										reg, _ = regexp.Compile(fmt.Sprintf(`<div style="display:none;visibility:hidden;" id="%s">(\S+)</div>`, k))
										arr = reg.FindStringSubmatch(body)
										if len(arr) > 1 {
											innerHTML := arr[1]
											reg, _ = regexp.Compile(`setTimeout\(function\(\)\{\n([\S\s]+'; \d+')\n[\S\s]+}, (\d+)`)
											arr = reg.FindStringSubmatch(body)
											if len(arr) > 2 {
												challenge, _ms := arr[1], arr[2]
												challenge = fmt.Sprintf(`
var document = {
    createElement: function () {
        return { firstChild: { href: "http://%s/" } }
    },
    getElementById: function () {
        return { "innerHTML": "%s" };
    }
};
%s; a.value
    `, domain, innerHTML, challenge)
												challenge = base64.StdEncoding.EncodeToString([]byte(challenge))
												js := fmt.Sprintf(`
var atob = Object.setPrototypeOf(function (str) {
    try {
        return Buffer.from("" + str, "base64").toString("binary");
    } catch (e) { }
}, null);
var challenge = atob("%s");
var context = Object.setPrototypeOf({ atob: atob }, null);
var options = {
    filename: "iuam-challenge.js",
    contextOrigin: "cloudflare:iuam-challenge.js",
    contextCodeGeneration: { strings: true, wasm: false },
    timeout: 5000
};
process.stdout.write(String(
    require("vm").runInNewContext(challenge, context, options)
));
    `, challenge)
												proc := exec.Command("node", "-e", js)
												result := bytes.NewBufferString("")
												proc.Stdout = result
												if proc.Start() == nil {
													state, err := proc.Process.Wait()
													if err == nil && state.Success() {
														ms, err := strconv.Atoi(_ms)
														if err == nil {
															time.Sleep(time.Duration(ms) * time.Millisecond)
															jschl_answer := result.String()
															_url := fmt.Sprintf("http://%s/%s", domain, action)
															r := nurl.QueryEscape(r)
															jschl_vc = nurl.QueryEscape(jschl_vc)
															body := fmt.Sprintf("r=%s&jschl_vc=%s&pass=%s&jschl_answer=%s", r, jschl_vc, pass, jschl_answer)
															data = []byte(body)
															req, err = NewPostRequest(_url, bytes.NewReader(data))
															if err == nil {
																resp, err = i.do(req, i.url)
																if err == nil && resp.StatusCode != 503 {
																	i.__cfduid = i.GetCookie(i.url, "__cfduid")
																	i.cf_clearance = i.GetCookie(i.url, "cf_clearance")
																	return i.__cfduid != "" && i.cf_clearance != ""
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		} else if err == nil {
			return true
		}
	}
	return false
}

func (i *cfBypass) GetTokens() []string {
	cookies := fmt.Sprintf("%s=%s; %s=%s", "__cfduid", i.__cfduid, "cf_clearance", i.cf_clearance)
	return []string{cookies, i.ua}
}

func (i *cfBypass) genUA() {
	list := []string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
	}
	x := utils.GetRandomIntN(len(list))
	i.ua = list[x]
	i.HttpClient.SetUserAgent(i.ua)
}

func (i *cfBypass) do(req *http.Request, referer string) (*http.Response, error) {
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
	if req.Method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return i.HttpClient.Do(req)
}

func (i *cfBypass) Do(req *http.Request) (*http.Response, error) {
	resp, err := i.HttpClient.Do(req)
	if err == nil && resp.StatusCode == 503 && resp.Header.Get("Server") == "cloudflare" {
		for x := 0; x < 5; x++ {
			if i.Bypass() {
				return i.HttpClient.Do(req)
			}
		}
	}
	return resp, err
}