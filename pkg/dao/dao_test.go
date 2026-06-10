package dao

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ========== 测试辅助结构体 ==========

// 测试用的数据模型
type TestDataModel struct {
	Key  string
	Val  any
	Info *TestDataModel
}

func (d *TestDataModel) GetKey() string {
	return d.Key
}

func (d *TestDataModel) GetEmptyValue() *TestDataModel {
	return &TestDataModel{
		Key:  d.Key,
		Val:  nil,
		Info: &TestDataModel{},
	}
}

func (d *TestDataModel) SetInfo(info *TestDataModel) {
	d.Info = info
}

func (d *TestDataModel) WhetherExist() bool {
	if d.Info == nil {
		return false
	}
	return d.Info.Val != nil
}

// 测试用的 DBContext
type TestDBContext struct {
	Mysql string
	Redis string
}

func (d *TestDBContext) GetClient() *TestDBContext {
	return d
}

// ========== 测试辅助函数 ==========

// 成功返回数据的 getCache
func successGetCache(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	data := &TestDataModel{
		Key:  info.GetKey(),
		Val:  "cached_value",
		Info: &TestDataModel{Val: "cached_value"},
	}
	return data, true, nil
}

// 返回缓存未命中的 getCache
func cacheMissGetCache(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	return nil, false, nil
}
func SuccessGetEmptyCache(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	data := &TestDataModel{
		Key:  info.GetKey(),
		Val:  nil,
		Info: &TestDataModel{},
	}
	return data, true, nil
}

// 缓存错误的 getCache
func cacheErrorGetCache(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	return nil, false, errors.New("redis connection failed")
}

// 成功返回数据的 getDB
func successGetDB(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	data := &TestDataModel{
		Key:  info.GetKey(),
		Val:  "db_value",
		Info: &TestDataModel{Val: "db_value"},
	}
	return data, true, nil
}

// 数据库未命中的 getDB
func dbMissGetDB(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	return nil, false, nil
}

// 数据库错误的 getDB
func dbErrorGetDB(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	return nil, false, errors.New("database connection failed")
}

// 慢查询的 getDB（模拟超时）
func slowGetDB(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	select {
	case <-time.After(500 * time.Millisecond):
		data := &TestDataModel{
			Key:  info.GetKey(),
			Val:  "slow_db_value",
			Info: &TestDataModel{Val: "slow_db_value"},
		}
		return data, true, nil
	case <-ctx.Done():
		return nil, false, ctx.Err()
	}
}

// 超时的 getDB（模拟上下文超时）
func timeoutGetDB(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	<-time.After(500 * time.Millisecond)
	return nil, false, context.DeadlineExceeded
}

// 成功的 setCache
func successSetCache(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) error {
	return nil
}

// 失败的 setCache
func errorSetCache(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) error {
	return errors.New("set cache failed")
}

// ========== 测试用例 ==========

// 1. 测试命中缓存
func TestGet_HitCache(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	info := &TestDataModel{Key: "test_key", Info: &TestDataModel{}}

	exist, err := Get[TestDataModel, TestDBContext](
		context.Background(),
		cache,
		dbCtx,
		info,
		successGetCache, // 缓存命中
		successSetCache,
		successGetDB,
	)

	if err != nil {
		t.Errorf("期望无错误，实际: %v", err)
	}
	if !exist {
		t.Errorf("期望 exist = true，实际: false")
	}
	if info.Info == nil || info.Info.Val != "cached_value" {
		t.Errorf("期望 info.Info.Val = cached_value，实际: %v", info.Info)
	}
	t.Log("✅ 命中缓存测试通过")
}

