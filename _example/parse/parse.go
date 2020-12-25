package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sgs921107/gdownloader"
	"github.com/sgs921107/gdownloader/parse"
)

func main() {
	client := redis.NewClient(&redis.Options{
		// 你自己的redis配置
		Addr:     gdownloader.SettingsDemo.RedisAddr,
		Password: gdownloader.SettingsDemo.RedisPassword,
		DB:       gdownloader.SettingsDemo.RedisDB,
	})
	page, err := client.LPop("example:items").Result()
	if err != nil {
		fmt.Println("ValueError: ", err.Error())
		return
	}
	parser := parse.NewParser()
	resp, err := parser.Unmarshal(page)
	if err != nil {
		fmt.Println("Unmarshal Failed: ", err.Error())
		return
	}
	parser.Parse(resp)
}