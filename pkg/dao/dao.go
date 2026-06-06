package commondao

import (
	"context"
	"sync"
	"time"
)

type entry struct {
	mu      sync.RWMutex
	loading bool
	ready   chan struct{}
	value   any
	epoch   int64
}

type L1Cache struct {
	items map[string]*entry
	mapMu       sync.RWMutex
	globalEpoch int64
}
type DBContext interface {
	GetDBContext(string) any
}
type Info interface {
	GetKey() string
	SetInfo(any)
	SetGetinfo(any)
	WhetherExist() bool
	GetEmptyValue() any
}

func NewL1Cache() *L1Cache {
	return &L1Cache{
		items: make(map[string]*entry),
	}
}

func DeleteCache(c *L1Cache, key string, epoch int64) {
	time.Sleep(1 * time.Second)
	c.mapMu.Lock()
	if c.items[key].epoch == epoch {
		delete(c.items, key)
	}
	c.mapMu.Unlock()
}
func Get[I Info, D DBContext, rawInfo any](ctx context.Context, c *L1Cache, dbContext D, info I, getCache func(context.Context, D, I) (rawInfo, bool, error), setCache func(context.Context, D, I) error, getDB func(context.Context, D, I) (rawInfo, bool, error)) (bool, error) {

	value, exist, err := getCache(ctx, dbContext, info)
	if !exist || err != nil {
		c.mapMu.RLock()
		k, ok := c.items[info.GetKey()]
		c.mapMu.RUnlock()
		if ok {
			k.mu.RLock()
			if !k.loading {
				defer k.mu.RUnlock()
				info.SetGetinfo(k.value)
				return info.WhetherExist(), nil
			} else {
				k.mu.RUnlock()
				select {
				case <-k.ready:
				case <-ctx.Done():
					return false, ctx.Err()
				}
				k.mu.RLock()
				defer k.mu.RUnlock()
				info.SetGetinfo(k.value)
				return info.WhetherExist(), err
			}
		} else {
			c.mapMu.Lock()
			c.globalEpoch++
			c.items[info.GetKey()] = &entry{
				ready:   make(chan struct{}),
				loading: true,
				epoch:   c.globalEpoch,
			}
			c.mapMu.Unlock()
			c.mapMu.RLock()
			k, _ = c.items[info.GetKey()]
			c.mapMu.RUnlock()
			value, exist, err = getDB(ctx, dbContext, info)
			k.mu.Lock()
			if err != nil {
				k.loading = false
				close(k.ready)
				k.mu.Unlock()
				return false, err
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
			err2 := setCache(ctx, dbContext, info)
			if err2 != nil {
				return false, err2
			}
			go DeleteCache(c, info.GetKey(), c.items[info.GetKey()].epoch)
			return exist, err
		}
	}
	info.SetGetinfo(value)
	return info.WhetherExist(), nil
}
