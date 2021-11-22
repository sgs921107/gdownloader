/*************************************************************************
	> File Name: settings.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月09日 星期三 09时49分57秒
 ************************************************************************/

package gdownloader

import (
	"sync"

	"github.com/sgs921107/gcommon"
	"github.com/sgs921107/gspider"
)

var (
	// DefaultEnvPath defalut env path
	defaultEnvPath = "/etc/gdownloader/.env"
	settings = DownloaderSettings{}
	settingsOnce sync.Once
)

// SpiderSettings spider settings type
type SpiderSettings = gspider.SpiderSettings

// DownloaderSettings downloader的配置结构
// 使数据结构简单，不继承自spider settings, 通过反射来生成spdier settings
type DownloaderSettings struct {
	// SpiderSettings
	SpiderSettings

	// download settings
	Downloader	struct {
		// 存储页面数据的最大数量  list元素超出将被裁剪, 避免内存过高
		MaxTopicSize	int64	`default:"10000"`
  		// 是否清除html页面的head内容, 只保留body数据
		ClearHead		bool 	`default:"false"`
		// 是否使用gzip对下载页面内容进行压缩
		GzipCompress	bool	`default:"false"`
	}
}


// NewDownloaderSettings new a downlaoder settings
func NewDownloaderSettings(envFiles ...string) DownloaderSettings {
	settingsOnce.Do(func(){
		if len(envFiles) == 0 {
			envFiles = append(envFiles, defaultEnvPath)
		}
		gcommon.OverLoadEnvFiles(envFiles...)
		gcommon.EnvIgnorePrefix()
		gcommon.EnvFill(&settings)
		gcommon.EnvFill(&settings.SpiderSettings)
	})
	return settings
}
