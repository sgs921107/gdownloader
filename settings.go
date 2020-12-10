/*************************************************************************
	> File Name: settings.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月09日 星期三 09时49分57秒
 ************************************************************************/

package gdownloader

import (
	"github.com/sgs921107/gspider"
)

// 起别名
type SpiderSettings = gspider.SpiderSettings

// downloader的配置结构
type DownloaderSettings struct {
	SpiderSettings
	RedisKey string
	Topic    string
}

// 配置实例demo
var Settings = DownloaderSettings{
	Topic:    "items",
	RedisKey: "start_urls",
	SpiderSettings: SpiderSettings{
		Debug: false,
		// 是否在启动前清空之前的数据
		FlushOnStart: false,
		// UserAgent bool
		UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
		ConcurrentReqs: 16,
		// 最大深度
		MaxDepth: 1,
		// 允许重复抓取
		DontFilter: true,
		// 启用异步
		Async:         true,
		EnableCookies: false,
		// 是否开启长连接 bool
		KeepAlive: true,
		// 超时  单位：秒
		Timeout: 5,
		// 最大连接数
		MaxConns: 100,
		// 空闲超时 单位: 秒
		// 如果不为""则使用redis存储数据
		RedisAddr:      "172.17.0.1:6379",
		RedisDB:        2,
		RedisPassword:  "qaz123",
		RedisPrefix:    "simple",
		MaxIdleTimeout: 10,
	},
}