// 2. 测试缓存未命中，但数据库命中
func TestGet_CacheMiss_DBHit(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	info := &TestDataModel{Key: "test_key", Info: &TestDataModel{}}

	exist, err := Get[TestDataModel, TestDBContext](
		context.Background(),
		cache,
		dbCtx,
		info,
		cacheMissGetCache, // 缓存未命中
		successSetCache,
		successGetDB, // 数据库命中
	)

	if err != nil {
		t.Errorf("期望无错误，实际: %v", err)
	}
	if !exist {
		t.Errorf("期望 exist = true，实际: false")
	}
	if info.Info == nil || info.Info.Val != "db_value" {
		t.Errorf("期望 info.Info.Val = db_value，实际: %v", info.Info)
	}

	// 验证缓存是否被正确写入
	cache.mapMu.RLock()
	entry, ok := cache.items[info.GetKey()]
	cache.mapMu.RUnlock()
	if !ok {
		t.Errorf("期望缓存条目存在，实际不存在")
	}
	if entry == nil || entry.value == nil {
		t.Errorf("期望缓存值不为空")
	}
	t.Log("✅ 缓存未命中-数据库命中测试通过")
}

// 3. 测试缓存和数据库都未命中
func TestGet_CacheMiss_DBMiss(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	info := &TestDataModel{Key: "test_key", Info: &TestDataModel{}}

	exist, err := Get[TestDataModel, TestDBContext](
		context.Background(),
		cache,
		dbCtx,
		info,
		cacheMissGetCache,
		successSetCache,
		dbMissGetDB, // 数据库也未命中
	)

	if err != nil {
		t.Errorf("期望无错误，实际: %v", err)
	}
	if exist {
		t.Errorf("期望 exist = false，实际: true")
	}
	// 验证空值是否被缓存
	cache.mapMu.RLock()
	_, ok := cache.items[info.GetKey()]
	cache.mapMu.RUnlock()
	if !ok {
		t.Errorf("期望空值缓存条目存在，实际不存在")
	}
	t.Log("✅ 缓存和数据库都未命中测试通过")
}

// 4. 测试缓存错误
func TestGet_CacheError(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	info := &TestDataModel{Key: "test_key", Info: &TestDataModel{}}

	_, err := Get[TestDataModel, TestDBContext](
		context.Background(),
		cache,
		dbCtx,
		info,
		cacheErrorGetCache, // 缓存返回错误
		successSetCache,
		successGetDB,
	)

	if err == nil {
		t.Errorf("期望有错误，实际无错误")
		return
	}
	if err.Error() != "redis connection failed" {
		t.Errorf("期望错误信息为 'redis connection failed'，实际: %v", err)
	}
	t.Log("✅ 缓存错误测试通过")
}

// 5. 测试数据库错误
func TestGet_DBError(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	info := &TestDataModel{Key: "test_key", Info: &TestDataModel{}}

	_, err := Get[TestDataModel, TestDBContext](
		context.Background(),
		cache,
		dbCtx,
		info,
		cacheMissGetCache,
		successSetCache,
		dbErrorGetDB, // 数据库返回错误
	)

	if err == nil {
		t.Errorf("期望有错误，实际无错误")
	}
	if err.Error() != "database connection failed" {
		t.Errorf("期望错误信息为 'database connection failed'，实际: %v", err)
	}
	t.Log("✅ 数据库错误测试通过")
}

