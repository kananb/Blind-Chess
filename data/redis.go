package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx context.Context
var client *redis.Client

func init() {
	ctx = context.Background()
	host, port, pass := "localhost", "6379", "" // os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS")
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", host, port),
		Password: pass,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatal("failed to connect to redis server ", err)
	} else {
		log.Println("connected to redis server")
	}
}

type RedisStateManager struct {
	remoteSubs map[string]*redis.PubSub
	localSubs  map[string]map[chan string]bool
	lock       sync.RWMutex
}

func NewRedisStateManager() *RedisStateManager {
	return &RedisStateManager{
		remoteSubs: map[string]*redis.PubSub{},
		localSubs:  map[string]map[chan string]bool{},
	}
}

func (r *RedisStateManager) genKey() string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var num int
	var key string
	for {
		num = rand.Intn(1 << (4 * 5)) // 2^(4bits*#digits)
		key = fmt.Sprintf("%05X", num)

		if _, ok := r.remoteSubs[key]; !ok {
			return key
		}
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
	val := state.String()
	_, err := client.Set(ctx, key, string(val), 1200*time.Second).Result()
	return err == nil
}
func (r *RedisStateManager) Del(key string) bool {
	_, err := client.Del(ctx, key).Result()
	return err == nil
}

func (r *RedisStateManager) Sub(channel string) chan string {
	r.lock.Lock()
	defer r.lock.Unlock()

	ch := make(chan string)
	if r.remoteSubs[channel] == nil {
		r.remoteSubs[channel] = client.Subscribe(ctx, channel)
		r.localSubs[channel] = map[chan string]bool{}
		r.Set(channel, NewGameState())

		go func() {
			for msg := range r.remoteSubs[channel].Channel() {
				r.pubLocal(msg.Payload, channel, nil)
			}
		}()
	}
	r.localSubs[channel][ch] = true

	return ch
}
func (r *RedisStateManager) Unsub(channel string, in chan string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	select {
	case _, ok := <-in:
		if ok {
			close(in)
		}
	default:
		if in != nil {
			close(in)
		}
	}
	delete(r.localSubs[channel], in)
	if len(r.localSubs[channel]) == 0 {
		if r.remoteSubs[channel] != nil {
			r.remoteSubs[channel].Close()
			r.Del(channel)
		}
		delete(r.localSubs, channel)
		delete(r.remoteSubs, channel)
	}
}
func (r *RedisStateManager) pubLocal(msg, channel string, in chan string) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for ch := range r.localSubs[channel] {
		if ch != in {
			ch <- msg
		}
	}
}
func (r *RedisStateManager) Pub(msg, channel string, in chan string) bool {
	err := client.Publish(ctx, channel, msg).Err()
	r.pubLocal(msg, channel, in)

	return err == nil
}
