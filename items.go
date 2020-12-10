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

// 下载器解析下载内容的结构
type DownloaderItem struct {
	Url      string
	ReqBody  string
	RespBody string
	Ctx      map[string]interface{}
	Depth    int
	Status   int
	Method   string
	Headers  http.Header
}

func (i DownloaderItem) ToMap() (item.ItemMap, error) {
	return item.ItemToMap(i)
}

func (i DownloaderItem) ToJson() ([]byte, error) {
	return item.ItemToJson(i)
}