// 6. 测试并发请求（缓存击穿防护）
func TestGet_ConcurrentRequests_CacheMiss(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	key := "concurrent_key"

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// 模拟 10 个并发请求
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testInfo := &TestDataModel{Key: key, Info: &TestDataModel{}}
			exist, err := Get[TestDataModel, TestDBContext](
				context.Background(),
				cache,
				dbCtx,
				testInfo,
				cacheMissGetCache,
				successSetCache,
				successGetDB,
			)
			if err == nil && exist {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	// 所有请求都应该成功获取数据
	if successCount != 10 {
		t.Errorf("期望 10 个请求都成功，实际: %d", successCount)
	}

	// 验证只有一个 entry 被创建
	cache.mapMu.RLock()
	entryCount := len(cache.items)
	cache.mapMu.RUnlock()
	if entryCount != 1 {
		t.Errorf("期望只有 1 个缓存条目，实际: %d", entryCount)
	}
	t.Log("✅ 并发请求测试通过")
}

// 7. 测试上下文超时
func TestGet_ContextTimeout(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	info := &TestDataModel{Key: "test_key", Info: &TestDataModel{}}

	// 创建 100ms 超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := Get[TestDataModel, TestDBContext](
		ctx,
		cache,
		dbCtx,
		info,
		cacheMissGetCache,
		successSetCache,
		slowGetDB, // 慢查询，500ms
	)

	if err == nil {
		t.Errorf("期望有错误，实际无错误")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("期望超时错误，实际: %v", err)
	}
	t.Log("✅ 上下文超时测试通过")
}

// 8. 测试数据库查询超时
func TestGet_DBTimeout(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	info := &TestDataModel{Key: "test_key", Info: &TestDataModel{}}

	ctx := context.Background()

	_, err := Get[TestDataModel, TestDBContext](
		ctx,
		cache,
		dbCtx,
		info,
		cacheMissGetCache,
		successSetCache,
		timeoutGetDB, // 数据库超时
	)

	if err == nil {
		t.Errorf("期望有错误，实际无错误")
	}
	t.Log("✅ 数据库超时测试通过")
}

// 9. 测试设置缓存失败
func TestGet_SetCacheError(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	info := &TestDataModel{Key: "test_key", Info: &TestDataModel{}}

	_, err := Get[TestDataModel, TestDBContext](
		context.Background(),
		cache,
		dbCtx,
		info,
		cacheMissGetCache,
		errorSetCache, // 设置缓存失败
		successGetDB,
	)

	if err == nil {
		t.Errorf("期望有错误，实际无错误")
	}
	if err.Error() != "set cache failed" {
		t.Errorf("期望错误信息为 'set cache failed'，实际: %v", err)
	}
	t.Log("✅ 设置缓存失败测试通过")
}

// 10. 测试重复请求同一 key（等待 loading 完成）
func TestGet_RequestWhileLoading(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	key := "loading_key"

	var wg sync.WaitGroup
	startCh := make(chan struct{})

	// 第一个请求会触发加载
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-startCh
		info := &TestDataModel{Key: key, Info: &TestDataModel{}}
		exist, err := Get[TestDataModel, TestDBContext](
			context.Background(),
			cache,
			dbCtx,
			info,
			cacheMissGetCache,
			successSetCache,
			slowGetDB, // 慢查询，500ms
		)
		if err != nil || !exist {
			t.Errorf("第一个请求失败: err=%v, exist=%v", err, exist)
		}
	}()

	// 第二个请求会在第一个请求的 loading 期间等待
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-startCh
		// 延迟 50ms 发起第二个请求，确保第一个请求已进入 loading 状态
		time.Sleep(50 * time.Millisecond)
		info := &TestDataModel{Key: key, Info: &TestDataModel{}}
		exist, err := Get[TestDataModel, TestDBContext](
			context.Background(),
			cache,
			dbCtx,
			info,
			cacheMissGetCache,
			successSetCache,
			successGetDB, // 这个不会被调用，因为会等待第一个请求的结果
		)
		if err != nil || !exist {
			t.Errorf("第二个请求失败: err=%v, exist=%v", err, exist)
		}
	}()

	close(startCh)
	wg.Wait()

	// 验证只有一个缓存条目
	cache.mapMu.RLock()
	entryCount := len(cache.items)
	entry, ok := cache.items[key]
	cache.mapMu.RUnlock()

	if entryCount != 1 {
		t.Errorf("期望只有 1 个缓存条目，实际: %d", entryCount)
	}
	if ok {
		if entry.value == nil {
			t.Errorf("期望缓存值不为空")
		}
	}
	t.Log("✅ 等待加载测试通过")
}

