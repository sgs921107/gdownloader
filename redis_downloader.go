/*************************************************************************
	> File Name: redis_downloader.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月09日 星期三 10时41分41秒
 ************************************************************************/
package gdownloader

import (
	"github.com/go-redis/redis"
	"github.com/sgs921107/gcommon"
	"github.com/sgs921107/gspider"
)

// 基于redis的分布式下载器
type RedisDownloader struct {
	BaseDownloader
	Client *redis.Client
	// RedisKey string
	Topic string
}

// 存储方法
func (d *RedisDownloader) Save(item DownloaderItem) {
	data, err := item.ToJson()
	if err != nil {
		d.Logger.Printf("serialize item failed: %s", err.Error())
		return
	}
	// 如果指定了存储的topic则存入指定的topic
	if topic, ok := item.Ctx["Topic"].(string); ok {
		d.Client.RPush(topic, string(data))
	} else {
		d.Client.RPush(d.Topic, string(data))
	}
}

// response钩子, 用于解析并存储每个请求的内容
func (d *RedisDownloader) OnResponse(response *gspider.Response) {
	item := d.Parse(response)
	item.Ctx["saveTime"] = gcommon.TimeStamp(1)
	d.Save(item)
}

// 实例化一个分布式下载器
func NewRedisDownloader(settings *DownloaderSettings) *RedisDownloader {
	spider := gspider.NewRedisSpider(settings.RedisKey, &settings.SpiderSettings)
	prefix := settings.RedisPrefix
	rd := &RedisDownloader{
		BaseDownloader: BaseDownloader{
			Spider: spider,
			Logger: spider.Logger,
		},
		Client: spider.Client,
		// RedisKey: prefix + ":" + settings.RedisKey,
		Topic: prefix + ":" + settings.Topic,
	}
	return rd
}
