/*
 * @Author: xiangcai
 * @Date: 2021-08-30 21:13:51
 * @LastEditors: xiangcai
 * @LastEditTime: 2021-08-30 21:29:10
 * @Description: file content
 */
package gdownloader

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sgs921107/gspider"
)

// SignalExtension  signal extension
type SignalExtension struct{}

func (e SignalExtension) Run(spider *gspider.BaseSpider) {
	//创建监听退出chan
	c := make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for sig := range c {
		logger := spider.Logger.With(
			"signal", sig,
		)
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			spider.Quit()
			logger.Warnw("Recv Exit Signal")
		case syscall.SIGUSR1:
			logger.Infow("Recv USR1 Signal")
		case syscall.SIGUSR2:
			logger.Infow("Recv USR2 Signal")
		default:
			logger.Infow("Recv Other Signal")
		}
	}
}

// NewSignalExtension new a signal extension
func NewSignalExtension() gspider.Extension {
	return &SignalExtension{}
}
