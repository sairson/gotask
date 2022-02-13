package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var RedisDB *redis.Client

// InitRedisPool 连接数据库
func InitRedisPool(ip string, port int, pass string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", ip, port),
		Password: pass,
		DB:       0,
	})
	// 判断时否连接成功
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	RedisDB = rdb // 赋值给RedisDB

	return nil
}

func RedisSet(key string, job Job) error {
	add, err := json.Marshal(&job) // 将job转json
	if err != nil {
		return err
	}
	RedisDB.Set(context.Background(), key, add, 0) // 添加一个键值
	return nil
}

func RedisDel(key string) error {
	_, err := RedisDB.Del(context.Background(),key).Result()
	if err != nil {
		return err
	}
	return nil
}

func RedisGet(key string) Job {
	result, err := RedisDB.Get(context.Background(), key).Result()
	if err != nil {
		// 返回一个空
		return Job{}
	}
	var job Job
	err = json.Unmarshal([]byte(result), &job)
	if err != nil {
		return Job{}
	}
	return job
}