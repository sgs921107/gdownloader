/*************************************************************************
	> File Name: bin/start.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月09日 星期三 14时51分35秒
 ************************************************************************/

package main

import (
	"github.com/sgs921107/gdownloader"
)

func main() {
	settings := gdownloader.NewDownloaderSettings("env_demo")
	// 生成一个redis downloader实例
	rd := gdownloader.NewRedisDownloader(settings)
	rd.Spider.Client.RPush(rd.Spider.RedisKey, "http://www.baidu.com")
	rd.Run()
}
