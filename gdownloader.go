/*************************************************************************
	> File Name: redis_downloader.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月08日 星期二 16时28分05秒
 ************************************************************************/

package gdownloader

import (
	"bytes"
	"compress/gzip"
	"strings"

	"github.com/antchfx/htmlquery"
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

// Downloader 定义下载器接口
type Downloader interface {
	// 解析方法
	parse(response *gspider.Response) (*DownloaderItem, error)
	// 存储方法
	save(item *DownloaderItem)
	// spider实例
	Spider() *gspider.RedisSpider
	// 启动
	Run()
}

// BaseDownloader 定义下载器的机构
type BaseDownloader struct {
	spider   *gspider.RedisSpider
	logger   *gspider.Logger
	settings *Settings
	save     func(item *DownloaderItem)
}

// Spider spider实例
func (d *BaseDownloader) Spider() *gspider.RedisSpider {
	return d.spider
}

func (d BaseDownloader) clearHead(text string) string {
	if doc, err := htmlquery.Parse(strings.NewReader(text)); err == nil {
		if bodyNode := htmlquery.FindOne(doc, "/html/body"); bodyNode != nil {
			return strings.TrimSpace(htmlquery.OutputHTML(bodyNode, false))
		}
	}
	return text
}

func (d BaseDownloader) compress(text string) (string, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	defer gz.Close()
	if _, err := gz.Write([]byte(text)); err != nil {
		return "", err
	}
	if err := gz.Flush(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// parse 解析方法
func (d BaseDownloader) parse(response *gspider.Response) (item *DownloaderItem, err error) {
	respBody := string(response.Body)
	// 清除head, 只保留body数据
	if d.settings.Downloader.ClearHead {
		respBody = d.clearHead(respBody)
	}
	if d.settings.Downloader.GzipCompress {
		respBody, err = d.compress(respBody)
		if err != nil {
			return item, err
		}
	}
	item = &DownloaderItem{}
	req := response.Request
	item.URL = req.URL.String()
	item.Method = req.Method
	item.Depth = req.Depth
	item.ReqBody = gcommon.ReaderToString(req.Body)
	item.RespBody = respBody
	item.Ctx = CtxToMap(response.Ctx)
	item.Status = response.StatusCode
	item.Headers = *response.Headers
	return item, nil
}

// example
func (d BaseDownloader) _save(item *DownloaderItem) {
	data, err := item.ToJSON()
	if err != nil {
		d.logger.Errorw("Serialize Item Failed",
			"errMsg", err.Error(),
		)
		return
	}
	d.logger.Debug(string(data))
}

// onResponse response钩子, 用于解析并存储每个请求的内容
func (d BaseDownloader) onResponse(response *gspider.Response) {
	item, err := d.parse(response)
	if err != nil {
		d.logger.Errorw("Parse Response Error",
			"errMsg", err.Error(),
		)
		return
	}
	item.Ctx["saveTime"] = gcommon.TimeStamp(1)
	d.save(item)
}

// addDownloadTime 记录开始下载时的时间, 单位: 毫秒
func (d BaseDownloader) addDownloadTime(r *gspider.Request) {
	r.Ctx.Put("downloadTime", gcommon.TimeStamp(1))
}

// addDownloadedTime 记录接收到返回时的时间, 单位: 毫秒
func (d BaseDownloader) addDownloadedTime(r *gspider.Response) {
	downloadedTime := gcommon.TimeStamp(1)
	r.Ctx.Put("downloadedTime", downloadedTime)
	req := r.Request
	d.logger.Infow(
		"downloaded",
		"status", r.StatusCode,
		"url", req.URL.String(),
		"method", req.Method,
		"headers", *req.Headers,
		"ctx", CtxToMap(r.Ctx),
		"ms", downloadedTime-r.Ctx.GetAny("downloadTime").(int64),
	)
}

// errorHandler 记录错误请求
func (d BaseDownloader) errorHandler(r *gspider.Response, err error) {
	req := r.Request
	d.logger.Errorw(
		"DownloadError",
		"errMsg", err.Error(),
		"status", r.StatusCode,
		"url", req.URL.String(),
		"method", req.Method,
		"headers", *req.Headers,
		"ctx", CtxToMap(r.Ctx),
		"ms", gcommon.TimeStamp(1)-r.Ctx.GetAny("downloadTime").(int64),
	)
}

func (d *BaseDownloader) init() {
	if d.spider == nil {
		panic("Spider Not Instance")
	}
	if d.save == nil {
		d.save = d._save
	}
	d.spider.OnRequest(d.addDownloadTime)
	d.spider.OnResponse(d.addDownloadedTime)
	d.spider.OnResponse(d.onResponse)
	d.spider.OnError(d.errorHandler)
}

// Run run downloader
func (d *BaseDownloader) Run() {
	d.init()
	// 添加信号扩展
	d.spider.AddExtension(NewSignalExtension())
	d.spider.Start()
}
