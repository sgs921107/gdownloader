/*************************************************************************
	> File Name: redis_downloader.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月08日 星期二 16时28分05秒
 ************************************************************************/

package gdownloader

import (
	"github.com/sgs921107/gcommon"
	"github.com/sgs921107/gspider"
	"log"
)

// 上下文的类型
type CtxMap map[string]interface{}

// 将response中的上下文转为自定义的上下文结构
func CtxToMap(ctx *gspider.Context) CtxMap {
	data := make(CtxMap)
	ctx.ForEach(func(k string, v interface{}) interface{} {
		data[k] = v
		return nil
	})
	return data
}

// 定义下载器的机构
type BaseDownloader struct {
	Spider   *gspider.RedisSpider
	Logger   *log.Logger
	settings *DownloaderSettings
}

// 解析方法
func (d *BaseDownloader) Parse(response *gspider.Response) DownloaderItem {
	item := DownloaderItem{}
	req := response.Request
	item.Url = req.URL.String()
	item.Method = req.Method
	item.Depth = req.Depth
	item.ReqBody = gcommon.ReaderToString(req.Body)
	item.RespBody = string(response.Body)
	item.Ctx = CtxToMap(response.Ctx)
	item.Status = response.StatusCode
	item.Headers = *response.Headers
	return item
}

// 存储方法
func (d *BaseDownloader) Save(item DownloaderItem) {
	data, err := item.ToJson()
	if err != nil {
		d.Logger.Printf("serialize item failed: %s", err.Error())
		return
	}
	d.Logger.Print(string(data))
}

// response钩子, 用于解析并存储每个请求的内容
func (d *BaseDownloader) OnResponse(response *gspider.Response) {
	item := d.Parse(response)
	d.Save(item)
}

// 记录开始下载时的时间, 单位: 纳秒
func (d *BaseDownloader) AddDownloadTime(r *gspider.Request) {
	r.Ctx.Put("downloadTime", gcommon.TimeStamp(1))
}

// 记录接收到返回时的时间, 单位: 纳秒
func (d *BaseDownloader) AddDownloadedTime(r *gspider.Response) {
	r.Ctx.Put("downloadedTime", gcommon.TimeStamp(1))
}
