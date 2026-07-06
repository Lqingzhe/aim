package dao

import (
	"context"
	"time"
)

func GetWithCache[dataModel any, dbContextModel any](ctx context.Context, c *Cache[dataModel], DBContext DBContextInterface[dbContextModel], info OperateWithDBAndCache[dataModel, dbContextModel]) (bool, error) {
	value, exist, err := info.GetCache(ctx, DBContext)
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
			value, exist, err2 = info.GetDB(ctx, DBContext)
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
			err2 = info.SetCache(ctx, DBContext)
			if err2 != nil {
				return false, err2
			}
			k.mu.RLock()
			go deleteL1Cache(c, info.GetKey(), c.items[info.GetKey()].epoch)
			k.mu.RUnlock()
			return exist, err
		}
	}
	info.SetInfo(value)
	return info.WhetherExist(), nil
}

func deleteL1Cache[dataModel any](c *Cache[dataModel], key string, epoch int64) {
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

func GetWithOutCache[dataModel any, dbContextModel any](ctx context.Context, DBContext DBContextInterface[dbContextModel], info OperateWithDB[dataModel, dbContextModel]) (bool, error) {
	getInfo, exist, err := info.GetDB(ctx, DBContext)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	info.SetInfo(getInfo)
	return true, nil
}
