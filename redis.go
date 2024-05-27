package dhredis

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	dhlog "github.com/lepingbeta/go-common-v2-dh-log"

	"github.com/gomodule/redigo/redis"
)

type redisConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Pass        string `json:"pass"`
	Network     string `json:"network"`
	MaxIdle     int    `json:"max-idle"`
	IdleTimeout int    `json:"idle-timeout"`
	MaxActive   int    `json:"max-active"`
	Db          int    `json:"db"`
	Prefix      string `json:"prefix"`
}

var (
	pool    *redis.Pool
	rConfig redisConfig
	//redisServer = flag.String("127.0.0.1", ":6379", "")
)

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:         rConfig.MaxIdle,
		MaxActive:       rConfig.MaxActive,
		IdleTimeout:     time.Duration(rConfig.IdleTimeout) * time.Second,
		MaxConnLifetime: time.Duration(rConfig.IdleTimeout*2) * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(rConfig.Network, addr)
			if err != nil {
				return nil, err
			}
			//此处1234对应redis密码
			if _, err := conn.Do("AUTH", rConfig.Pass); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, err
		},
	}
}

func IsKeyExists(key string) (bool, error) {
	key = GetRealKey(key)
	// dhlog.Info("real key: %s", key)
	is_key_exit, err := redis.Bool(GetConn().Do("EXISTS", key))
	if err != nil {
		dhlog.Error("error: %v", err.Error())
		return false, err
	}
	return is_key_exit, nil
}

func Del(k string) {
	k = GetRealKey(k)
	DelByRealKey(k)
}

func DelByRealKey(k string) {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("DEL", k)
	if err != nil {
		dhlog.Error("Redis SET error: %s", err.Error())
	}
}

func Set(k, v string, ttl int) {
	k = GetRealKey(k)
	conn := GetConn()
	defer conn.Close()
	// replySet, err := conn.Do("SET", k, v)
	_, err := conn.Do("SET", k, v)
	if err != nil {
		dhlog.Error("Redis SET error: %s", err.Error())
	}
	// conn.Do("EXPIRE", k, ttl)
	if ttl > 0 {
		expire(k, ttl)
	}

	// dhlog.Infoln(replySet)
}

/**
 * @description: 扫描redis匹配的所有keys（本函数使用scan命令，不能用keys命令。因为keys命令会阻塞线上其它redis请求）
 * @param {string} pattern
 * @return {*}
 */
func ScanKeys(pattern string) ([]string, error) {
	match := GetRealKey(pattern)
	conn := GetConn()
	defer conn.Close()
	var (
		cursor = 0 // 默认从0开始
		args   = make([]interface{}, 0, 3)
	)

	args = append(args, cursor)
	args = append(args, "MATCH", match)
	// count := 4
	// args = append(args, "COUNT", count)

	var keyList []string
	for {
		args[0] = cursor
		res, err := redis.Values(conn.Do("SCAN", args...))
		if err != nil {
			return keyList, err
		}

		var tmp []string
		_, err = redis.Scan(res, &cursor, &tmp)

		if err != nil {
			return keyList, err
		}

		if lt := len(tmp); lt > 0 {

			keyList = append(keyList, tmp...)
		}

		if cursor == 0 {
			break // 查询结束
		}
	}
	return keyList, nil
}

/**
 * @description: 设置key过期时间，给外部调用
 * @param {string} k
 * @param {int} ttl
 * @return {*}
 */
func Expire(k string, ttl int) {
	k = GetRealKey(k)
	expire(k, ttl)
}

/**
 * @description: 设置key过期时间，给内部调用
 * @param {string} k
 * @param {int} ttl
 * @return {*}
 */
func expire(k string, ttl int) {
	conn := GetConn()
	defer conn.Close()
	conn.Do("EXPIRE", k, ttl)
}

/**
 * @description: 从redis获取值
 * @param {string} k
 * @return {*}
 */
