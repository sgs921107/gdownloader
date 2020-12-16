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
	"github.com/go-redis/redis"
	"log"
)

type BaseSender struct {
	Client    *redis.Client
	UrlsQueue string
	ReqsQueue string
}

func (s *BaseSender) AddURL(url string) {
	s.Client.RPush(s.UrlsQueue, url)
}

func (s *BaseSender) AddRequest(req *Request) {
	sr, err := req.Marshal()
	if err != nil {
		log.Print("Serialize request failed: req: %v, err msg: %s", req, err.Error())
		return
	}
	s.Client.RPush(s.ReqsQueue, sr)
}

func NewSender(client *redis.Client, urlsQueue string, reqsQueue string) Sender {
	var sender Sender
	sender = &BaseSender{
		Client:    client,
		UrlsQueue: urlsQueue,
		ReqsQueue: reqsQueue,
	}
	return sender
}
