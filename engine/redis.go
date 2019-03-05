package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/nishrit/aero/key"
	"gopkg.in/redis.v3"
)

type Redis struct {
	key.AsIsFormat
	r    *redis.Client
	name string
}

func NewRedis(host string, port int, db int) Redis {
	serv := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
		DB:   int64(db),
	})
	_, err := serv.Ping().Result()
	if err != nil {
		panic(err)
	}

	return Redis{
		r: serv,
	}
}

func NewRedis2(host string, port int, db int, name string) Redis {
	serv := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
		DB:   int64(db),
	})
	_, err := serv.Ping().Result()
	if err != nil {
		panic(err)
	}

	return Redis{
		r:    serv,
		name: name,
	}
}

func (rd Redis) Set(key string, data []byte, expireIn time.Duration) {
	key = rd.Format(key)
	rd.r.Set(key, data, expireIn)
}

func (rd Redis) Get(key string) ([]byte, error) {
	key = rd.Format(key)
	data, err := rd.r.Get(key).Bytes()
	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (rd Redis) Delete(key string) error {
	key = rd.Format(key)
	return rd.r.Del(key).Err()
}

func (rd Redis) Close() {
	err := rd.r.Close()
	if err != nil {
		panic(err)
	}
}

func (rd Redis) Push(data []byte) error {
	return rd.r.LPush(rd.name, string(data)).Err()
}

func (rd Redis) Pop() ([]byte, error) {
	s, e := rd.r.RPop(rd.name).Result()
	return []byte(s), e
}

func (rd Redis) PopWait(dur time.Duration) ([]byte, error) {
	sa, e := rd.r.BRPop(dur, rd.name).Result()

	// did the time expire?
	if len(sa) == 0 && strings.Contains(e.Error(), "nil") {
		return []byte(""), nil
	}

	if len(sa) == 0 {
		return []byte(""), e
	} else {
		return []byte(sa[1]), e
	}
}

func (rd Redis) Len() (int, error) {
	i, e := rd.r.LLen(rd.name).Result()
	return int(i), e
}
