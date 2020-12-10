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
	// 生成一个redis downloader实例
	rd := gdownloader.NewRedisDownloader(&gdownloader.Settings)
	rd.Spider.OnRequest(rd.AddDownloadTime)
	rd.Spider.OnResponse(rd.AddDownloadedTime)
	rd.Spider.OnResponse(rd.OnResponse)
	rd.Spider.Start()
}
