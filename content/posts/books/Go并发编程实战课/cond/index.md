---
title: "Cond 条件变量"
date: 2023-05-30T12:49:56+08:00
---

Go 标准库提供 Cond 原语的目的是，为等待/通知场景下的并发问题提供支持。

# Cond 的基本用法

```go
type Cond
  func NeWCond(l Locker) *Cond
  func (c *Cond) Broadcast()
  func (c *Cond) Signal()
  func (c *Cond) Wait()
```

- Cond 关联的 Locker 实例可以通过 c.L 访问，它内部维护着一个先入先出的等待队列。
- Signal 方法，允许调用者 Caller 唤醒一个等待此 Cond 的 goroutine。如果此时没有等待的 goroutine，显然无需通知 waiter；如果 Cond 等待队
  列中有一个或者多个等待的 goroutine，则需要从等待队列中移除第一个 goroutine 并把它唤醒。
  - 调用 Signal 方法时，不强求你一定要持有 c.L 的锁。
- Broadcast 方法，允许调用者 Caller 唤醒所有等待此 Cond 的 goroutine。如果此时没有等待的 goroutine，显然无需通知 waiter；如果 Cond 等
  待队列中有一个或者多个等待的 goroutine，则清空所有等待的 goroutine，并全部唤醒。
  - 调用 Broadcast 方法时，也不强求你一定持有 c.L 的锁。
- Wait 方法，会把调用者 Caller 放入 Cond 的等待队列中并阻塞，直到被 Signal 或者 Broadcast 的方法从等待队列中移除并唤醒。
  - **调用 Wait 方法时必须要持有 c.L 的锁。**

```go
func main() {
    c := sync.NewCond(&sync.Mutex{})
    var ready int

    for i := 0; i < 10; i++ {
        go func(i int) {
            time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)

            // 加锁更改等待条件
            c.L.Lock()
            ready++
            c.L.Unlock()

            log.Printf("运动员#%d 已准备就绪\n", i)
            // 广播唤醒所有的等待者
            c.Broadcast()
        }(i)
    }

    c.L.Lock()
    for ready != 10 {
        c.Wait()
        log.Println("裁判员被唤醒一次")
    }
    c.L.Unlock()

    //所有的运动员是否就绪
    log.Println("所有运动员都准备就绪。比赛开始，3，2，1, ......")
}
```

Cond 的使用其实没那么简单。它的复杂在于：

- 这段代码有时候需要加锁，有时候可以不加；
- Wait 唤醒后需要检查条件；
- 条件变量的更改，其实是需要原子操作或者互斥锁保护的。

# Cond 的实现原理

```go
type Cond struct {
    noCopy noCopy

    // 当观察或者修改等待条件的时候需要加锁
    L Locker

    // 等待队列
    notify  notifyList
    checker copyChecker
}

func NewCond(l Locker) *Cond {
    return &Cond{L: l}
}

func (c *Cond) Wait() {
    c.checker.check()
    // 增加到等待队列中
    t := runtime_notifyListAdd(&c.notify)
    c.L.Unlock()
    // 阻塞休眠直到被唤醒
    runtime_notifyListWait(&c.notify, t)
    c.L.Lock()
}

func (c *Cond) Signal() {
    c.checker.check()
    runtime_notifyListNotifyOne(&c.notify)
}

func (c *Cond) Broadcast() {
    c.checker.check()
    runtime_notifyListNotifyAll(&c.notify)
}
```

- runtime_notifyListXXX 是运行时实现的方法，实现了一个等待 / 通知的队列。
- copyChecker 是一个辅助结构，可以在运行时检查 Cond 是否被复制使用。
- Signal 和 Broadcast 只涉及到 notifyList 数据结构，不涉及到锁。
- Wait 把调用者加入到等待队列时会释放锁，在被唤醒之后还会请求锁。在阻塞休眠期间，调用者是不持有锁的，这样能让其他 goroutine 有机会检查或者更新
  等待变量。

# Cond为什么不能被Channel替代

- Cond 和一个 Locker 关联，可以利用这个 Locker 对相关的依赖条件更改提供保护。
- Cond 可以**同时**支持 Signal 和 Broadcast 方法，而 Channel 只能同时支持其中一种。
- Cond 的 Broadcast 方法可以被重复调用。等待条件再次变成不满足的状态后，我们又可以调用 Broadcast 再次唤醒等待的 goroutine。这也是
  Channel 不能支持的，Channel 被 close 掉了之后不支持再 open。

# 有限容量队列实现

{{< embedcode go "cap-limited-queue.go" >}}
