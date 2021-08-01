package main

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sgs921107/gdownloader"
	"github.com/sgs921107/gdownloader/parse"
	"github.com/sgs921107/gredis"
)

func main() {
	settings := gdownloader.NewDownloaderSettings("env_demo")
	client := redis.NewClient(&redis.Options{
		// 你自己的redis配置
		Addr:     settings.Redis.Addr,
		Password: settings.Redis.Password,
		DB:       settings.Redis.DB,
	})
	parser := parse.NewParser()
	prefix := settings.Redis.Prefix
	for {
		page, err := client.LPop(prefix + ":items:www.example.com").Result()
		if err == gredis.RedisNil{
			break
		} else if err != nil {
			fmt.Println("ValueError: ", err.Error())
			continue
		}
		resp, err := parser.Unmarshal(page)
		if err != nil {
			fmt.Println("Unmarshal Failed: ", err.Error())
			continue
		}
		parser.Parse(resp)
	}
}