---
title: "Gotchas in the Go Network Packages Defaults"
date: 2023-06-19T14:10:11+08:00
tags:
  - go
---

[原文](https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/)

# Timeouts

## Client timeouts

```go
c := &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			// This is the TCP connect timeout in this instance.
			Timeout: 2500 * time.Millisecond,
		}).DialContext,
		TLSHandshakeTimeout: 2500 * time.Millisecond,
	},
}
```

## Server timeouts

```go
s := &http.Server{
	ReadTimeout:  2500 * time.Millisecond,
	WriteTimeout: 5 * time.Second,
}
```

# HTTP Response Bodies

```go
res, err := client.Do(req)
if err != nil {
	return err
}
defer res.Body.Close()
...
```

```go
_, err := io.Copy(ioutil.Discard, res.Body)
```

# HTTP/1.x Keep-alives

不重用连接的方式

```go
client := &http.Client{
    &http.Transport{
        DisableKeepAlives: true
    }
}
```

```go
req.Close = true
```

```go
req.Header.Add("Connection", "close")
```

# Connection Pooling

```go
var DefaultTransport RoundTripper = &Transport{
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
}
...
const DefaultMaxIdleConnsPerHost = 2
```
