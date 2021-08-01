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
)

// CtxMap 上下文的类型
type CtxMap map[string]interface{}

// CtxToMap 将response中的上下文转为自定义的上下文结构
func CtxToMap(ctx *gspider.Context) CtxMap {
	data := make(CtxMap)
	ctx.ForEach(func(k string, v interface{}) interface{} {
		data[k] = v
		return nil
	})
	return data
}

// BaseDownloader 定义下载器的机构
type BaseDownloader struct {
	Spider   *gspider.RedisSpider
	Logger	 *gspider.Logger
	settings DownloaderSettings
	onResponse	 func(resp *gspider.Response)
}

// Parse 解析方法
func (d *BaseDownloader) Parse(response *gspider.Response) *DownloaderItem {
	item := &DownloaderItem{}
	req := response.Request
	item.URL = req.URL.String()
	item.Method = req.Method
	item.Depth = req.Depth
	item.ReqBody = gcommon.ReaderToString(req.Body)
	item.RespBody = string(response.Body)
	item.Ctx = CtxToMap(response.Ctx)
	item.Status = response.StatusCode
	item.Headers = *response.Headers
	return item
}

// Save 存储方法
func (d *BaseDownloader) save(item *DownloaderItem) {
	data, err := item.ToJSON()
	if err != nil {
		d.Logger.WithFields(gspider.LogFields{
			"errMsg": err.Error(),
		}).Error("Serialize item failed")
		return
	}
	d.Logger.Debug(string(data))
}

// OnResponse response钩子, 用于解析并存储每个请求的内容
func (d *BaseDownloader) OnResponse(response *gspider.Response) {
	item := d.Parse(response)
	d.save(item)
}

// AddDownloadTime 记录开始下载时的时间, 单位: 纳秒
func (d *BaseDownloader) AddDownloadTime(r *gspider.Request) {
	r.Ctx.Put("downloadTime", gcommon.TimeStamp(1))
}

// AddDownloadedTime 记录接收到返回时的时间, 单位: 纳秒
func (d *BaseDownloader) AddDownloadedTime(r *gspider.Response) {
	r.Ctx.Put("downloadedTime", gcommon.TimeStamp(1))
}

func (d *BaseDownloader) init() {
	if d.Spider == nil {
		panic("Spider Not Instance")
	}
	if d.onResponse == nil {
		d.onResponse = d.OnResponse
	}
	d.Spider.OnRequest(d.AddDownloadTime)
	d.Spider.OnResponse(d.AddDownloadedTime)
	d.Spider.OnResponse(d.onResponse)
}

// Run run downloader
func (d *BaseDownloader) Run() {
	d.init()
	d.Spider.Start()
}
