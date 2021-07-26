package main

import (
	"fmt"
	"github.com/go-redis/redis"
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

func counterWindowLimit(key string, duration time.Duration, limit int64) bool{
	now := time.Now().Unix()
	tick := now-(now%int64(duration.Seconds()))
	formatKey := fmt.Sprintf("%s_%d_%d_%d", key, duration, limit, tick)
	_, err := client.SetNX(formatKey, 0, duration).Result()
	if err != nil{
		panic(err)
	}
	currNum, err := client.Incr(formatKey).Result()
	if err != nil{
		panic(err)
	}
	if currNum>limit{
		return false
	}
	return true
}

func test1()  {
	for i:=0;i<10;i++{
		res := counterWindowLimit("test1", 1*time.Minute, 5)
		fmt.Printf("conn num %d, result is %t \n", i, res)
	}
}
func main()  {
	test1()
}



