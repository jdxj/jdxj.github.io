---
title: "curl"
date: 2022-12-06T10:28:45+08:00
draft: false
---

# Options

## `-K, --config <file>`

在文件中指定 curl 的参数

- 选项和其参数可以用 `空格`, `:`, `=` 分隔
- 如果选项前有 `-` 或 `--`, 则可以省略 `:`, `=`
- 长选项可省略 `--`
- 参数中包含空格 (或者以 `:`, `=` 开头), 则需要用引号包围

## `-x, --proxy [protocol://]host[:port]`

## `-A, --user-agent <name>`

- `-A ""`, 请求中将不会有该 header
- `-A " "`, 有该 header, 但是值为空白

## `-H, --header <header/@file>`

清除某个 header
```
-H "Host:"
```

发送不带值的自定义 header, 必须以 `;` 结尾
```
-H "X-Custom-Header;"
```

使用 `@filename` 可以在文件中读取 header
```
// todo 例子
```

## `-b, --cookie <data|filename>`

在命令行中所指定的 cookie 的格式
```
"NAME1=VALUE1; NAME2=VALUE2"
```

从文件中读取
- 文件的 cookie 格式应该是 [Netscape/Mozilla cookie file format](https://everything.curl.dev/http/cookies/fileformat)
```shell
curl -b cookiefile https://example.com
```

从文件中读取 cookie 并写回
```shell
curl -b cookiefile -c cookiefile https://example.com
```

## `-c, --cookie-jar <filename>`

存储 cookie



## `-L, --location`

跟随重定向



## `-d, --data <data>`

发送 `application/x-www-form-urlencoded` 数据
```shell
curl -d "name=curl" -d "tool=cmdline" https://example.com
curl -d @filename https://example.com
```

## `--data-binary <data>`

发送 `application/x-www-form-urlencoded` 数据



## `-F, --form <name=content>`

发送 `multipart/form-data` 数据

`@`, 作为文件上传
```shell
curl -F profile=@portrait.jpg https://example.com/upload.cgi
# Content-Disposition: form-data; name="file"; filename="nameinpost"
curl -F "file=@localfile;filename=nameinpost" example.com
```

发送 form
```shell
curl -F name=John -F shoesize=11 https://example.com/
```

`<`, 从文件中读 form
```shell
curl -F "story=<hugefile.txt" https://example.com/
```



## `-e, --referer <URL>`

添加 referer header
```shell
curl --referer "https://fake.example;auto" -L https://example.com
```



## `-G, --get`

更改 HTTP 方法为 GET



## `-i, --include`

输出包括响应 header



## `-v, --verbose`

更多传输信息



## `-k, --insecure`

跳过证书检查



## `-w, --write-out <format>`



## `-D, --dump-header <filename>`

保存 header 到文件
```shell
curl --json '{ "drink":' --json ' "coffe" }' https://example.com
curl --json @prepared https://example.com
```


## `--json <data>`


# Files

## `~/.curlrc`

# Environment

# Proxy protocol prefixes