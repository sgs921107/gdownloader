/*************************************************************************
	> File Name: items.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月08日 星期二 16时45分51秒
 ************************************************************************/

package gdownloader

import (
	"github.com/sgs921107/gspider/item"
	"net/http"
)

// DownloaderItem 下载器解析下载内容的结构
type DownloaderItem struct {
	URL      string
	ReqBody  string
	RespBody string
	Ctx      map[string]interface{}
	Depth    int
	Status   int
	Method   string
	Headers  http.Header
}

// ToMap item to map
func (i DownloaderItem) ToMap() (item.Map, error) {
	return item.ToMap(i)
}

// ToJSON item to json
func (i DownloaderItem) ToJSON() ([]byte, error) {
	return item.ToJSON(i)
}
