/*
 * @Author       : Symphony zhangleping@cezhiqiu.com
 * @Date         : 2022-11-19 23:50:11
 * @LastEditors  : Symphony zhangleping@cezhiqiu.com
 * @LastEditTime : 2024-05-28 02:35:19
 * @FilePath     : /v2/go-common-v2-dh-redis/redis_test.go
 * @Description  :
 *
 * Copyright (c) 2022 by 大合前研, All Rights Reserved.
 */
/*
 * @Author: dev.cezhiqiu.cn dev@cezhiqiu.cn
 * @Date: 2022-06-25 21:42:36
 * @LastEditors: dev.cezhiqiu.cn dev@cezhiqiu.cn
 * @LastEditTime: 2022-06-26 14:19:04
 * @FilePath: /api/common/redis/redis_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package dhredis

import (
	"testing"

	dhlog "github.com/lepingbeta/go-common-v2-dh-log"
)

func TestGet(t *testing.T) {
	jsonStr := `{
        "host": "dev.cezhiqiu.cn",
        "port": 6379,
        "pass": "dahe_redis_888",
        "network": "tcp",
        "max-idle": 3,
        "idle-timeout": 240,
        "max-active": 0,
        "db": 6,
        "prefix": "redis-test"
    }`
	InitRedis(jsonStr)
	for i := 0; i < 10; i++ {
		c := GetConn()
		defer c.Close()
		Set("zzz", "444", -1)
		dhlog.Info("%v %v %v %v", c, "active count: ", ActiveCount(), "idle count: ", IdleCount())
		// k := "kkk"
		// Set(k, "vvv", 60)
		// r := TTL(k)
		// if r != 60 {
		// 	t.Errorf("ttl 函数错误")
		// }
		// r2 := TTL("abcdefgh")
		// if r2 > 0 {
		// 	t.Errorf("ttl 取值异常")
		// }
		// r3 := Get(k)
		// dhlog.Infoln(r3)

		// r4 := Get("abc")
		// dhlog.Infoln(r4, r4 == "")

		// res, err := c.Do("ping")
		// if err != nil {
		// 	t.Errorf(err.Error())
		// 	dhlog.Errorln("go---1---error:", err)
		// 	return
		// }
		// dhlog.Infoln("go---1---:", res.(string))
	}
}

// func TestScanKeys(t *testing.T) {

// 	l, e := ScanKeys("lock-list*")
// 	dhlog.Infoln(l, len(l), e)
// }

// func TestPush(t *testing.T) {
// 	key := "aa"
// 	msg := "bbb"
// 	Push(key, msg)
// }

// func TestPop(t *testing.T) {
// 	key := "aa"
// 	r, e := Pop(key)
// 	t.Logf("111")
// 	dhlog.Info("%v, %v", r, e)
// }

// func TestBRPop(t *testing.T) {
// 	key := "aa"
// 	r, err := BRPop(key, 5)

// 	if err != nil {
// 		if e, ok := err.(redis.Error); ok {
// 			if e.Error()[:7] == "timeout" {
// 				// 处理超时错误
// 				fmt.Println("111111 Timeout error:", e.Error())
// 			} else {
// 				// 处理其他错误
// 				fmt.Println("222222222 Redis error:", err)
// 			}
// 		} else {
// 			// 处理其他错误
// 			fmt.Println("33333333 Redis error:", err)
// 		}

// 	}

// 	dhlog.Infoln(r)
// 	dhlog.Infoln(err)
// }
