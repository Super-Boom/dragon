package httpclient

import (
	"dragon/core/dragon/dlogger"
	"dragon/core/dragon/tracker"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// DefaultClient default http client pool
var DefaultClient = NewClient(nil)

type Client struct {
	TrackWriter *http.Request
	HttpCli     *http.Client
}

// NewClient new a client
func NewClient(req *http.Request) *Client {
	trans := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxConnsPerHost:       100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// 新建一个http的client
	httpCli := &http.Client{
		Transport:     trans,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       10 * time.Second,
	}
	return &Client{
		TrackWriter: req,
		HttpCli:     httpCli,
	}
}

// Response Client response struct
type Response struct {
	Content string
	Status  int
	Err     error
}

// GET send get request
func (c *Client) GET(url string, params map[string]string, headers map[string]string) *Response {
	return c.send(url, params, "GET", headers)
}

// POST send post request
func (c *Client) POST(url string, params map[string]string, headers map[string]string) *Response {
	return c.send(url, params, "POST", headers)
}

// PUT send put request
func (c *Client) PUT(url string, params map[string]string, headers map[string]string) *Response {
	return c.send(url, params, "PUT", headers)
}

// DELETE send delete request
func (c *Client) DELETE(url string, params map[string]string, headers map[string]string) *Response {
	return c.send(url, params, "DELETE", headers)
}

// PATCH send patch request
func (c *Client) PATCH(url string, params map[string]string, headers map[string]string) *Response {
	return c.send(url, params, "PATCH", headers)
}

func (c *Client) send(url string, params map[string]string, method string, headers map[string]string) (resp *Response) {
	defer func() {
		dlogger.Info("httpclient:"+method, map[string]interface{}{
			"url":     url,
			"params":  params,
			"headers": headers,
		})
	}()
	// 跟踪器
	var trackMan *tracker.Tracker
	if c.TrackWriter != nil {
		trackInfo := c.TrackWriter.Header.Get(tracker.TrackKey)
		trackMan = tracker.UnMarshal(trackInfo)
	}

	paramsStr := ""
	for k, v := range params {
		paramsStr += k + "=" + v + "&"
	}
	if paramsStr != "" {
		paramsStr = paramsStr[:len(paramsStr)-1]
	}
	var req *http.Request
	if method == "GET" {
		if paramsStr != "" {
			url += "?" + paramsStr
		}
		req, _ = http.NewRequest(method, url, nil)
	} else {
		req, _ = http.NewRequest(method, url, strings.NewReader(paramsStr))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	// add reqdata headers
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if trackMan != nil {
		//trackMan.Service.Req = reqdata todo req直接结构体不行
		trackMan.HttpClient.Req.Uri = req.URL.String()
		trackMan.HttpClient.Req.Params = params // 记录请求内容
	}

	rsp, err := c.HttpCli.Do(req)
	if err != nil {
		log.Println(err)
		resp = &Response{
			"",
			http.StatusInternalServerError,
			err,
		}
		if trackMan != nil {
			trackMan.ErrInfo = err
			c.TrackWriter.Header.Set(tracker.TrackKey, trackMan.Marshal()) // 最后写日志跟踪
		} // todo 继续完成
		return
	}
	defer rsp.Body.Close()

	content, errR := ioutil.ReadAll(rsp.Body)
	contentStr := string(content)
	if trackMan != nil {
		trackMan.HttpClient.Resp = contentStr
	}

	if errR != nil {
		resp = &Response{
			contentStr,
			http.StatusInternalServerError,
			errR,
		}

		if trackMan != nil {
			trackMan.ErrInfo = errR
			c.TrackWriter.Header.Set(tracker.TrackKey, trackMan.Marshal()) // 最后写日志跟踪
		}
		return
	}
	// service返回

	resp = &Response{
		contentStr,
		rsp.StatusCode,
		errR,
	}
	if trackMan != nil {
		c.TrackWriter.Header.Set(tracker.TrackKey, trackMan.Marshal()) // 最后写日志跟踪
	}
	return
}

//send postJson
func (c *Client) POSTJson(url string, paramsStr string) (resp *Response) {
	var req *http.Request
	req, _ = http.NewRequest("POST", url, strings.NewReader(paramsStr))
	req.Header.Add("Content-Type", "application/json")
	rsp, err := c.HttpCli.Do(req)

	if err != nil {
		resp = &Response{
			"",
			http.StatusInternalServerError,
			err,
		}
		return
	}
	defer rsp.Body.Close()

	content, errR := ioutil.ReadAll(rsp.Body)

	if errR != nil {
		log.Println(err)
		resp = &Response{
			string(content),
			http.StatusInternalServerError,
			errR,
		}
		return
	}

	resp = &Response{
		string(content),
		rsp.StatusCode,
		errR,
	}
	return
}
