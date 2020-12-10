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
	"github.com/go-redis/redis"
	"github.com/sgs921107/gdownloader"
	"github.com/sgs921107/gdownloader/send"
)

func main() {
	var params = map[string]string{
		"langx": "zh-cn",
		"gid":   "4584749",
		"ltype": "4",
		"date":  "2020-12-09",
	}

	var req = &send.Request{
		URL:    "http://www.example.com/app/member/get_game_allbets.php",
		Method: "POST",
		Body:   params,
		Headers: map[string]string{
			"Origin": "http://www.example.com",
		},
		Ctx: map[string]interface{}{
			"Age": 18,
		},
	}
	var url = "https://www.example.com"
	client := redis.NewClient(&redis.Options{
		// 你自己的redis配置
		Addr:     gdownloader.Settings.RedisAddr,
		Password: gdownloader.Settings.RedisPassword,
		DB:       gdownloader.Settings.RedisDB,
	})
	sender := send.NewSender(client, "example:start_urls", "example:queue")
	// 添加一个链接到start urls
	sender.AddURL(url)
	// 添加一个Post请求到任务队列
	sender.AddRequest(req)
}