// 11. 测试不同的 key（隔离性）
func TestGet_DifferentKeys(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}

	keys := []string{"key1", "key2", "key3"}

	for _, key := range keys {
		info := &TestDataModel{Key: key, Info: &TestDataModel{}}
		exist, err := Get[TestDataModel, TestDBContext](
			context.Background(),
			cache,
			dbCtx,
			info,
			cacheMissGetCache,
			successSetCache,
			successGetDB,
		)
		if err != nil {
			t.Errorf("key %s 失败: %v", key, err)
		}
		if !exist {
			t.Errorf("key %s 期望 exist=true", key)
		}
	}

	// 验证每个 key 都有独立的缓存条目
	cache.mapMu.RLock()
	entryCount := len(cache.items)
	cache.mapMu.RUnlock()

	if entryCount != 3 {
		t.Errorf("期望 3 个缓存条目，实际: %d", entryCount)
	}
	t.Log("✅ 不同 key 隔离性测试通过")
}

// 12. 测试空值缓存（防止缓存穿透）
func TestGet_EmptyValueCaching(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	key := "nonexistent_key"
	info := &TestDataModel{Key: key, Info: &TestDataModel{}}

	// 第一次请求：数据库也不存在
	exist, err := Get[TestDataModel, TestDBContext](
		context.Background(),
		cache,
		dbCtx,
		info,
		cacheMissGetCache,
		successSetCache,
		dbMissGetDB, // 数据库不存在
	)

	if err != nil {
		t.Errorf("第一次请求错误: %v", err)
		return
	}
	if exist {
		t.Errorf("第一次请求期望 exist=false")
		return
	}

	// 第二次请求：应该命中空值缓存
	info2 := &TestDataModel{Key: key, Info: &TestDataModel{}}
	exist2, err := Get[TestDataModel, TestDBContext](
		context.Background(),
		cache,
		dbCtx,
		info2,
		SuccessGetEmptyCache, // 这里会返回空值缓存
		successSetCache,
		successGetDB,
	)

	if err != nil {
		t.Errorf("第二次请求错误: %v", err)
		return
	}
	if exist2 {
		t.Errorf("第二次请求期望 exist=false")
		return
	}
	t.Log("✅ 空值缓存测试通过")
}

// 13. 测试缓存过期清理（DeleteCache）
func TestGet_DeleteCache(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	key := "expire_key"
	info := &TestDataModel{Key: key, Info: &TestDataModel{}}

	// 写入缓存
	_, err := Get[TestDataModel, TestDBContext](
		context.Background(),
		cache,
		dbCtx,
		info,
		cacheMissGetCache,
		successSetCache,
		successGetDB,
	)
	if err != nil {
		t.Errorf("写入缓存失败: %v", err)
	}

	// 验证缓存存在
	cache.mapMu.RLock()
	_, ok := cache.items[key]
	cache.mapMu.RUnlock()
	if !ok {
		t.Errorf("期望缓存存在")
	}

	// 触发删除（DeleteCache 会等待 1 秒）
	DeleteCache(cache, key, cache.items[key].epoch)

	// 等待删除完成
	time.Sleep(1100 * time.Millisecond)

	// 验证缓存已被删除
	cache.mapMu.RLock()
	_, ok = cache.items[key]
	cache.mapMu.RUnlock()
	if ok {
		t.Errorf("期望缓存已被删除")
	}
	t.Log("✅ 缓存过期清理测试通过")
}

// 14. 测试 context 取消
func TestGet_ContextCancel(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	info := &TestDataModel{Key: "test_key", Info: &TestDataModel{}}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	_, err := Get[TestDataModel, TestDBContext](
		ctx,
		cache,
		dbCtx,
		info,
		slowGetDB,
		successSetCache,
		slowGetDB,
	)

	if err == nil {
		t.Errorf("期望有错误")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("期望 context.Canceled，实际: %v", err)
	}
	t.Log("✅ context 取消测试通过")
}

