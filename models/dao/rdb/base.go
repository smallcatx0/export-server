package rdb

import (
	"math/rand"
	"strings"
	"time"

	"export-server/models/dao"
)

// K key拼接
func K(keys ...string) string {
	return strings.Join(keys, ":")
}

// BlurTTL 随机过期时间
func BlurTTL(sec int) int {
	if sec < 600 {
		return sec
	}
	rand.Seed(time.Now().UnixNano())
	return sec - 30 + rand.Intn(60)
}

// RateLimit 流量控制
func RateLimit(key string, sec, max int) bool {
	rdb := dao.Rdb
	res := rdb.Get(rdb.Context(), key)
	if res.Err() != nil {
		// 如果没有此key 创建 过期时间为time =》 true
		rdb.Set(rdb.Context(), key, 1, time.Second*time.Duration(sec))
		return true
	}
	curr, _ := res.Int()
	if max > curr {
		rdb.Incr(rdb.Context(), key)
		return true
	}
	return false
}
