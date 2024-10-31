/*************************************************************************
	> File Name: send.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月10日 星期四 11时25分34秒
*************************************************************************/
/*
向downloader监听的队列投放任务的例子
*/

package main

import (
	"bytes"
	"github.com/go-redis/redis"
	"github.com/sgs921107/gdownloader"
	"github.com/sgs921107/gdownloader/send"
)

func main() {

	settings, _ := gdownloader.NewSettingsFromEnvFile("/etc/gdownloader/.env")
	body := `{"invoke_info":{"pos_1":[{}],"pos_2":[{}],"pos_3":[{}]}}`
	var req = &send.Request{
		URL:    "https://ug.baidu.com/mcp/pc/pcsearch",
		Method: "POST",
		Body:   bytes.NewBufferString(body).Bytes(),
		Headers: map[string]string{
			"Origin":       "https://www.baidu.com",
			"Referer":      "https://www.baidu.com",
			"Content-Type": "application/json",
		},
		Ctx: map[string]interface{}{
			"clearHead":    true,
			"gzipCompress": true,
		},
	}
	var url = "https://translate.google.com"
	client := redis.NewClient(&redis.Options{
		// 你自己的redis配置
		Addr:     settings.Redis.Addr,
		Password: settings.Redis.Password,
		DB:       settings.Redis.DB,
	})
	prefix := settings.Redis.Prefix
	sender := send.NewSender(client, prefix+":start_urls", prefix+":queue", prefix+":items")
	// 添加一个链接到start urls
	sender.AddURL(url)
	// 添加一个Post请求到任务队列
	sender.AddRequest(req)
}