// 15. 测试极端并发下的缓存穿透
func TestGet_ExtremeConcurrentPenetration(t *testing.T) {
	cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	key := "penetration_key"

	// 计数器：统计数据库被调用的次数
	dbCallCount := 0
	var dbMu sync.Mutex

	// 自定义 getDB 来统计调用次数
	countingGetDB := func(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
		dbMu.Lock()
		dbCallCount++
		dbMu.Unlock()
		time.Sleep(100 * time.Millisecond) // 模拟数据库查询耗时
		data := &TestDataModel{
			Key:  info.GetKey(),
			Val:  "db_value",
			Info: &TestDataModel{Val: "db_value"},
		}
		return data, true, nil
	}

	var wg sync.WaitGroup
	goroutineCount := 100

	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			info := &TestDataModel{Key: key, Info: &TestDataModel{}}
			exist, err := Get[TestDataModel, TestDBContext](
				context.Background(),
				cache,
				dbCtx,
				info,
				cacheMissGetCache,
				successSetCache,
				countingGetDB,
			)
			if err != nil || !exist {
				t.Errorf("请求失败: err=%v", err)
			}
		}()
	}
	wg.Wait()

	// 验证数据库只被调用了一次（缓存击穿防护生效）
	if dbCallCount != 1 {
		t.Errorf("期望数据库只被调用 1 次，实际: %d", dbCallCount)
	}
	t.Log("✅ 极端并发缓存穿透测试通过")
}

// ========== 压力测试辅助结构体 ==========

type testCounter struct {
	l1HitCount  int32
	l1MissCount int32
	l2HitCount  int32
	l2MissCount int32
	dbCallCount int32
}

// ========== L2 缓存模拟函数（Redis） ==========

// 模拟 L2 缓存命中（有效值）
func l2CacheHit(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	return &TestDataModel{
		Key:  info.GetKey(),
		Val:  "cached_value",
		Info: &TestDataModel{},
	}, true, nil
}

// 模拟 L2 缓存命中（空值）
func l2CacheEmptyHit(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	return &TestDataModel{
		Key:  info.GetKey(),
		Val:  nil,
		Info: &TestDataModel{},
	}, true, nil
}

// 模拟 L2 缓存未命中
func l2CacheMiss(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	return nil, false, nil
}

// 模拟 L2 缓存错误
func l2CacheError(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	return nil, false, errors.New("redis connection failed")
}

// 模拟 L2 缓存写入成功
func l2CacheSetSuccess(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) error {
	return nil
}

// 模拟 L2 缓存写入失败
func l2CacheSetError(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) error {
	return errors.New("redis write failed")
}

// ========== 数据库模拟函数 ==========

func dbSuccess(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	return &TestDataModel{
		Key:  info.GetKey(),
		Val:  "db_value",
		Info: &TestDataModel{},
	}, true, nil
}

func dbMiss(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	return nil, false, nil
}

func dbError(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	return nil, false, errors.New("database connection failed")
}

func dbSlow(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
	time.Sleep(100 * time.Millisecond)
	return &TestDataModel{
		Key:  info.GetKey(),
		Val:  "db_slow_value",
		Info: &TestDataModel{},
	}, true, nil
}

// ========== 压力测试 ==========

