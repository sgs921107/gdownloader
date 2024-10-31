/*************************************************************************
	> File Name: send.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月10日 星期四 10时30分28秒
*************************************************************************/
/*
downloader的任务只负责对start urls队列和真实的请求队列queue进行监听和下载
它没有自己的下载目标，也不负责扩链
真实的下载需求应该单独放在一个项目进行处理
假设你启动了一个downloader服务, 它正监听start urls队列simple:start_urls和任务队列simple:queue
你可以直接向simple:start_urls队列投放url(仅支持GET请求),
或者可以直接向simple:queue放入序列化后的请求

该模块不具有实际意义，仅展示该如何向downloader投放任务
*/

package send

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// SerializableRequest 序列化request结构
type SerializableRequest struct {
	URL      string
	Method   string
	Depth    int
	Body     []byte
	ID       uint32
	Ctx      map[string]interface{}
	Headers  *http.Header
	ProxyURL string
}

// Request 请求结构
type Request struct {
	URL      string
	Method   string
	Depth    int
	Body     []byte
	ID       uint32
	Ctx      map[string]interface{}
	Headers  map[string]string
	ProxyURL string
}

// 格式化headers
func (r *Request) formatHeaders() *http.Header {
	headers := make(http.Header)
	for k, v := range r.Headers {
		headers.Set(k, v)
	}
	return &headers
}

// Marshal 将请求进行序列化
func (r *Request) Marshal() ([]byte, error) {
	// 对url进行校验
	u, err := url.Parse(r.URL)
	if err != nil {
		return []byte{}, nil
	}
	headers := r.formatHeaders()
	sr := &SerializableRequest{
		URL:      u.String(),
		Method:   r.Method,
		Depth:    r.Depth,
		Body:     r.Body,
		ID:       r.ID,
		Ctx:      r.Ctx,
		Headers:  headers,
		ProxyURL: r.ProxyURL,
	}
	return json.Marshal(sr)
}
