package db

import (
	"context"
	"encoding/json"
	"github.com/hqbobo/frame/common/log"

	"github.com/go-redis/redis/v8"
)

type redisStore struct {
	client     *redis.Client
	clusterCLi *redis.ClusterClient
	ip         string
	pass       string
	name       string
	cluster    bool
}

func newRedis(ip string, pass string) *redisStore {
	s := new(redisStore)
	s.pass = pass
	s.client = redis.NewClient(&redis.Options{
		Addr:     ip,
		Password: pass, // no password set
		DB:       0,
	})
	s.cluster = false
	return s
}

func newRedisCluster(ip []string, pass string) *redisStore {
	s := new(redisStore)
	s.pass = pass
	s.clusterCLi = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    ip,
		Password: pass, // no password set
	})
	s.cluster = true
	return s
}

var ctx = context.Background()

// CacheStore is a interface to store cache
type redisCacheStore struct {
	store *redisStore
}

func (rcs *redisCacheStore) Put(key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if rcs.store.cluster {
		return rcs.store.clusterCLi.Set(ctx, key, string(val), 0).Err()
	}
	return rcs.store.client.Set(ctx, key, string(val), 0).Err()
}
func (rcs *redisCacheStore) Get(key string) (interface{}, error) {
	var obj interface{}
	var result []byte
	var err error
	if rcs.store.cluster {
		result, err = rcs.store.clusterCLi.Get(ctx, key).Bytes()
		if err != nil {
			return nil, err
		}
	}
	result, err = rcs.store.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (rcs *redisCacheStore) Del(key string) error {
	if rcs.store.cluster {
		return rcs.store.clusterCLi.Del(ctx, key).Err()
	}
	return rcs.store.client.Del(ctx, key).Err()
}

func newCacheStore(addrs []string, pass string) *redisCacheStore {
	s := new(redisCacheStore)
	if len(addrs) > 1 {
		log.Infof("初始化redisStore %s", (addrs[0]))
		s.store = newRedisCluster(addrs, pass)
	} else if len(addrs) == 1 {
		log.Infof("初始化redisStore %s", (addrs[0]))
		s.store = newRedis(addrs[0], pass)
	}
	return s
}
