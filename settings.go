/*************************************************************************
	> File Name: settings.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月09日 星期三 09时49分57秒
 ************************************************************************/

package gdownloader

import (
	"github.com/sgs921107/gcommon"
	"github.com/sgs921107/gspider"
)

// SpiderSettings spider settings type
type SpiderSettings = gspider.SpiderSettings

// Settings downloader的配置结构
// 使数据结构简单，不继承自spider settings, 通过反射来生成spdier settings
type Settings struct {
	// SpiderSettings
	SpiderSettings

	// download settings
	Downloader struct {
		// 存储页面数据的最大数量  list元素超出将被裁剪, 避免内存过高
		// 0则不进行裁剪
		MaxTopicSize int64 `default:"10000"`
	}
}

// NewSettingsFromEnv new a downlaoder settings from env
func NewSettingsFromEnv() (*Settings, error) {
	var settings Settings
	gcommon.EnvIgnorePrefix()
	if err := gcommon.EnvFill(&settings); err != nil {
		return nil, err
	}
	if err := gcommon.EnvFill(&settings.SpiderSettings); err != nil {
		return nil, err
	}
	return &settings, nil
}

// NewSettingsFromEnvFile new a downloader settings from env file
func NewSettingsFromEnvFile(envFile string) (*Settings, error) {
	if err := gcommon.LoadEnvFile(envFile, true); err != nil {
		return nil, err
	}
	return NewSettingsFromEnv()
}
