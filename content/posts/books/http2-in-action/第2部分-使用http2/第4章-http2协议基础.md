---
title: "第4章 HTTP/2协议基础"
date: 2023-03-21T21:26:37+08:00
draft: true
---

## 4.1 为什么是HTTP/2而不是HTTP/1.2

新版本的协议与原来的协议有很大的不同，新增了如下概念：

- 二进制协议
- 多路复用
- 流量控制功能
- 数据流优先级
- 首部压缩
- 服务端推送

HTTP/1.0的Web服务器可以支持HTTP/1.1的消息，并可以忽略后来的版本中新增的功能，但在HTTP/2中，就不能兼容了，HTTP/2使用了不同的数据结构和格式。
出于这个原因，HTTP/2被视为主版本更新。

HTTP/2.0还是HTTP/2

- HTTP/2定义了新版本HTTP的主要部分（二进制、多路复用等），并且未来的任何实现或变更（如果有HTTP/2.1的话），此规范都兼容。

> 所以小版本号不太重要.

- 此外，与HTTP/1消息不同，在HTTP/2请求中未明确声明版本号。
  - 例如，HTTP/2中没有GET /index.html HTTP/1.1形式的请求。
  - 但是，许多实现会在日志文件中使用次要版本号（.0）。例如，在Apache日志文件中，版本号显示为HTTP/2.0，其甚至会伪造HTTP/1形式的请求

