package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)
var client *redis.Client

func init()  {
	client = redis.NewClient(
		&redis.Options{
			Addr:       "localhost:6379",
			Password:   "",
			DB:         0,
			PoolSize:   3,
			MaxRetries: 3,
		},
	)
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func slidingWindowLimit(key string, duration time.Duration, segmentNum, limit int64) bool {
	now := time.Now().Unix()
	segmentDuration := int64(duration.Seconds())/segmentNum
	tick := now-(now%segmentDuration)
	formatKey := fmt.Sprintf("%s_%d_%d_%d_%d", key, duration, limit, segmentNum, tick)
	client.SetNX(formatKey, 0, duration)
	currNum, err := client.Incr(formatKey).Result()
	if err != nil{
		panic(err)
	}
	for i:=int64(1);i<=segmentNum;i++{
		preKey := fmt.Sprintf("%s_%d_%d_%d_%d", key, duration, limit, segmentNum, tick-segmentDuration*i)
		preNum, err := client.Get(preKey).Result()
		if err != nil{
			if err == redis.Nil{
				preNum = "0"
			} else {
				panic(err)
			}
		}
		num, err := strconv.ParseInt(preNum, 0, 64)
		if err != nil{
			panic(err)
		}
		currNum += num
		if currNum> limit{
			client.Decr(formatKey)
			return false
		}
	}
	return true
}

func main()  {
	for i:=0;i<30;i++{
		now := time.Now().Format("2006-1-2 15:04:05")
		res := slidingWindowLimit("test2", time.Minute, 6, 3)
		fmt.Printf("now is %s, res is %t\n", now, res)
		time.Sleep(time.Second*5)
	}
}
