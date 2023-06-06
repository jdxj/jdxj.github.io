---
title: "Pool"
date: 2023-06-06T00:16:57+08:00
---

sync.Pool 数据类型用来保存一组可独立访问的临时对象。

- sync.Pool 本身就是线程安全的
- sync.Pool 不可在使用之后再复制使用
- Get 方法的返回值还可能会是一个 nil（Pool.New 字段没有设置，又没有空闲元素可以返回），所以你在使用的时候，可能需要判断
- 如果 Put 一个 nil 值，Pool 就会忽略这个值

# 实现原理

![](https://static001.geekbang.org/resource/image/f4/96/f4003704663ea081230760098f8af696.jpg?wh=3659*2186)

```go
func poolCleanup() {
    // 丢弃当前victim, STW所以不用加锁
    for _, p := range oldPools {
        p.victim = nil
        p.victimSize = 0
    }

    // 将local复制给victim, 并将原local置为nil
    for _, p := range allPools {
        p.victim = p.local
        p.victimSize = p.localSize
        p.local = nil
        p.localSize = 0
    }

    oldPools, allPools = allPools, nil
}
```

请求元素时也是优先从 local 字段中查找可用的元素

而 poolLocalInternal 也包含两个字段：private 和 shared。

- private，代表一个缓存的元素，而且只能由相应的一个 P 存取。因为一个 P 同时只能执行一个 goroutine，所以不会有并发的问题。
- shared，可以由任意的 P 访问，但是只有本地的 P 才能 pushHead/popHead，其它 P 可以 popTail，相当于只有一个本地的 P
  作为生产者（Producer），多个 P 作为消费者（Consumer），它是使用一个 local-free 的 queue 列表实现的。

## Get

```go
func (p *Pool) Get() interface{} {
    // 把当前goroutine固定在当前的P上
    l, pid := p.pin()
    x := l.private // 优先从local的private字段取，快速
    l.private = nil
    if x == nil {
        // 从当前的local.shared弹出一个，注意是从head读取并移除
        x, _ = l.shared.popHead()
        if x == nil { // 如果没有，则去偷一个
            x = p.getSlow(pid) 
        }
    }
    runtime_procUnpin()
    // 如果没有获取到，尝试使用New函数生成一个新的
    if x == nil && p.New != nil {
        x = p.New()
    }
    return x
}
```

重点是 getSlow 方法，我们来分析下。看名字也就知道了，它的耗时可能比较长。它首先要遍历所有的 local，尝试从它们的 shared 弹出一个元素。如果还没
找到一个，那么，就开始对 victim 下手了。

```go
func (p *Pool) getSlow(pid int) interface{} {

    size := atomic.LoadUintptr(&p.localSize)
    locals := p.local                       
    // 从其它proc中尝试偷取一个元素
    for i := 0; i < int(size); i++ {
        l := indexLocal(locals, (pid+i+1)%int(size))
        if x, _ := l.shared.popTail(); x != nil {
            return x
        }
    }

    // 如果其它proc也没有可用元素，那么尝试从vintim中获取
    size = atomic.LoadUintptr(&p.victimSize)
    if uintptr(pid) >= size {
        return nil
    }
    locals = p.victim
    l := indexLocal(locals, pid)
    if x := l.private; x != nil { // 同样的逻辑，先从vintim中的local private获取
        l.private = nil
        return x
    }
    for i := 0; i < int(size); i++ { // 从vintim其它proc尝试偷取
        l := indexLocal(locals, (pid+i)%int(size))
        if x, _ := l.shared.popTail(); x != nil {
            return x
        }
    }

    // 如果victim中都没有，则把这个victim标记为空，以后的查找可以快速跳过了
    atomic.StoreUintptr(&p.victimSize, 0)

    return nil
}
```

## Put

```go
func (p *Pool) Put(x interface{}) {
    if x == nil { // nil值直接丢弃
        return
    }
    l, _ := p.pin()
    if l.private == nil { // 如果本地private没有值，直接设置这个值即可
        l.private = x
        x = nil
    }
    if x != nil { // 否则加入到本地队列中
        l.shared.pushHead(x)
    }
    runtime_procUnpin()
}
```

# sync.Pool 的坑

## 内存泄漏

将容量已经变得很大的 Buffer 再放回 Pool 中，导致内存泄漏。

在使用 sync.Pool 回收 buffer 的时候，一定要检查回收的对象的大小。如果 buffer 太大，就不要回收了，否则就太浪费了。

## 内存浪费

可以将 buffer 池分成几层。

- 首先，小于 512 byte 的元素的 buffer 占一个池子；
- 其次，小于 1K byte 大小的元素占一个池子；
- 再次，小于 4K byte 大小的元素占一个池子。

这样分成几个池子以后，就可以根据需要，到所需大小的池子中获取 buffer 了。

# 连接池

- http.Client 实现连接池的代码是在 Transport 类型中
- tcp连接池[fatih/pool](https://github.com/fatih/pool)
- 标准库 sql.DB 还提供了一个通用的数据库的连接池，通过 MaxOpenConns 和 MaxIdleConns 控制最大的连接数和最大的 idle 的连接数。
- [Memcached Client](https://github.com/bradfitz/gomemcache) 连接池采用 Mutex+Slice 实现 Pool

# Worker Pool

- [fasthttp workerpool](https://github.com/valyala/fasthttp/blob/9f11af296864153ee45341d3f2fe0f5178fd6210/workerpool.go#L16)
- [workerpool](https://pkg.go.dev/github.com/gammazero/workerpool)

# 参考

- 更加通用的多层buffer池[bucketpool](https://github.com/vitessio/vitess/blob/main/go/bucketpool/bucketpool.go)
- [bytebufferpool](https://github.com/valyala/bytebufferpool)
- [bpool](https://github.com/oxtoacart/bpool)
