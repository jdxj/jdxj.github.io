---
title: "mutex"
date: 2023-04-21T10:38:19+08:00
---

# 演化过程

![](https://static001.geekbang.org/resource/image/c2/35/c28531b47ff7f220d5bc3c9650180835.jpg?wh=4000*2250)

# 初版

```go
// CAS操作，当时还没有抽象出atomic包
func cas(val *int32, old, new int32) bool
func semacquire(*int32)
func semrelease(*int32)
// 互斥锁的结构，包含两个字段
type Mutex struct {
    key  int32 // 锁是否被持有的标识
    sema int32 // 信号量专用，用以阻塞/唤醒goroutine
}

// 保证成功在val上增加delta的值
func xadd(val *int32, delta int32) (new int32) {
    for {
        v := *val
        if cas(val, v, v+delta) {
            return v + delta
        }
    }
    panic("unreached")
}

// 请求锁
func (m *Mutex) Lock() {
    if xadd(&m.key, 1) == 1 { //标识加1，如果等于1，成功获取到锁
        return
    }
    semacquire(&m.sema) // 否则阻塞等待
}

func (m *Mutex) Unlock() {
    if xadd(&m.key, -1) == 0 { // 将标识减去1，如果等于0，则没有其它等待者
        return
    }
    semrelease(&m.sema) // 唤醒其它阻塞的goroutine
}    
```

![](https://static001.geekbang.org/resource/image/82/25/825e23e1af96e78f3773e0b45de38e25.jpg?wh=3032*2199)

Unlock 方法可以被任意的 goroutine 调用释放锁，即使是没持有这个互斥锁的 goroutine，也可以进行这个操作。这是因为，Mutex
本身并没有包含持有这把锁的 goroutine 的信息，所以，Unlock 也不会对此进行检查。Mutex 的这个设计一直保持至今。

# 给新人机会

```go
type Mutex struct {
    state int32
    sema  uint32
}

const (
    mutexLocked = 1 << iota // mutex is locked
    mutexWoken // 是否存在被唤醒的goroutine
    mutexWaiterShift = iota
)
```

![](https://static001.geekbang.org/resource/image/4c/15/4c4a3dd2310059821f41af7b84925615.jpg?wh=2978*1756)

```go
func (m *Mutex) Lock() {
    // Fast path: 幸运case，能够直接获取到锁
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        return
    }

    awoke := false
    for {
        old := m.state
        new := old | mutexLocked // 新状态加锁
        if old&mutexLocked != 0 { // 锁未释放
            new = old + 1<<mutexWaiterShift //等待者数量加一
        }
        if awoke {
            // goroutine是被唤醒的，
            // 新状态清除唤醒标志
            new &^= mutexWoken
        }
        if atomic.CompareAndSwapInt32(&m.state, old, new) {//设置新状态
            if old&mutexLocked == 0 { // 锁原状态未加锁
                break
            }
            runtime.Semacquire(&m.sema) // 请求信号量
            awoke = true
        }
    }
}
```

[&^含义](https://zhuanlan.zhihu.com/p/86765504)

![](https://static001.geekbang.org/resource/image/de/ac/de281d2a9e022b2e026bb6126f28c9ac.jpg?wh=1690x620)

```go
func (m *Mutex) Unlock() {
    // Fast path: drop lock bit.
    new := atomic.AddInt32(&m.state, -mutexLocked) //去掉锁标志
    if (new+mutexLocked)&mutexLocked == 0 { //本来就没有加锁
        panic("sync: unlock of unlocked mutex")
    }

    old := new
    for {
        if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken) != 0 { // 没有等待者，或者有唤醒的waiter，或者锁原来已加锁
            return
        }
        new = (old - 1<<mutexWaiterShift) | mutexWoken // 新状态，准备唤醒goroutine，并设置唤醒标志
        if atomic.CompareAndSwapInt32(&m.state, old, new) {
            runtime.Semrelease(&m.sema)
            return
        }
        old = m.state
    }
}
```

# 多给些机会

```go
func (m *Mutex) Lock() {
    // Fast path: 幸运之路，正好获取到锁
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        return
    }

    awoke := false
    iter := 0
    for { // 不管是新来的请求锁的goroutine, 还是被唤醒的goroutine，都不断尝试请求锁
        old := m.state // 先保存当前锁的状态
        new := old | mutexLocked // 新状态设置加锁标志
        if old&mutexLocked != 0 { // 锁还没被释放
            if runtime_canSpin(iter) { // 还可以自旋
                if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
                    atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
                    awoke = true
                }
                runtime_doSpin()
                iter++
                continue // 自旋，再次尝试请求锁
            }
            new = old + 1<<mutexWaiterShift
        }
        if awoke { // 唤醒状态
            if new&mutexWoken == 0 {
                panic("sync: inconsistent mutex state")
            }
            new &^= mutexWoken // 新状态清除唤醒标记
        }
        if atomic.CompareAndSwapInt32(&m.state, old, new) {
            if old&mutexLocked == 0 { // 旧状态锁已释放，新状态成功持有了锁，直接返回
                break
            }
            runtime_Semacquire(&m.sema) // 阻塞等待
            awoke = true // 被唤醒
            iter = 0
        }
    }
}
```

# 解决饥饿

在极端情况下，等待中的 goroutine 可能会一直获取不到锁

![](https://static001.geekbang.org/resource/image/e0/76/e0c23794c8a1d355a7a183400c036276.jpg?wh=3096*1644)

```go
type Mutex struct {
    state int32
    sema  uint32
}

const (
    mutexLocked = 1 << iota // mutex is locked
    mutexWoken
    mutexStarving // 从state字段中分出一个饥饿标记
    mutexWaiterShift = iota

    starvationThresholdNs = 1e6
)

func (m *Mutex) Lock() {
    // Fast path: 幸运之路，一下就获取到了锁
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        return
    }
    // Slow path：缓慢之路，尝试自旋竞争或饥饿状态下饥饿goroutine竞争
    m.lockSlow()
}

func (m *Mutex) lockSlow() {
    var waitStartTime int64
    starving := false // 此goroutine的饥饿标记
    awoke := false // 唤醒标记
    iter := 0 // 自旋次数
    old := m.state // 当前的锁的状态
    for {
        // 锁是非饥饿状态，锁还没被释放，尝试自旋
        if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
            if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
                atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
                awoke = true
            }
            runtime_doSpin()
            iter++
            old = m.state // 再次获取锁的状态，之后会检查是否锁被释放了
            continue
        }
        new := old
        if old&mutexStarving == 0 {
            new |= mutexLocked // 非饥饿状态，加锁
        }
        if old&(mutexLocked|mutexStarving) != 0 {
            new += 1 << mutexWaiterShift // waiter数量加1
        }
        if starving && old&mutexLocked != 0 {
            new |= mutexStarving // 设置饥饿状态
        }
        if awoke {
            if new&mutexWoken == 0 {
                throw("sync: inconsistent mutex state")
            }
            new &^= mutexWoken // 新状态清除唤醒标记
        }
        // 成功设置新状态
        if atomic.CompareAndSwapInt32(&m.state, old, new) {
            // 原来锁的状态已释放，并且不是饥饿状态，正常请求到了锁，返回
            if old&(mutexLocked|mutexStarving) == 0 {
                break // locked the mutex with CAS
            }
            // 处理饥饿状态

            // 如果以前就在队列里面，加入到队列头
            queueLifo := waitStartTime != 0
            if waitStartTime == 0 {
                waitStartTime = runtime_nanotime()
            }
            // 阻塞等待
            runtime_SemacquireMutex(&m.sema, queueLifo, 1)
            // 唤醒之后检查锁是否应该处于饥饿状态
            starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
            old = m.state
            // 如果锁已经处于饥饿状态，直接抢到锁，返回
            if old&mutexStarving != 0 {
                if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
                    throw("sync: inconsistent mutex state")
                }
                // 有点绕，加锁并且将waiter数减1
                delta := int32(mutexLocked - 1<<mutexWaiterShift)
                if !starving || old>>mutexWaiterShift == 1 {
                    delta -= mutexStarving // 最后一个waiter或者已经不饥饿了，清除饥饿标记
                }
                atomic.AddInt32(&m.state, delta)
                break
            }
            awoke = true
            iter = 0
        } else {
            old = m.state
        }
    }
}

func (m *Mutex) Unlock() {
    // Fast path: drop lock bit.
    new := atomic.AddInt32(&m.state, -mutexLocked)
    if new != 0 {
        m.unlockSlow(new)
    }
}

func (m *Mutex) unlockSlow(new int32) {
    if (new+mutexLocked)&mutexLocked == 0 {
        throw("sync: unlock of unlocked mutex")
    }
    if new&mutexStarving == 0 {
        old := new
        for {
            if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
                return
            }
            new = (old - 1<<mutexWaiterShift) | mutexWoken
            if atomic.CompareAndSwapInt32(&m.state, old, new) {
                runtime_Semrelease(&m.sema, false, 1)
                return
            }
            old = m.state
        }
    } else {
        runtime_Semrelease(&m.sema, true, 1)
    }
}
```

# vet发现Mutex复制原理

检查是通过[copylock](https://github.com/golang/tools/blob/master/go/analysis/passes/copylock/copylock.go)分析器静态分析实现的。
这个分析器会分析函数调用、range 遍历、复制、声明、函数返回值等位置，有没有锁的值 copy 的情景，以此来判断有没有问题。可以说，只要是实现了
Locker 接口，就会被分析。
