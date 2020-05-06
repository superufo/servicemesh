package redis

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-mesh/openlogging"
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v2"
)

type RedisCfg struct {
	Addr        string `yaml:"addr"`        // Redis地址
	Password    string `yaml:"pwd"`         // Redis账号
	DB          int    `yaml:"dbnum"`       // Redis库
	PoolSize    int    `yaml:"poolsize"`    // Redis连接池大小
	MaxRetries  int    `yaml:"maxRetries"`  // 最大重试次数
	IdleTimeout int    `yaml:"idleTimeout"` // 空闲链接超时时间
}

var Client *redis.Client

func init() {
	file := "libs/common/redis/redis.yaml"
	yamlContent := make([]byte, 1000)
	var error error
	for i := 1; i < 10; i++ {
		file = "../" + file
		_, fileErr := os.Stat(file)
		if fileErr == nil {
			yamlContent, error = ioutil.ReadFile(file)
			break
		}
	}

	//yamlContent, err := ioutil.ReadFile("../../../../libs/common/redis/redis.yaml")
	if error != nil {
		panic("读取redis.yaml失败")
		return
	}

	redisOption := RedisCfg{}
	err := yaml.Unmarshal([]byte(yamlContent), &redisOption)
	if err != nil {
		openlogging.Error("error", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
		panic("解析redis.yaml失败")
		return
	}

	Client = redis.NewClient(&redis.Options{
		Addr:       redisOption.Addr,
		Password:   redisOption.Password,
		DB:         redisOption.DB,
		PoolSize:   redisOption.PoolSize,
		MaxRetries: redisOption.MaxRetries,
		// redisOption.IdleTimeout,
	})
	// Client = redis.NewClient(&redis.Options{
	// 	Addr:       "127.0.0.1",
	// 	Password:   "",
	// 	DB:         0,
	// 	PoolSize:   2048,
	// 	MaxRetries: 3,
	// })
	fmt.Println(Client)

	ping, err := Client.Ping().Result()
	if err == redis.Nil {
		openlogging.Error("Redis异常", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
	} else if err != nil {
		openlogging.Error("Redis失败", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
	} else {
		openlogging.Info(ping)
	}
}

// HashSet 向key的hash中添加元素field的值
func HashSet(key, field string, data interface{}) {
	err := Client.HSet(key, field, data)
	if err != nil {
		openlogging.Error("Redis HSet Error", openlogging.WithTags(openlogging.Tags{
			"err": err,
		}))
	}
}

// BatchHashSet 批量向key的hash添加对应元素field的值
func BatchHashSet(key string, fields map[string]interface{}) string {
	val, err := Client.HMSet(key, fields).Result()
	if err != nil {
		openlogging.Error("Redis HMSet Error", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
	}

	return string(val)
}

// HashGet 通过key获取hash的元素值
func HashGet(key, field string) string {
	result := ""
	val, err := Client.HGet(key, field).Result()
	if err == redis.Nil {
		openlogging.Error("Key Doesn't Exists", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
		return result
	} else if err != nil {
		openlogging.Error("Redis HGet Error", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
		return result
	}
	return val
}

// BatchHashGet 批量获取key的hash中对应多元素值
func BatchHashGet(key string, fields ...string) map[string]interface{} {
	resMap := make(map[string]interface{})
	for _, field := range fields {
		var result interface{}
		val, err := Client.HGet(key, fmt.Sprintf("%s", field)).Result()
		if err == redis.Nil {
			openlogging.Error("Key Doesn't Exists", openlogging.WithTags(openlogging.Tags{
				"err": err.Error(),
			}))
			resMap[field] = result
		} else if err != nil {
			openlogging.Error("Redis HMGet Error:", openlogging.WithTags(openlogging.Tags{
				"err": err.Error(),
			}))
			resMap[field] = result
		}
		if val != "" {
			resMap[field] = val
		} else {
			resMap[field] = result
		}
	}
	return resMap
}

// Incr 获取自增唯一ID
func Incr(key string) int {
	val, err := Client.Incr(key).Result()
	if err != nil {
		openlogging.Error("Redis Incr Error:", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
	}
	return int(val)
}

// SetAdd 添加集合数据
func SetAdd(key, val string) {
	Client.SAdd(key, val)
}

// SetGet  从集合中获取数据
func SetGet(key string) []string {
	val, err := Client.SMembers(key).Result()
	if err != nil {
		openlogging.Error("Redis SMembers Error:", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
	}
	return val
}
