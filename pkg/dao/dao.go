package dao

import (
	"context"
	"sync"
	"time"
)

type entry[dataModel any] struct {
	mu      sync.RWMutex
	loading bool
	ready   chan struct{}
	value   *dataModel
	epoch   int64
}

type Cache[dataModel any] struct {
	items map[string]*entry[dataModel]
	mapMu sync.RWMutex
	epoch int64
}

func NewCache[dataModel any]() *Cache[dataModel] {
	return &Cache[dataModel]{
		items: make(map[string]*entry[dataModel]),
	}
}

type DBContextInterface[Context any] interface {
	GetClient() *Context
}
type Info[dataModel any] interface {
	GetKey() string
	GetEmptyValue() *dataModel
	SetInfo(*dataModel)
	WhetherExist() bool
}

func Get[dataModel any, dbContextModel any](ctx context.Context, c *Cache[dataModel], DBContext DBContextInterface[dbContextModel], info Info[dataModel], getCache func(context.Context, DBContextInterface[dbContextModel], Info[dataModel]) (*dataModel, bool, error), setCache func(context.Context, DBContextInterface[dbContextModel], Info[dataModel]) error, getDB func(context.Context, DBContextInterface[dbContextModel], Info[dataModel]) (*dataModel, bool, error)) (bool, error) {
	value, exist, err := getCache(ctx, DBContext, info)
	if !exist || err != nil {
		c.mapMu.RLock()
		k, ok := c.items[info.GetKey()]
		c.mapMu.RUnlock()
		if ok {
			k.mu.RLock()
			if !k.loading {
				defer k.mu.RUnlock()
				info.SetInfo(k.value)
				return info.WhetherExist(), err
			} else {
				k.mu.RUnlock()
				select {
				case <-k.ready:
				case <-ctx.Done():
					return false, ctx.Err()
				}
				k.mu.RLock()
				defer k.mu.RUnlock()
				info.SetInfo(k.value)
				return info.WhetherExist(), err
			}
		} else {
			c.mapMu.Lock()
			k, ok = c.items[info.GetKey()]
			if ok {
				c.mapMu.Unlock()
				info.SetInfo(k.value)
				return info.WhetherExist(), err
			}
			c.epoch++
			c.items[info.GetKey()] = &entry[dataModel]{
				ready:   make(chan struct{}),
				loading: true,
				epoch:   c.epoch,
			}
			c.mapMu.Unlock()
			c.mapMu.RLock()
			k, _ = c.items[info.GetKey()]
			c.mapMu.RUnlock()
			var err2 error
			value, exist, err2 = getDB(ctx, DBContext, info)
			k.mu.Lock()
			if err2 != nil {
				k.loading = false
				close(k.ready)
				k.mu.Unlock()
				return false, err2
			}
			if exist {
				k.value = value
			} else {
				k.value = info.GetEmptyValue()
			}
			k.loading = false
			close(k.ready)
			k.mu.Unlock()
			info.SetInfo(value)
			err2 = setCache(ctx, DBContext, info)
			if err2 != nil {
				return false, err2
			}
			k.mu.RLock()
			go DeleteCache(c, info.GetKey(), c.items[info.GetKey()].epoch)
			k.mu.RUnlock()
			return exist, err
		}
	}
	info.SetInfo(value)
	return info.WhetherExist(), nil
}

func DeleteCache[dataModel any](c *Cache[dataModel], key string, epoch int64) {
	time.Sleep(1 * time.Second)
	c.mapMu.Lock()
	info, ok := c.items[key]
	if ok && info != nil {
		if info.epoch == epoch {
			delete(c.items, key)
		}
	}
	c.mapMu.Unlock()
}
