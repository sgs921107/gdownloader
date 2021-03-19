/*************************************************************************
	> File Name: items.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月08日 星期二 16时45分51秒
 ************************************************************************/

package gdownloader

import (
	"net/http"
	"encoding/json"

	"github.com/sgs921107/gcommon"
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

// ToMapSA item to map[string]interface{}
func (i DownloaderItem) ToMapSA() (gcommon.MapSA, error) {
	return gcommon.StructToMapSA(i)
}

// ToJSON item to json
func (i *DownloaderItem) ToJSON() ([]byte, error) {
	return json.Marshal(i)
}
