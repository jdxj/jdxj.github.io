---
title: "Designing Go Libraries"
date: 2022-12-20T10:17:23+08:00
draft: false
tags:
  - go
---

## [原文](https://abhinavg.net/2022/12/06/designing-go-libraries/)

# 1. Primary Concerns

## 1.1. Usability

- 建立惯例使库的特性可发现
- 潜在的错误使用
- 易于完成常见任务

### Case Study: net/http

看起来有些繁琐
```go {title="Sending a GET request"}
// import "net/http"

req, err := http.NewRequest(http.MethodGet, "http://example.com", nil /* body */)
if err != nil {
  return err
}

var client http.Client
res, err := client.Do(req)
// ...
```

更简单的方法
```go {title="Sending a GET request—easiest way"}
// import "net/http"

res, err := http.Get("http://example.com")
```

> 但是不应该使用全局 HTTP Client.

## 1.2. Readability

在实现之前写伪代码和文档有助于设计出可读性好的 API. 随着经验的增加可能会减少这个习惯.

## 1.3. Flexibility

灵活性决定能否加入新功能来满足新用例, 允许用户自定义和扩展来适应他们的需要.
- 生态系统
- 第三方扩展

### Case Study: net/http

```go
var client http.Client
res, err := client.Get("http://example.com")
```

```go
package http

type Client struct {
  // Transport specifies the mechanism by which individual
  // HTTP requests are made.
  // If nil, DefaultTransport is used.
  Transport RoundTripper
  // ...
}

type RoundTripper interface {
  RoundTrip(*Request) (*Response, error)
}
```

封装一层, 添加日志功能
```go
type loggingRT struct {
  http.RoundTripper
}

func (rt *loggingRT) RoundTrip(req *http.Request) (*http.Response, error) {
  log.Printf("%v %v", req.Method, req.URL)
  return rt.RoundTripper.RoundTrip(req)
}

roundTripper := &loggingRT{
  RoundTripper: http.DefaultTransport,
}
client := http.Client{
  Transport: roundTripper,
}
res, err := client.Get("http://example.com")
```

## 1.4. Testability

不要事后想到可测试性, 而是在设计时.





# 2. Backwards compatibility

新版本软件可以使用旧版本软件的数据.

- 使用某个库的重要因素是该库承诺后向兼容
- 绝对的后向兼容有缺点
  - 不能使用语言特性
  - 不能使用标准库的新 API
  - 偿还技术债务

设置向后兼容的范围, 例子
- 维护当前和前一版本的 releases
- 遵循 Semantic versioning

## 2.1. Breaking changes