// 测试1：高并发 - L2缓存命中场景
func TestHighConcurrency_L2CacheHit(t *testing.T) {
	l1Cache := NewCache[TestDataModel]() // L1 本地缓存
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	//counter := &testCounter{}

	key := "hot_key"
	goroutineCount := 100

	var wg sync.WaitGroup
	startCh := make(chan struct{})

	// 真实场景：先从 L1 查，再从 L2 查，最后查 DB
	// 这里直接模拟 L2 命中
	hitCount := int32(0)

	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startCh

			info := &TestDataModel{Key: key, Info: &TestDataModel{}}
			exist, err := Get[TestDataModel, TestDBContext](
				context.Background(),
				l1Cache,
				dbCtx,
				info,
				l2CacheHit, // L2 命中
				l2CacheSetSuccess,
				dbSuccess,
			)

			if err != nil {
				t.Logf("错误: %v", err)
			}
			if exist {
				atomic.AddInt32(&hitCount, 1)
			}
		}()
	}

	close(startCh)
	wg.Wait()

	t.Logf("=== L2缓存命中压测 ===")
	t.Logf("请求数: %d", goroutineCount)
	t.Logf("成功数: %d", atomic.LoadInt32(&hitCount))

	if atomic.LoadInt32(&hitCount) != int32(goroutineCount) {
		t.Errorf("期望全部成功，实际: %d", atomic.LoadInt32(&hitCount))
	} else {
		t.Logf("✅ L2缓存命中测试通过")
	}
}

// 测试2：高并发 - L2缓存未命中，DB命中
func TestHighConcurrency_L2CacheMiss_DBHit(t *testing.T) {
	l1Cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}

	key := "miss_key"
	goroutineCount := 100

	var wg sync.WaitGroup
	startCh := make(chan struct{})
	dbCallCount := int32(0)

	getDB := func(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
		atomic.AddInt32(&dbCallCount, 1)
		time.Sleep(20 * time.Millisecond)
		return &TestDataModel{
			Key:  info.GetKey(),
			Val:  "db_value",
			Info: &TestDataModel{},
		}, true, nil
	}

	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startCh

			info := &TestDataModel{Key: key, Info: &TestDataModel{}}
			_, err := Get[TestDataModel, TestDBContext](
				context.Background(),
				l1Cache,
				dbCtx,
				info,
				l2CacheMiss, // L2 未命中
				l2CacheSetSuccess,
				getDB,
			)

			if err != nil {
				t.Logf("错误: %v", err)
			}
		}()
	}

	close(startCh)
	wg.Wait()

	t.Logf("=== L2缓存未命中，DB命中压测 ===")
	t.Logf("请求数: %d", goroutineCount)
	t.Logf("DB调用次数: %d", atomic.LoadInt32(&dbCallCount))

	// 由于 L1 缓存击穿防护，DB 应该只被调用 1 次
	if atomic.LoadInt32(&dbCallCount) != 1 {
		t.Errorf("期望 DB 调用 1 次，实际: %d", atomic.LoadInt32(&dbCallCount))
	} else {
		t.Logf("✅ L1缓存击穿防护有效，DB只被调用 %d 次", atomic.LoadInt32(&dbCallCount))
	}
}

// 测试3：高并发 - L2缓存未命中，DB也不存在（空值缓存）
func TestHighConcurrency_L2CacheMiss_DBMiss(t *testing.T) {
	l1Cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}

	key := "empty_key"
	goroutineCount := 100000

	var wg sync.WaitGroup
	startCh := make(chan struct{})
	dbCallCount := int32(0)

	getDB := func(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
		atomic.AddInt32(&dbCallCount, 1)
		return nil, false, nil
	}

	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startCh

			info := &TestDataModel{Key: key, Info: &TestDataModel{}}
			_, err := Get[TestDataModel, TestDBContext](
				context.Background(),
				l1Cache,
				dbCtx,
				info,
				l2CacheMiss,
				l2CacheSetSuccess,
				getDB,
			)

			if err != nil {
				t.Logf("错误: %v", err)
			}
		}()
	}

	close(startCh)
	wg.Wait()

	t.Logf("=== L2缓存未命中，DB不存在压测 ===")
	t.Logf("请求数: %d", goroutineCount)
	t.Logf("DB调用次数: %d", atomic.LoadInt32(&dbCallCount))

	// 空值缓存也应该只查一次 DB
	if atomic.LoadInt32(&dbCallCount) != 1 {
		t.Errorf("期望 DB 调用 1 次，实际: %d", atomic.LoadInt32(&dbCallCount))
	} else {
		t.Logf("✅ 空值缓存有效，DB只被调用 %d 次", atomic.LoadInt32(&dbCallCount))
	}
}

