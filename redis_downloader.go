/*************************************************************************
	> File Name: redis_downloader.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月09日 星期三 10时41分41秒
 ************************************************************************/

package gdownloader

import (
	"github.com/sgs921107/gcommon"
	"github.com/sgs921107/gredis"
	"github.com/sgs921107/gspider"
)

// RedisDownloader 基于redis的分布式下载器
type RedisDownloader struct {
	BaseDownloader
	Client *gredis.Client
}

// Save 存储方法
func (d *RedisDownloader) save(item *DownloaderItem) {
	data, err := item.ToJSON()
	if err != nil {
		d.logger.Errorw("Serialize Item Failed",
			"errMsg", err.Error(),
		)
		return
	}
	// 将解析出的数据存入指定的topic, 如果未指定则不存储
	topic, ok := item.Ctx["topic"].(string)
	if !ok {
		prefix := d.settings.Redis.Prefix
		if host, err := gcommon.FetchURLHost(item.URL); err != nil {
			d.logger.Errorw("Save Item Error: fetch url host failed",
				"errMsg", err.Error(),
				"url", item.URL,
			)
			return
		} else {
			topic = prefix + ":items:" + host
		}
	}
	if size := d.settings.Downloader.MaxTopicSize; size > 0 {
		d.Client.RPushTrim(topic, size, string(data))
	} else {
		d.Client.RPush(topic, string(data))
	}
}

// NewRedisDownloader 实例化一个分布式下载器
func NewRedisDownloader(settings *Settings) Downloader {
	spiderSettings := settings.SpiderSettings
	spider := gspider.NewRedisSpider(spiderSettings)
	rd := &RedisDownloader{
		BaseDownloader: BaseDownloader{
			spider:   spider,
			logger:   spider.Logger,
			settings: settings,
		},
		Client: spider.Client,
	}
	rd.BaseDownloader.save = rd.save
	return rd
}
