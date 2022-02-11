package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var ctx context.Context
var client *redis.Client

func init() {
	host, pass := os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASS")
	client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatal("failed to connect to redis server", err)
	}
}

type RedisStateManager struct {
	subs map[string]*redis.PubSub
	lock sync.RWMutex
}

func NewRedisStateManager() *RedisStateManager {
	return &RedisStateManager{subs: map[string]*redis.PubSub{}}
}

func (r *RedisStateManager) GenKey() string {
	var num int
	var key string
	for {
		num = rand.Intn(1048576)
		key = fmt.Sprintf("%05X", num)

		r.lock.RLock()
		if _, ok := r.subs[key]; !ok {
			return key
		}
		r.lock.RUnlock()
	}
}
func (r *RedisStateManager) Get(key string) *GameState {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	state := new(GameState)
	err = json.Unmarshal([]byte(val), state)
	if err != nil {
		return nil
	}

	return state
}
func (r *RedisStateManager) Set(key string, state *GameState) bool {
	val, err := json.Marshal(state)
	if err != nil {
		return false
	}

	_, err = client.Set(ctx, key, string(val), 1200).Result()
	return err == nil
}
func (r *RedisStateManager) Del(key string) bool {
	return true
}

func (r *RedisStateManager) Sub(channel string) <-chan string {
	r.lock.Lock()
	defer r.lock.Unlock()

	sub, ok := r.subs[channel]
	if !ok {
		sub = client.Subscribe(ctx, channel)
		r.subs[channel] = sub
	}

	ch := make(chan string)
	go func() {
		for msg := range sub.Channel() {
			ch <- msg.Payload
		}
		close(ch)
	}()
	return ch
}
func (r *RedisStateManager) Unsub(channel string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	sub, ok := r.subs[channel]
	if !ok {
		return false
	}

	sub.Close()

	return true
}
func (r *RedisStateManager) Pub(msg, channel string) bool {
	err := client.Publish(ctx, channel, msg).Err()
	return err == nil
}
