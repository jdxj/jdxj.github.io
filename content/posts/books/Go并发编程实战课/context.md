---
title: "Context"
date: 2023-06-06T18:32:59+08:00
---

在使用 Context 的时候，有一些约定俗成的规则。

- 一般函数使用 Context 的时候，会把这个参数放在第一个参数的位置。
- 从来不把 nil 当做 Context 类型的参数值，可以使用 context.Background() 创建一个空的上下文对象，也不要使用 nil。
- Context 只用来临时做函数之间的上下文透传，不能持久化 Context 或者把 Context 长久保存。把 Context 持久化到数据库、本地文件或者全局变量、
  缓存中都是错误的用法。
- key 的类型不应该是字符串类型或者其它内建类型，否则容易在包之间使用 Context 时候产生冲突。使用 WithValue 时，key 的类型应该是自己定义的类型。
- 常常使用 struct{}作为底层类型定义 key 的类型。对于 exported key 的静态类型，常常是接口或者指针。这样可以尽量减少内存分配。

# WithValue

```go
type valueCtx struct {
    Context
    key, val interface{}
}
```

# WithCancel

不是只有你想中途放弃，才去调用 cancel，只要你的任务正常完成了，就需要调用 cancel

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
    c := newCancelCtx(parent)
    propagateCancel(parent, &c)// 把c朝上传播
    return &c, func() { c.cancel(true, Canceled) }
}

// newCancelCtx returns an initialized cancelCtx.
func newCancelCtx(parent Context) cancelCtx {
    return cancelCtx{Context: parent}
}
```

propagateCancel 方法会顺着 parent 路径往上找，直到找到一个 cancelCtx，或者为 nil。如果不为空，就把自己加入到这个 cancelCtx 的 child，
以便这个 cancelCtx 被取消的时候通知自己。如果为空，会新起一个 goroutine，由它来监听 parent 的 Done 是否已关闭。

cancel 是向下传递的，如果一个 WithCancel 生成的 Context 被 cancel 时，如果它的子 Context（也有可能是孙，或者更低，依赖子的类型）也是
cancelCtx 类型的，就会被 cancel，但是不会向上传递。parent Context 不会因为子 Context 被 cancel 而 cancel。

# WithTimeout

WithTimeout 其实是和 WithDeadline 一样，只不过一个参数是超时时间，一个参数是截止时间。超时时间加上当前时间，其实就是截止时间，因此，
WithTimeout 的实现是：

```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
    // 当前时间+timeout就是deadline
    return WithDeadline(parent, time.Now().Add(timeout))
}
```

# WithDeadline

```go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
    // 如果parent的截止时间更早，直接返回一个cancelCtx即可
    if cur, ok := parent.Deadline(); ok && cur.Before(d) {
        return WithCancel(parent)
    }
    c := &timerCtx{
        cancelCtx: newCancelCtx(parent),
        deadline:  d,
    }
    propagateCancel(parent, c) // 同cancelCtx的处理逻辑
    dur := time.Until(d)
    if dur <= 0 { //当前时间已经超过了截止时间，直接cancel
        c.cancel(true, DeadlineExceeded)
        return c, func() { c.cancel(false, Canceled) }
    }
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.err == nil {
        // 设置一个定时器，到截止时间后取消
        c.timer = time.AfterFunc(dur, func() {
            c.cancel(true, DeadlineExceeded)
        })
    }
    return c, func() { c.cancel(true, Canceled) }
}
```