![](https://res.weread.qq.com/wrepub/epub_32517945_78)

### 4.1.1 使用二进制格式替换文本格式

使用基于文本的协议，要先发完请求，并接收完响应之后，才能开始下一个请求。

HTTP的小改进

- HTTP/1.0引入了二进制的HTTP消息体，支持在响应中发送图片或其他媒体文件。
- HTTP/1.1引入了管道化和分块编码。
  - 都有队头阻塞（HOL）的问题——在队列首部的消息会阻塞后面消息的发送，更不用说，管道化在实际应用中并没有得到很好的支持。

### 4.1.2 多路复用代替同步请求

HTTP/1是一种同步的、独占的请求-响应协议。HTTP/1的主要解决方法是打开多个连接，并且使用资源合并以减少请求数，但这两种解决方法都会引入其他的问题和
带来性能开销。

图4.1　并发多个HTTP/1请求，需要多个TCP连接

![](https://res.weread.qq.com/wrepub/epub_32517945_79)

HTTP/2允许在单个连接上同时执行多个请求，每个HTTP请求或响应使用不同的流。通过使用二进制分帧层，给每个帧分配一个流标识符，以支持同时发出多个独立请
求。当接收到该流的所有帧时，接收方可以将帧组合成完整消息。

图4.2　使用多路复用技术的HTTP/2连接请求三个资源

![](https://res.weread.qq.com/wrepub/epub_32517945_80)

- **HTTP/2连接在请求发出后不需要阻塞到响应返回**
- 服务器发送响应的顺序完全取决于服务器，但客户端可以指定优先级。
- 每个请求都有一个新的、自增的流ID（如图4.2中所示的流5、7和9）。返回响应时使用相同的流ID
- 响应完成后，流会被丢弃而且不能重用
- 为了防止流ID冲突，客户端发起的请求使用奇数流ID, 服务器发起的请求使用偶数流ID。
- 请注意，在写作本书时，从技术上讲**服务器不能新建一个流**，除非是特殊情况（服务端推送，但也要客户端先发起请求）
- ID为0的流（图中未显示出）是客户端和服务器用于管理连接的控制流。

小结

- HTTP/2使用多个二进制帧发送HTTP请求和响应，使用单个TCP连接，以流的方式多路复用。
- HTTP/2与HTTP/1的不同主要在消息发送的层面上，在更上层，HTTP的核心概念不变。例如，请求包含一个方法（例如GET）、想要获取的资源
  （例如/styles.css）、首部、正文、状态码（例如200、404）、缓存、Cookie等，这些都与HTTP/1保持一致。

图4.3　HTTP/2中的流和HTTP/1中的连接相似

![](https://res.weread.qq.com/wrepub/epub_32517945_81)

### 4.1.3 流的优先级和流量控制

现在HTTP/2对并发的请求数量的限制放宽了很多（在许多实现中，默认情况下允许同时存在100个活跃的流），因此许多请求不再需要浏览器来排队，可以立即发送它
们。这可能导致带宽浪费在较低优先级的资源（例如图像）上，从而导致在HTTP/2下页面的加载速度变慢。所以需要控制流的优先级，使用更高的优先级发送最关键
的资源。

流的优先级控制是通过这种方式实现的：当数据帧在排队时，服务器会给高优先级的请求发送更多的帧。

流量控制是在同一个连接上使用多个流的另一种方式。如果接收方处理消息的速度慢于发送方，就会存在积压，需要将数据放入缓冲区。而当缓冲区满时会导致丢包，
需要重新发送。在连接层，TCP支持限流，但HTTP/2要在流的层面实现流量控制。

### 4.1.4 首部压缩

HTTP首部（包括请求首部和响应首部）用于发送与请求和响应相关的额外信息。在这些首部中，有很多信息是重复的，多个资源使用的首部经常相同。

- Cookie
- User-Agent
- Host
- Accept
- Accept-Encoding

### 4.1.5 服务端推送

HTTP/2添加了服务端推送的概念，它**允许服务端给一个请求返回多个响应**。

HTTP/2服务端推送是HTTP协议中的新概念，如果使用不当，它很容易浪费带宽。浏览器并不需要推送的资源，特别是，在之前已经请求过的服务器推送的资源，放在
浏览器缓存中。决定什么时候推送、如何推送，是充分利用服务端推送的关键。

## 4.2 如何创建一个HTTP/2连接

- 使用HTTPS协商。
- 使用HTTP Upgrade首部。
- 和之前的连接保持一致。

理论上，HTTP/2支持基于未加密的HTTP（也就是h2c）创建连接，也支持基于加密的HTTPS（即h2）创建连接。实际上，所有的Web浏览器仅支持基于HTTPS（h2）
建立HTTP/2连接，所以浏览器使用第一个方法来协商HTTP/2。服务器之间的HTTP/2连接可以基于未加密的HTTP（h2c）或者HTTPS（h2）。

### 4.2.1 使用HTTPS协商

HTTPS需要经过一个协议协商阶段来建立连接，在建立连接并交换HTTP消息之前，它们需要协商SSL/TLS协议、加密的密码，以及其他的设置。这个过程比较灵活，
可以引入新的HTTPS协议和密码，只要客户端和服务端都支持就行。在HTTPS握手的过程中，可以同时完成HTTP/2协商，这就不需要在建立连接时增加一次跳转。

HTTPS握手

公钥私钥加密被称为非对称加密, 但它比较慢，所以这种加密方式用于协商一个对称加密的密钥，以便在创建连接之后使用对称密钥加密消息。

图4.4　HTTPS握手 (TLSv1.2, 与TLSv1.3略微不同)

![](https://res.weread.qq.com/wrepub/epub_32517945_82)

握手过程涉及4类消息：

1. 客户端发送一个ClientHello消息，用于详细说明自己的加密能力。不加密此消息，因为加密方法还没有达成一致。
2. 服务器返回一个SeverHello消息，用于选择客户端所支持的HTTPS协议（如TLSv1.2）。基于客户端在ClientHello中声明的密码，和服务器本身支持的密码，
   服务器返回此连接的加密密码（如ECDHE-RSA-AES128-GCM-SHA256）。
   - 之后提供服务端HTTPS证书（ServerCertificate）。
   - 然后是基于所选密码加密的密钥信息（ServerKeyExchange）
   - 以及是否需要客户端发送客户端证书（CertificateRequest，大多数网站不需要）的说明。
   - 最后，服务端宣告本步骤结束（ServerHelloDone）。
3. 客户端校验服务端证书，如果需要发送客户端证书（ClientCertificate，大多数网站不需要）。
   - 然后发送密钥信息（ClientKeyExchange）。这些信息通过服务端证书中的公钥加密，所以只有服务端可以通过密钥解密消息。
   - 如果使用客户端证书，则会发送一个**CertificateVerify消息，此消息使用私钥签名**，以证明客户端对证书的拥有权。
   - 客户端使用ServerKeyExchange和ClientKeyExchange信息来定义一个加密过的对称加密密钥，然后发送一个ChangeCipherSpec消息通知服务端加密开
     始
   - 最后发送一个Finished消息。

> 配置客户端证书需要私钥吗?

4. 服务端也切换到加密连接上（ChangeCipherSpec），然后发送一个加密过的Finished消息。

当HTTPS会话建立完成后，在同一个连接上的HTTP消息就不再需要这个协商过程了。类似地，后续的连接（不管是并发的额外连接，还是后来重新打开的连接）可以
跳过其中的某些步骤 —— 如果它复用上次的加密密钥，这个过程就叫作**TLS会话恢复**。

TLSv1.3可以将协商过程中的消息往返减少到1个（如果复用之前的协商结果，则可以降到0个）

**ALPN**

ALPN[5]给ClientHello和ServerHello消息添加了功能扩展，客户端可以用它来声明应用层协议支持（“嗨，我支持h2和http/1，你用哪个都行。”），服务端
可以用它来确认在HTTPS协商之后所使用的应用层协议（“好的，我们用h2吧”）。

图4.5　使用ALPN的HTTPS握手

![](https://res.weread.qq.com/wrepub/epub_32517945_83)

NPN

在使用NPN时，客户端决定最终使用的协议，而在使用ALPN时，服务端决定最终使用的协议。

图4.6　使用NPN的HTTPS握手

![](https://res.weread.qq.com/wrepub/epub_32517945_84)

现在不再推荐使用NPN，应该使用ALPN。

使用ALPN进行HTTPS握手的示例(curl)

![](https://res.weread.qq.com/wrepub/epub_32517945_85)

### 4.2.2 使用HTTP Upgrade首部

通过发送Upgrade首部，客户端可以请求将现有的HTTP/1.1连接升级为HTTP/2。这个首部应该只用于未加密的HTTP连接（h2c）。基于加密的HTTPS连接的
HTTP/2（h2）不得使用此方法进行HTTP/2协商，它必须使用ALPN。我们已经说过多次，Web浏览器只支持基于加密连接的HTTP/2，所以它们不会使用这个方法。

示例1：一个不成功的Upgrade请求

客户端支持并想要使用HTTP/2，发送一个带Upgrade首部的请求：

![](https://res.weread.qq.com/wrepub/epub_32517945_87)

这样的请求必须包含一个HTTP2-Settings首部，它是一个Base-64编码的HTTP/2 SETTINGS帧

不支持HTTP/2的服务器可以像之前一样返回一个HTTP/1.1消息，就像Upgrade首部没有发送一样：

![](https://res.weread.qq.com/wrepub/epub_32517945_88)

示例2：一个成功的Upgrade请求

支持HTTP/2的服务器可以返回一个HTTP/1.1 101响应以表明它将切换协议，而不是忽略升级请求，并返回HTTP/1.1 200响应：

![](https://res.weread.qq.com/wrepub/epub_32517945_89)

然后服务器直接切换到HTTP/2，发送SETTINGS帧（见4.3.3节），之后以HTTP/2格式发送响应。

示例3：服务端请求的升级

当客户端认为服务器不支持HTTP/2时，它会发送不带Upgrade的请求：

![](https://res.weread.qq.com/wrepub/epub_32517945_90)

一个支持HTTP/2的服务端可以返回一个200响应，但是在响应首部中添加Upgrade来说明自己支持HTTP/2。这个时候，它是一个升级建议，而不是升级请求，因为只
有客户端才发起升级请求。

如下是一个服务端宣告支持h2（基于HTTPS的HTTP/2）和h2c（基于纯文本的HTTP/2）的示例：

![](https://res.weread.qq.com/wrepub/epub_32517945_91)

客户端可以利用这个信息来完成协议升级，并在下一个请求中发送一个Upgrade首部

![](https://res.weread.qq.com/wrepub/epub_32517945_92)

发送Upgrade首部的问题

由于所有的浏览器都只支持基于HTTPS的HTTP/2，因此这个Upgrade方法可能永远不会被浏览器使用，这会带来问题。

- 应用服务器可能会发送一个Upgrade首部，帮助升级到HTTP/2以提升性能。反向代理Web服务器可能会透传这个首部。浏览器会收到升级建议，并决定升级。但是
  与客户端直接连接的这个反向代理Web服务器并不支持HTTP/2。
- 在类似的场景中，可能反向代理已经和Web浏览器使用HTTP/2交互，但使用HTTP/1.1将请求代理到后端应用服务器。应用服务器可能会发出升级建议，如果其被反
  向代理透传，浏览器就会困惑，因为当前已经使用HTTP/2通信了，服务端还在建议升级到HTTP/2。

### 4.2.3 使用先验知识

有不同的方法可以让客户端事先知道服务器是否支持HTTP/2。如果你使用反向代理来卸载HTTPS，则可能会通过基于纯文本的HTTP/2（h2c）与后端服务器通信，因
为你知道它们支持HTTP/2。或者，可以根据Alt-Svc首部（HTTP/1.1）或ALTSVC帧（参见4.3.4节）推断先前的连接信息。

### 4.2.4 HTTP Alternative Services

第4种方法是使用HTTP Alternative Services（替代服务），它没有被包含在原来的标准中，在HTTP/2发布之后，将其列为单独的标准。此标准允许服务器使用
HTTP/1.1协议（通过Alt-Svc HTTP首部）通知客户端，它所请求的资源在另一个位置（例如，另一个IP或端口），可以使用不同的协议访问它们。该协议可以使用
先验知识启用HTTP/2。

### 4.2.5 HTTP/2前奏消息

不管使用哪种方法启用HTTP/2连接，在HTTP/2连接上发送的第一个消息必须是HTTP/2连接前奏，或者说是“魔法”字符串。此消息是客户端在HTTP/2连接上发送的
第一个消息。它是一个24个八位字节的序列，以十六进制表示法显示如下：

这个序列被转换为ASCII字符串后如下所示：

![](https://res.weread.qq.com/wrepub/epub_32517945_95)
![](https://res.weread.qq.com/wrepub/epub_32517945_96)

这个无意义的看起来像HTTP/1样式的消息，目的是兼容，客户端向不支持HTTP/2的服务端发送HTTP/2消息的情况。然后服务器会尝试解析此消息，就像收到其他
HTTP消息时一样。因为它无法识别这个无意义的方法（PRI）和HTTP版本（HTTP/2.0），所以解析会失败，从而拒绝此消息。注意，此消息前奏是官方规范中唯一一
处使用HTTP/2.0的地方，在其他地方都是HTTP/2，正如4.1节中所讨论的。而对于支持HTTP/2的服务器，可以根据这个收到的前奏消息推断出客户端支持HTTP/2，
它不会拒绝这个神奇的消息，它必须发送SETTINGS帧作为其第一条消息（可以为空）。

为什么是PRI和SM

在早期的草稿中，HTTP/2规范中的消息前奏使用FOO和BAR或者BA表示，它们是编程中常见的占位符。但是在规范草稿的第4个版本中，这个占位符变成了PRI SM，
但是没有说为什么。

## 4.3 HTTP/2帧

### 4.3.1 查看HTTP/2帧

使用Chrome net-export

- 抓包[chrome://net-export](https://netlog-viewer.appspot.com/#import)
- 查看日志 https://netlog-viewer.appspot.com/#import

图4.7　在Chrome中查看HTTP/2帧

![](https://res.weread.qq.com/wrepub/epub_32517945_98)

以下输出来自一个SETTINGS帧：

![](https://res.weread.qq.com/wrepub/epub_32517945_99)

使用[nghttp](https://nghttp2.org/)

![](https://res.weread.qq.com/wrepub/epub_32517945_100)

使用Wireshark

- 需要告诉Wireshark HTTPS密钥

![](https://res.weread.qq.com/wrepub/epub_32517945_101)

- 对于macOS，设置SSLKEYLOGFILE环境变量

![](https://res.weread.qq.com/wrepub/epub_32517945_102)

- 或者直接作为命令行参数提供

![](https://res.weread.qq.com/wrepub/epub_32517945_103)

- 在Wireshark中加载密钥

图4.8　设置Wireshark HTTPS密钥文件

![](https://res.weread.qq.com/wrepub/epub_32517945_104)

图4.9　Wireshark中显示的HTTP/2魔法字符串

![](https://res.weread.qq.com/wrepub/epub_32517945_105)

图4.10　Wireshark中的ClientHello消息中的ALPN扩展

![](https://res.weread.qq.com/wrepub/epub_32517945_107)

### 4.3.2 HTTP/2帧数据格式

每个HTTP/2帧由一个固定长度的头部和不定长度的负载组成。

表4.1　HTTP/2帧头部格式

![](https://res.weread.qq.com/wrepub/epub_32517945_108)

- Stream Identifier, 将此字段限制为31位的其中一个原因是考虑到Java的兼容性，因为它没有32位无符号整数

### 4.3.3 HTTP/2消息流示例

![](https://res.weread.qq.com/wrepub/epub_32517945_109)
![](https://res.weread.qq.com/wrepub/epub_32517945_110)
![](https://res.weread.qq.com/wrepub/epub_32517945_111)

可以使用-n参数来隐藏数据，仅显示帧头部：

![](https://res.weread.qq.com/wrepub/epub_32517945_112)

1. 首先，通过HTTPS（h2）协商建立HTTP/2连接。nghttp不输出HTTPS建立过程和HTTP/2前奏/“魔术”消息，因此我们首先看到SETTINGS帧：

![](https://res.weread.qq.com/wrepub/epub_32517945_113)

**SETTINGS帧**

SETTINGS帧（0x4）是服务器和客户端必须发送的第一个帧（在HTTP/2前奏/“魔术”消息之后）。该帧不包含数据，或只包含若干键/值对

表4.2　SETTINGS帧格式

![](https://res.weread.qq.com/wrepub/epub_32517945_114)

> 规范中的Value是默认值, 顺序与Identifier对应, e.g. 0x1->4096.

再回头看第一个消息

![](https://res.weread.qq.com/wrepub/epub_32517945_115)

收到的SETTINGS帧有30个8位字节数据，没有设置标志位（因此不是确认帧），使用的流ID为0。流ID 0是保留数字，用于控制消息（SETTINGS和
WINDOW_UPDATE帧），所以服务器使用流ID 0发送此SETTINGS帧是合理的。

在此示例中有5个设置项（niv=5），每个设置项长度为16位（标识符）+32位（值）。也就是说，每项设置有48位即6字节

2. 查看接下来的3个SETTINGS帧

![](https://res.weread.qq.com/wrepub/epub_32517945_117)

nghttp接收初始服务器SETTINGS帧（刚讨论过），然后，客户端发送带有几个设置项的SETTINGS帧。接下来，客户端确认服务器的SETTINGS帧。确认SETTINGS
帧非常简单，只有一个ACK（0x01）标志，长度为0，因此只有0设置（niv=0）。再接下来是服务器确认客户端的SETTINGS帧，格式同样简单。

**WINDOW_UPDATE帧**

![](https://res.weread.qq.com/wrepub/epub_32517945_118)

WINDOW_UPDATE帧（0x8）用于流量控制，比如限制发送数据的数量，防止接收端处理不完。在HTTP/1下，同时只能有一个请求。如果客户端无法及时处理数据，它
会停止处理TCP数据包，然后TCP流量控制（类似HTTP/2流量控制）开始工作，降低发送数据的流量，直到接收方可以正常处理为止。在HTTP/2下，在同一个连接上
有多个流，所以不能依赖TCP流量控制，必须自己实现针对每个流的减速方法。

表4.3　WINDOW_UPDATE帧格式

![](https://res.weread.qq.com/wrepub/epub_32517945_119)

WINDOW_UPDATE帧未定义标志位，该设置用于给定的流，如果流ID指定为0，则应用于整个HTTP/2连接。发送方必须跟踪每个流和整个连接。

> `如果流ID指定为0`中的流ID应该是帧首部中的`Stream Identifier`.

HTTP/2流量控制设置仅应用于DATA帧，所有其他类型的帧（至少目前定义的），就算超出了窗口大小的限制也可以继续发送。这个特性可以防止重要的控制消息（比
如WINDOW_UPDATE帧自己）被较大的DATA帧阻塞。同时DATA帧也是唯一可以为任意大小的帧。

**PRIORITY帧**

![](https://res.weread.qq.com/wrepub/epub_32517945_120)

通过dep_stream_id，它将其他流悬挂在开始时创建的流之下(该流依赖dep_stream_id所指定的流)。使用之前创建的流的优先级，可以方便地对请求进行优先级
排序，无须为每个后续新创建的流明确指定优先级。并非所有HTTP/2客户端都给流**预定义优先级**

表4.4　PRIORITY帧格式

![](https://res.weread.qq.com/wrepub/epub_32517945_121)

PRIORITY帧（0x2）长度固定，没有定义标志位。

**HEADERS帧**

**一个HTTP/2请求以HEADERS帧开始发送（0x1）**

![](https://res.weread.qq.com/wrepub/epub_32517945_122)

HTTP/2定义了新的伪首部（以冒号开始），以定义HTTP请求中的不同部分：

![](https://res.weread.qq.com/wrepub/epub_32517945_124)

- :authority伪首部代替了原来HTTP/1.1的Host首部
- HTTP/2伪首部定义严格，不像标准的HTTP首部那样可以在其中添加新的自定义首部

不能这样创建新的伪首部

```
:barry: value
```

如果应用需要，还得用普通的HTTP首部，没有开头的冒号

```
barry: value
```

可以依照新的规范来创建新的伪首部

- 在Bootstrapping WebSockets with HTTP/2 RFC中添加:protocol伪首部。应用新的伪首部需要使用新的SETTINGS参数，也需要客户端和服务端的支持。
- 可以在客户端工具中查看这些伪首部，它们表明正在使用HTTP/2请求

图4.11　Chrome开发者工具中的伪首部

![](https://res.weread.qq.com/wrepub/epub_32517945_127)

**HTTP/2强制将HTTP首部名称小写**, HTTP首部的值可以包含不同的大小写字母

- HTTP/2对HTTP首部的格式要求也更严格。开头的空格、双冒号或者换行，在HTTP/2中都会带来问题
- 当客户端发现首部格式不正确时，报错信息通常含义不明（比如Chrome中的ERR_SPDY_PROTOCOL_ERROR）

表4.5　HEADERS帧格式

![](https://res.weread.qq.com/wrepub/epub_32517945_128)

- 添加Pad Length和Padding字段是出于安全原因，用以隐藏真实的消息长度。
- Header Block Fragment（首部块片段）字段包含所有的首部（和伪首部）。这个字段不是纯文本的，不像nghttp里所显示的那样。(首部压缩了)

HEADERS首部定义了4个标志位

- END_STREAM(0x1)，如果当前HEADERS帧后面没有其他帧（比如POST请求，后面会跟DATA帧），设置此标志。有点违反直觉的是，CONTINUATION帧不受此限
  制，它们由END_HEADERS标志控制，被当作HEADERS帧的延续，而不是额外的帧。
- END_HEADERS（0x4），它表明所有的HTTP首部都已经包含在此帧中，后面没有CONTINUATION帧了。
- PADDED(0x8)，当使用数据填充时设置此标志位。这个标志表明，DATA帧的前8位代表HEADERS帧中填充的内容长度。
- PRIORITY(0x20)，表明在帧中设置了E、Stream Dependency和Weight字段。

如果HTTP首部尺寸超出一个帧的容量，则需要使用一个CONTINUATION帧（紧接着是一个HEADERS帧），而不是使用另外一个HEADERS帧。

- 这个过程相较于HTTP正文来说好像过于复杂，HTTP正文会使用多个DATA帧。因为表4.5中的其他字段只能使用一次，所以如果同一个请求有多个HEADERS帧，并
  且它们的其他字段值不同，就会带来一些问题。
- 要求**CONTINUATION帧紧跟在HEADERS帧后面，其中不能插入其他帧**，这影响了HTTP/2的多路复用，人们正考虑其他替代方案。
- 实际上CONTINUATION帧很少使用，大多数请求都不会超出一个HEADERS帧的容量。

再回头看这些日志输出

![](https://res.weread.qq.com/wrepub/epub_32517945_129)

- 每个新的请求都会被分配一个独立的流ID，其值在上一个流ID的基础上自增（在这个示例中上一个流ID是11，它是nghttp创建的PRIORITY帧，所以这个帧使用
  流ID13创建，偶数12是服务端使用的）。
- 同时设置了多个标志位，组合起来的十六进制数为0x25
  - 其中的END_STREAM（0x1）和END_HEADERS（0x4）标志位说明，当前帧包含完整的请求，没有DATA帧（可能用于POST请求）。
  - PRIORITY标志位（0x20）表明，此帧使用了优先级策略。
  - 将这些十六进制数加起来（0x1 + 0x4 + 0x20），结果是0x25，在帧首部中显示。
- 这个流依赖流11，所以被分配了对应的优先级，权重为16。
- nghttp的注释说，这个流是新建的（Open new stream），
- 然后列出了多个HTTP伪首部和HTTP请求首部

