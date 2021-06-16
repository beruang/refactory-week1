package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"refactory/notes/internal/config"
	"time"
)

type Client interface {
	Conn() *redis.Client
	Cache() *cache.Cache
	Close() error
}

func NewClient() (Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Cfg().RdrHost, config.Cfg().RdrPort),
		DB:       config.Cfg().RdrDb,
		PoolSize: config.Cfg().RdrPool,
	})

	if err := db.Ping(context.Background()).Err(); nil != err {
		return nil, err
	}

	dbcache := cache.New(&cache.Options{
		Redis:      db,
		LocalCache: cache.NewTinyLFU(1000, time.Minute*30),
	})

	return &client{db, dbcache}, nil
}

type client struct {
	db      *redis.Client
	dbcache *cache.Cache
}

func (c *client) Conn() *redis.Client {
	return c.db
}

func (c *client) Cache() *cache.Cache {
	return c.dbcache
}

func (c *client) Close() error {
	return c.db.Close()
}
