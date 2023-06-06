---
title: "map"
date: 2023-06-04T15:16:56+08:00
---

key必须是可比较的

- [不可比较的类型](https://gfw.go101.org/article/type-system-overview.html#:~:text=%E4%BA%8B%E5%AE%9E%EF%BC%9A-,%E5%8F%AF%E6%AF%94%E8%BE%83%E7%B1%BB%E5%9E%8B%E5%92%8C%E4%B8%8D%E5%8F%AF%E6%AF%94%E8%BE%83%E7%B1%BB%E5%9E%8B,-%E7%9B%AE%E5%89%8D%EF%BC%88Go%201.20)
- 如果 struct 的某个字段值修改了，查询 map 时无法获取它 add 进去的值

有序的map [orderedmap](https://github.com/elliotchance/orderedmap)

不能并发读写

# 实现线程安全map

- 加读写锁

{{< embedcode go "rwMap.go" >}}

- 分片加锁[concurrent-map](https://github.com/orcaman/concurrent-map)

```go
var SHARD_COUNT = 32

// 分成SHARD_COUNT个分片的map
type ConcurrentMap []*ConcurrentMapShared

// 通过RWMutex保护的线程安全的分片，包含一个map
type ConcurrentMapShared struct {
  items        map[string]interface{}
  sync.RWMutex // Read Write mutex, guards access to internal map.
}

// 创建并发map
func New() ConcurrentMap {
  m := make(ConcurrentMap, SHARD_COUNT)
  for i := 0; i < SHARD_COUNT; i++ {
    m[i] = &ConcurrentMapShared{items: make(map[string]interface{})}
  }
  return m
}

// 根据key计算分片索引
func (m ConcurrentMap) GetShard(key string) *ConcurrentMapShared {
  return m[uint(fnv32(key))%uint(SHARD_COUNT)]
}

func (m ConcurrentMap) Set(key string, value interface{}) {
  // 根据key计算出对应的分片
  shard := m.GetShard(key)
  shard.Lock() //对这个分片加锁，执行业务操作
  shard.items[key] = value
  shard.Unlock()
}

func (m ConcurrentMap) Get(key string) (interface{}, bool) {
  // 根据key计算出对应的分片
  shard := m.GetShard(key)
  shard.RLock()
  // 从这个分片读取key的值
  val, ok := shard.items[key]
  shard.RUnlock()
  return val, ok
}
```

# sync.Map

sync.Map 并不是用来替换内建的 map 类型的，它只能被应用在一些特殊的场景里。

[官方的文档](https://pkg.go.dev/sync#Map)中指出，在以下两个场景中使用 sync.Map，会比使用 map+RWMutex 的方式，性能要好得多：

- 只会增长的缓存系统中，一个 key 只写入一次而被读很多次；
- 多个 goroutine 为不相交的键集读、写和重写键值对。

这两个场景说得都比较笼统，而且，这些场景中还包含了一些特殊的情况。所以，官方建议你针对自己的场景做性能评测，如果确实能够显著提高性能，再使用
sync.Map。
