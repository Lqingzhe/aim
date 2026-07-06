package dao

import "context"

func SetWithCacheWithoutTX[dataModel any, dbContextModel any](ctx context.Context, DBContext DBContextInterface[dbContextModel], info OperateWithDBAndCache[dataModel, dbContextModel]) (err error) {
	err1 := info.SetDB(ctx, DBContext)
	err2 := info.SetCache(ctx, DBContext)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}
func SetWithCacheTx[dataModel any, dbContextModel any](ctx context.Context, DBContext DBContextInterface[dbContextModel], info OperateWithDBAndCache[dataModel, dbContextModel]) (err error) {
	err = info.SetDB(ctx, DBContext)
	if err != nil {
		return err
	}
	err = info.SetCache(ctx, DBContext)
	if err != nil {
		return err
	}
	return nil
} //尽量传入事务，在上层在执行回滚

func SetWithoutCache[dataModel any, dbContextModel any](ctx context.Context, DBContext DBContextInterface[dbContextModel], info OperateWithDB[dataModel, dbContextModel]) error {
	return info.SetDB(ctx, DBContext)
}