[Go 1 and the Future of Go Programs](https://go.dev/doc/go1compat) 给出了一些 breaking changes 定义

向公开的 interface 添加方法是 breaking change

## 2.2. Semantic versioning

1.0之前被视为不稳定版本




# 3. Recommendations

## 3.1. Work backwards

在实现之前要考虑
- API 可能的使用方法
- 能否被误用
- 如何测试功能符合预期
- 考虑灵活性以适应未来可能的新需求
- API 对输入的要求是什么, 对输出的保证是什么

## 3.2. Minimize surface area

![](library-surface-area.svg)

- 公开的"表面积"越小, 你的库在实现上会获得很大灵活性
- 公开的"表面积"越大, 需要保证稳定的东西越多

### 3.2.1 Internal packages

使用 internal 包减少表面积
- 不要在公开的 API 中包含 internal 中的实体

## 3.3. Avoid unknown outputs

- 对输入进行严格检查
- 必要情况下 copy 传入的 slice/map, 以免受元素内部状态变更所导致的问题
- 函数出错时返回类型的零值&err
- 考虑返回的 slice 的元素顺序

## 3.4. No global state

全局状态的缺点
- 不易测试
- 高耦合
- 减少灵活性

Bad
```go
var cache map[string]*Value 

func init() {
  cache = make(map[string]*Value) 
}

func Lookup(name string) (*Value, error) {
  if v, ok := cache[name]; ok {
    return v, nil 
  }
  // ...
  v := /* ... */
  cache[name] = v 
  return v, nil
}
```

Good
```go
type Handle struct { 
  cache map[string]*Value 
}

func NewHandle() *Handle { 
  return &Handle{
    cache: make(map[string]*Value),
  }
}

func (h *Handle) Lookup(name string) (*Value, error) { 
  if v, ok := h.cache[name]; ok {
    return v, nil
  }
  // ...
  v := /* ... */
  h.cache[name] = v
  return v, nil
}
```

## 3.5. Accept, don’t instantiate

Bad
```go
func New(fname string) (*Parser, error) {
  f, err := os.Open(fname)
  // ...
}
```

Good
```go
func New(f *os.File) (*Parser, error) { 
  // ...
}
```

## 3.6. Accept interfaces

Bad
```go
func New(f *os.File) (*Parser, error) {
  // ...
}
```

Good
```go
func New(r io.Reader) (*Parser, error) {
  // ...
}
```

组合接口
```go
type Source interface {
  io.Reader

  Name() string
}

var _ Source = (*os.File)(nil) 

func New(src Source) (*Parser, error) {
  // ...
}
```

## 3.7. Interfaces are forever

对接口的增删改都是 break change.

## 3.8. Return structs

返回 struct 而不是 interface 可以获得灵活性和向后兼容能力

Bad
```go
type Client interface { 
  Set(k string, v []byte) error
  Get(k string) ([]byte, error)
}

type clientImpl struct {  
  // ...
}
func (c *clientImpl) Set(...) error
func (c *clientImpl) Get(...) ([]byte, error)

func New(/* ... */) Client { 
  return &clientImpl{
    // ...
  }
}
```

Good
```go
type Client struct { 
  // ...
}
func (c *Client) Set(...) error
func (c *Client) Get(...) ([]byte, error)

func New(/* ... */) *Client { 
  return &Client{
    // ...
  }
}
```
## 3.9. Upgrade with upcasting

不能给公开的 interface 添加方法
```go
 type Source interface {
   io.Reader

   Name() string
+  Offset() int64 // bad: breaking change
 }
```

要创建新的 interface
```go
type OffsetSource interface {
  Source

  Offset() int64
}
```

向后兼容
```go
func New(src Source) *Parser {
  osrc, ok := src.(OffsetSource)
  if !ok {
    osrc = &nopOffsetSource{src} 
  }

  return &Parser{
    osrc: osrc,
    // ...
  }
}

type nopOffsetSource struct{ Source }

func (*nopOffsetSource) Offset() int64 {
  return 0
}
```

## 3.10. Parameter objects

Bad
```go
func New(url string) *Client { 
  // ...
}
```

Good
```go
type Config struct { 
  URL string
}

func New(c *Config) *Client { 
  // ...
}
```

## 3.11. Functional options

```go
package db

type Option /* ... */ 

func Connect(addr string, opts ...Option) (*Connection, error)

func WithTimeout(time.Duration) Option { /* ... */ }

func WithCache() Option { /* ... */ }
```

```go
db.Connect(addr)
db.Connect(addr, db.WithTimeout(time.Second))
db.Connect(addr, db.WithCache())
db.Connect(addr,
  db.WithTimeout(time.Second),
  db.WithCache(),
)
```

### 3.11.1. How to implement functional options

```go
// 定义 option
type connectOptions struct {
  timeout time.Duration
  cache   bool
}

// 定义操作 option 的接口
type Option interface {
  apply(*connectOptions)
}

// 实现 option 接口
type timeoutOption struct{ d time.Duration }
func (t timeoutOption) apply(o *connectOptions) {
  o.timeout = t.d
}

func WithTimeout(d time.Duration) Option {
  return timeoutOption{d: d} 
}

// 默认值
func Connect(addr string, os ...Option) (*Connection, error) {
  opts := connectOptions{
    timeout: time.Second,
  }
  for _, o := range os {
    o.apply(&opts)
  }

// ... 
}
```

### 3.11.2. Planning for functional options

预留
```go
// Option configures the behavior of Connect.
//
// There are no options at this time.
type Option interface {
  unimplemented() 
}

func Connect(addr string, opts ..Option) *Connection
```

## 3.12. Result objects

向后兼容的
```go
type UpsertResponse struct { 
  Entries  []*Entry
}

func (c *Client) Upsert(ctx context.Context, req *UpsertRequest) (*UpsertResponse, error) { 
  // ...
  return &UpsertResponse{
    Entries:  entries,
  }, nil
}
```

## 3.13. Errors

- 要么返回 err, 要么打印 log, 不要都有
- 不要使用 [pkg/errors](https://github.com/pkg/errors), 有性能问题

## 3.14. Goroutines

不要无限启 goroutine
- 处理一个请求所启动的 goroutine 数量应该与请求的内容无关
- 使用 goroutine 池
注意 goroutine 泄露
- 应该能够 graceful stop goroutine
- [goleak](https://pkg.go.dev/go.uber.org/goleak)

## 3.15. Reflection

谨慎地使用反射.

## 3.16. Naming

- [Effective Go > Names](https://go.dev/doc/effective_go#names)
- [What’s in a name?](https://go.dev/talks/2014/names.slide#1)

不要使用 common, util.

## 3.17. Documentation

有文档也是采用某个库的重要原因.

- 使用文档而不是开发文档
- 使用段落和列表突出信息
- 提供例子
- 不要写没用的注释

## 3.18. Keep a changelog

- 面向用户的 changelog
- [changelog 格式](https://keepachangelog.com/en/1.0.0/)
