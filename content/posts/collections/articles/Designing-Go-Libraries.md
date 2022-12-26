---
title: "Designing Go Libraries"
date: 2022-12-20T10:17:23+08:00
draft: true
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

