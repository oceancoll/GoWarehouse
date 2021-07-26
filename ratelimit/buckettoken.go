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

func min(x,y int64) int64 {
	if x<y{
		return x
	}
	return y
}
func buckettokenLimit(key string, duration time.Duration, limit, capicaty int64) bool{
	formatKey := fmt.Sprintf("%s_%d_%d_%d", key, duration, limit, capicaty)
	lastTimeKey := "lastTime"
	numKey := "num"
	now := time.Now().Unix()
	client.HSetNX(formatKey, lastTimeKey, now)
	client.HSetNX(formatKey, numKey, capicaty)

	currData, err := client.HMGet(formatKey, lastTimeKey, numKey).Result()
	if err != nil{
		panic(err)
	}
	lastTime,_ := strconv.ParseInt(currData[0].(string), 0, 64)
	number,_ := strconv.ParseInt(currData[1].(string), 0, 64)
	rate := limit/int64(duration.Seconds())
	effectNum := min((now-lastTime)*rate+number, capicaty)
	if effectNum>0{
		newData := map[string]interface{}{}
		newData[lastTimeKey] = now
		newData[numKey] = effectNum-1
		client.HMSet(formatKey, newData)
		return true
	}
	return false
}

func main()  {
	for i:=0;i<20;i++{
		now := time.Now().Format("2006-1-2 15:04:05")
		res := buckettokenLimit("test3", time.Second, 5, 10)
		fmt.Printf("i is %d, now is %s, res is %t\n", i, now, res)
	}
}

