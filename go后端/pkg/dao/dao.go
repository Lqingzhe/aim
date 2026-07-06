package dao

import (
	"context"
	"sync"
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
type info[dataModel any] interface {
	GetKey() string
	GetEmptyValue() *dataModel
	SetInfo(*dataModel)
	WhetherExist() bool
}
type cacheOperation[dataModel any, dbContextModel any] interface {
	GetCache(context.Context, DBContextInterface[dbContextModel]) (*dataModel, bool, error)
	SetCache(context.Context, DBContextInterface[dbContextModel]) error
	DeleteCache(context.Context, DBContextInterface[dbContextModel]) error
}

type dbOperation[dataModel any, dbContextModel any] interface {
	GetDB(context.Context, DBContextInterface[dbContextModel]) (*dataModel, bool, error)
	SetDB(context.Context, DBContextInterface[dbContextModel]) error
	DeleteDB(context.Context, DBContextInterface[dbContextModel]) error
	UpdateDB(context.Context, DBContextInterface[dbContextModel]) (bool, error)
}

type OperateWithDBAndCache[dataModel any, dbContextModel any] interface {
	info[dataModel]
	dbOperation[dataModel, dbContextModel]
	cacheOperation[dataModel, dbContextModel]
}
type OperateWithDB[dataModel any, dbContextModel any] interface {
	info[dataModel]
	dbOperation[dataModel, dbContextModel]
}

type DBSetter[dbContextModel any] interface {
	Set(context.Context, *dbContextModel) error
}

func Set[dbContextModel any](ctx context.Context, DBContext DBContextInterface[dbContextModel], Info DBSetter[dbContextModel]) error {
	return Info.Set(ctx, DBContext.GetClient())
}

type DBGetter[dataModel any, dbContextModel any] interface {
	Get(context.Context, *dbContextModel) (*dataModel, bool, error)
}

func Get[dataModel any, dbContextModel any](ctx context.Context, DBContext DBContextInterface[dbContextModel], Info DBGetter[dataModel, dbContextModel]) (*dataModel, bool, error) {
	return Info.Get(ctx, DBContext.GetClient())
}

type DBUpdater[dbContextModel any] interface {
	Update(context.Context, *dbContextModel) (bool, error)
}

func Update[dbContextModel any](ctx context.Context, DBContext DBContextInterface[dbContextModel], Info DBUpdater[dbContextModel]) (bool, error) {
	return Info.Update(ctx, DBContext.GetClient())
}

type DBDeleter[dbContextModel any] interface {
	Delete(context.Context, *dbContextModel) error
}

func Delete[dbContextModel any](ctx context.Context, DBContext DBContextInterface[dbContextModel], Info DBDeleter[dbContextModel]) error {
	return Info.Delete(ctx, DBContext.GetClient())
}
