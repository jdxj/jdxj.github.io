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

