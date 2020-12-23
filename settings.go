/*************************************************************************
	> File Name: settings.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月09日 星期三 09时49分57秒
 ************************************************************************/

package gdownloader

import (
	"fmt"
	"time"
	"reflect"

	"github.com/sgs921107/gspider"
)

// SpiderSettings spider settings type
type SpiderSettings = gspider.SpiderSettings

// DownloaderSettings downloader的配置结构
// 使数据结构简单，不继承自spider settings, 通过反射来生成spdier settings
type DownloaderSettings struct {
	// SpiderSettings
	//------------------------------------------------------------------
	Debug       	bool
	LogLevel		string
	LogFile     	string
	RotationTime	time.Duration
	RotationMaxAge	time.Duration
	LogFlag        	int
	FlushOnStart   	bool // 开始前清空之前的数据
	UserAgent      	string
	ConcurrentReqs 	int  // 并发
	MaxDepth       	int  // 最大深度
	DontFilter     	bool // 不过滤
	EnableCookies  	bool // 启用cookies
	Async          	bool // 启用异步
	KeepAlive      	bool
	Timeout        	time.Duration
	MaxConns       	int
	// 以下使用redis spider时需要配置
	RedisAddr      	string
	RedisDB        	int
	RedisPassword  	string
	RedisPrefix    	string
	MaxIdleTimeout 	time.Duration // 最大闲置时间, redis spider使用 0表示一直运行
	//------------------------------------------------------------------
	RedisKey 		string
	// 存储页面数据的最大数量  list元素超出将被裁剪, 避免内存过高
	MaxTopicSize 	int64
}

// 通过反射生成spidersettings
func (s DownloaderSettings) createSpiderSettings() *SpiderSettings {
	spiderSettings := &SpiderSettings{}
	dsv := reflect.ValueOf(s)
	sst := reflect.TypeOf(spiderSettings).Elem()
	ssv := reflect.ValueOf(spiderSettings).Elem()
	for i := 0; i < sst.NumField(); i++ {
		field := sst.Field(i)
		name := field.Name
		val := dsv.FieldByName(name)
		// 如果值是无效的则跳过
		if !val.IsValid() {
			continue
		}
		switch field.Type.Name() {
		case "string":
			ssv.FieldByName(name).SetString(val.String())
		case "int":
			ssv.FieldByName(name).SetInt(val.Int())
		case "bool":
			ssv.FieldByName(name).SetBool(val.Bool())
		case "Duration":
			ssv.FieldByName(name).SetInt(val.Int())
		default:
			fmt.Printf("Warning: miss a option: %s, val: %v", name, val)
		}
	}
	return spiderSettings
}

// SettingsDemo 配置实例demo
var SettingsDemo = DownloaderSettings{
	RedisKey:     "start_urls",
	MaxTopicSize: 50000,
	Debug:        false,
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
	Timeout: time.Second * 5,
	// 最大连接数
	MaxConns: 100,
	// 如果不为""则使用redis存储数据
	RedisAddr:      "172.17.0.1:6379",
	RedisDB:        2,
	RedisPassword:  "qaz123",
	RedisPrefix:    "simple",
	// 空闲超时
	MaxIdleTimeout: time.Second * 10,
}