// 测试4：多 key 并发访问
func TestHighConcurrency_MultipleKeys(t *testing.T) {
	l1Cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}

	keyCount := 50
	goroutinesPerKey := 10000

	var wg sync.WaitGroup
	startCh := make(chan struct{})
	dbCallCount := int32(0)

	getDB := func(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
		atomic.AddInt32(&dbCallCount, 1)
		return &TestDataModel{
			Key:  info.GetKey(),
			Val:  "db_value",
			Info: &TestDataModel{},
		}, true, nil
	}

	for k := 0; k < keyCount; k++ {
		key := string(rune('a'+k%26)) + string(rune('0'+k/26))
		for g := 0; g < goroutinesPerKey; g++ {
			wg.Add(1)
			go func(key string) {
				defer wg.Done()
				<-startCh

				info := &TestDataModel{Key: key, Info: &TestDataModel{}}
				Get[TestDataModel, TestDBContext](
					context.Background(),
					l1Cache,
					dbCtx,
					info,
					l2CacheMiss,
					l2CacheSetSuccess,
					getDB,
				)
			}(key)
		}
	}

	close(startCh)
	wg.Wait()

	t.Logf("=== 多 key 并发压测 ===")
	t.Logf("Key 数量: %d", keyCount)
	t.Logf("每 Key 并发: %d", goroutinesPerKey)
	t.Logf("总请求数: %d", keyCount*goroutinesPerKey)
	t.Logf("DB调用次数: %d", atomic.LoadInt32(&dbCallCount))

	// 每个 key 应该只调用一次 DB
	expectedDBCalls := keyCount
	if atomic.LoadInt32(&dbCallCount) > int32(expectedDBCalls*2) {
		t.Errorf("DB调用次数过多: %d > %d", atomic.LoadInt32(&dbCallCount), expectedDBCalls*2)
	} else {
		t.Logf("✅ 多 key 压测通过，DB调用 %d 次", atomic.LoadInt32(&dbCallCount))
	}
}

// 测试5：性能基准 - L2缓存命中
func BenchmarkL2CacheHit(b *testing.B) {
	l1Cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}
	key := "bench_key"

	// 预热 L1 缓存
	info := &TestDataModel{Key: key, Info: &TestDataModel{}}
	Get[TestDataModel, TestDBContext](
		context.Background(),
		l1Cache,
		dbCtx,
		info,
		l2CacheHit,
		l2CacheSetSuccess,
		dbSuccess,
	)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			info := &TestDataModel{Key: key, Info: &TestDataModel{}}
			Get[TestDataModel, TestDBContext](
				context.Background(),
				l1Cache,
				dbCtx,
				info,
				l2CacheHit,
				l2CacheSetSuccess,
				dbSuccess,
			)
		}
	})
}

// 测试6：性能基准 - L2缓存未命中，DB命中
func BenchmarkL2CacheMiss_DBHit(b *testing.B) {
	l1Cache := NewCache[TestDataModel]()
	dbCtx := &TestDBContext{Mysql: "test", Redis: "test"}

	getDB := func(ctx context.Context, db DBContextInterface[TestDBContext], info Info[TestDataModel]) (*TestDataModel, bool, error) {
		return &TestDataModel{
			Key:  info.GetKey(),
			Val:  "db_value",
			Info: &TestDataModel{},
		}, true, nil
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := "bench_key"
			info := &TestDataModel{Key: key, Info: &TestDataModel{}}
			_, _ = Get[TestDataModel, TestDBContext](
				context.Background(),
				l1Cache,
				dbCtx,
				info,
				l2CacheMiss,
				l2CacheSetSuccess,
				getDB,
			)
		}
	})
}