func Get(k string) string {
	k = GetRealKey(k)
	return GetByRealKey(k)
}

func GetByRealKey(k string) string {
	conn := GetConn()
	defer conn.Close()
	reply, err := redis.String(conn.Do("GET", k))
	if err != nil {
		dhlog.Error("Redis GET error: %s", err.Error())
		dhlog.Error("%v", reply)
		return ""
	}
	// dhlog.Infoln(reply)
	return reply
	// return fmt.Sprintf("%v", reply)
}

func GetRealKey(k string) string {
	return fmt.Sprintf("%s:%s", rConfig.Prefix, k)
}

func TTL(k string) int {
	k = GetRealKey(k)
	conn := GetConn()
	defer conn.Close()
	reply, err := conn.Do("TTL", k)
	dhlog.Info("续租输出reply(当前剩余有效期)是：%v", reply)
	if err != nil {
		dhlog.Error("get token ttl error %s", err.Error())
	}
	ttl, _ := strconv.Atoi(fmt.Sprintf("%v", reply))
	return ttl
}

func GetConn() redis.Conn {
	conn := pool.Get()
	conn.Do("SELECT", rConfig.Db)
	return conn
}

func ActiveCount() int {
	return pool.ActiveCount()
}

func IdleCount() int {
	return pool.IdleCount()
}

/**
 * @description: 主要作为互斥锁用
 * @param {string} key
 * @param {interface{}} val
 * @param {int} timeout
 * @return {*}
 */
func SetExNx(key string, val interface{}, timeout int) (reply interface{}, err error) {
	conn := GetConn()
	defer conn.Close()
	key = GetRealKey(key)

	reply, err = conn.Do("SET", key, val, "EX", timeout, "NX")
	dhlog.Info("key: %s timeout: %v reply: %v, err: %v", key, timeout, reply, err)
	return reply, err
}

func initRedis(jsonStr string) {
	// hecos.LoadConfigPlus("common", "redis", &rConfig)

	// Unmarshal the JSON data into the structure
	err := json.Unmarshal([]byte(jsonStr), &rConfig)
	if err != nil {
		dhlog.Error(err.Error())
	}

	dhlog.DebugAny(rConfig)

	pool = newPool(fmt.Sprintf("%s:%v", rConfig.Host, rConfig.Port))
}

func GetConfig() redisConfig {
	return rConfig
}

// Push pushes a message onto the Redis List.
func Push(key string, message string) error {
	conn := GetConn()
	defer conn.Close()
	key = GetRealKey(key)

	_, err := conn.Do("LPUSH", key, message)
	if err != nil {
		dhlog.Error("failed to push message: %s", err.Error())
		return err
	}
	return nil
}

// Pop pops a message from the Redis List.
func Pop(key string) (string, error) {
	conn := GetConn()
	defer conn.Close()
	key = GetRealKey(key)

	result, err := conn.Do("RPOP", key)
	if err != nil {
		dhlog.Error("failed to pop message: %s", err.Error())
		return "", err
	}
	if result == nil {
		return "", nil
	}
	return fmt.Sprintf("%s", result), nil
}

// BRPop 从 Redis List 的尾部弹出一个元素，如果 List 中没有元素则阻塞，直到有元素可用。
func BRPop(key string, timeoutSec int) (string, error) {
	conn := GetConn()
	defer conn.Close()
	key = GetRealKey(key)

	buffReader, err := redis.Strings(conn.Do("brpop", key, timeoutSec))
	if err != nil {
		log.Println("BRPop msg error:", err)
		return "BRPop msg error:" + err.Error(), err
	}

	if len(buffReader) > 1 { //数组第一个元素取出来是队列名
		//log.Println("BRPop msg success:", buffReader[1], " BRPop queue: ", buffReader[0])
		return buffReader[1], err
	}

	// 返回结果
	return "", errors.New("empty data")
}
