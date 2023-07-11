---
title: "用Go转换字符编码"
date: 2023-07-11T11:02:46+08:00
draft: true
tags:
  - go
  - unicode
  - utf
---

计算机字符集中的每个字符都有两个属性：码点（code point）和表示这个码点的内存编码（位模式，表示这个字符码点的二进制比特串）。

- 所谓码点（这里借用了Unicode字符集中码点的概念）是指将字符集中所有字符“排成一队”，字符在队伍中的唯一序号值。

表52-1　ASCII字符集中字符的码点与内存编码

- ASCII字符集中每个字符的码点与其内存编码表示是一致的。

![](https://res.weread.qq.com/wrepub/epub_42557147_42)

图52-1　同一内存编码表示对应不同字符集中的不同字符

![](https://res.weread.qq.com/wrepub/epub_42557147_43)

图52-2　Unicode字符集码点表

![](https://res.weread.qq.com/wrepub/epub_42557147_44)

Unicode的前128个码点与ASCII字符码点是一一对应的

Unicode采用的内存编码表示方案

- utf-16

该方案使用2或4字节表示每个Unicode字符码点。它的优点是编解码简单，因为所有字符都用偶数字节表示；其不足也很明显，比如存在字节序问题、不兼容ASCII
字符内存表示以及空间效率不高等。

- utf-32

该方案固定使用4字节表示每个Unicode字符码点。它的优点也是编解码简单，因为所有字符都用4字节表示；其不足也和UTF-16一样明显，同样存在字节序问题、不
兼容ASCII字符内存表示，且空间效率是这三种方案中最差的。

- utf-8

UTF-8编码使用的字节从1到4不等。前128个与ASCII字符重合的码点（U+0000~U+007F）使用1字节表示；带变音符号的拉丁文、希腊文、西里尔字母、阿拉伯文
等使用2字节来表示；而东亚文字（包括汉字）使用3字节表示；其他极少使用的语言的字符则使用4字节表示。

Unicode规范中对字节序标记的约定如下

- 如果没有提供字节序标记，则默认采用大端字节序解码。
- 由于UTF-8没有字节序问题，因此这个BOM只是用于表明该数据流采用的是UTF-8编码方案

```
FF FE         UTF-16 小端字节序
FE FF         UTF-16 大端字节序
FF FE 00 00   UTF-32 小端字节序
00 00 FE FF   UTF-32 大端字节序
EF BB BF      UTF-8
```

在Go中，每个rune对应一个Unicode字符的码点，而Unicode字符在内存中的编码表示则放在[]byte类型中。从rune类型转换为[]byte类型，称为“编码”
（encode），而反过来则称为“解码”（decode）

图52-4　Go语言中Unicode字符的编解码

![](https://res.weread.qq.com/wrepub/epub_42557147_46)

```go
// chapter9/sources/go-character-set-encoding/rune_encode_and_decode.go

// rune -> []byte
func encodeRune() {
    var r rune = 0x4E2D // 0x4E2D为Unicode字符"中"的码点
    buf := make([]byte, 3)
    n := utf8.EncodeRune(buf, r)

    fmt.Printf("the byte slice after encoding rune 0x4E2D is ")
    fmt.Printf("[ ")
    for i := 0; i < n; i++ {
        fmt.Printf("0x%X ", buf[i])
    }
    fmt.Printf("]\n")
    fmt.Printf("the unicode charactor is %s\n", string(buf))
}

// []byte -> rune
func decodeRune() {
    var buf = []byte{0xE4, 0xB8, 0xAD}
    r, _ := utf8.DecodeRune(buf)
    fmt.Printf("the rune after decoding [0xE4, 0xB8, 0xAD] is 0x%X\n", r)
}

func main() {
    encodeRune()
    decodeRune()
}
```

图52-5　将“中国人”三个字符从UTF-8编码表示转换为GB18030编码表示

![](https://res.weread.qq.com/wrepub/epub_42557147_47)

```go
// chapter9/sources/go-character-set-encoding/convert_utf8_to_gb18030.go
package main

import (
    "bytes"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "unicode/utf8"

    "golang.org/x/text/encoding/simplifiedchinese"
    "golang.org/x/text/transform"
)

func dumpToFile(in []byte, filename string) error {
    f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
    if err != nil {
        return err
    }
    defer f.Close()
    _, err = f.Write(in)
    if err != nil {
        return err
    }
    return nil
}

func utf8ToGB18030(in []byte) ([]byte, error) {
    if !utf8.Valid(in) {
        return nil, errors.New("invalid utf-8 runes")
    }

    r := bytes.NewReader(in)
	// utf8 -> gb18030
    t := transform.NewReader(r, simplifiedchinese.GB18030.NewEncoder())
    out, err := ioutil.ReadAll(t)
    if err != nil {
        return nil, err
    }
    return out, nil
}

func main() {
    var src = "中国人" // <=> "\u4E2D\u56FD\u4EBA"
    var dst []byte

    for i, v := range src {
        fmt.Printf("Unicode字符: %s <=> 码点(rune): %X <=> UTF8编码内存表示: ", string(v), v)
        s := src[i : i+3]
        for _, v := range []byte(s) {
            fmt.Printf("0x%X ", v)
        }

        t, _ := utf8ToGB18030([]byte(s))
        fmt.Printf("<=> GB18030编码内存表示: ")
        for _, v := range t {
            fmt.Printf("0x%X ", v)
        }
        fmt.Printf("\n")

        dst = append(dst, t...)
    }

    dumpToFile(dst, "gb18030.txt")
}
```

真正执行UTF-8到GB18030编码形式转换的是simplifiedchinese.GB18030.NewEncoder方法，它读取以UTF-8编码表示形式存在的字节流（[]byte），并将
其转换为以GB18030编码表示形式的字节流返回。

将GB18030编码数据转换为UTF-16和UTF-32的示例

```go
// chapter9/sources/go-character-set-encoding/convert_gb18030_to_utf16_and_utf32.go

func catFile(filename string) ([]byte, error) {
    f, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    return ioutil.ReadAll(f)
}

func gb18030ToUtf16BE(in []byte) ([]byte, error) {
    r := bytes.NewReader(in) //gb18030
	// gb18030 -> utf8
    s := transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	// utf8 -> utf16
    d := transform.NewReader(s,
          unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder())

    out, err := ioutil.ReadAll(d)
    if err != nil {
        return nil, err
    }
    return out, nil
}

func gb18030ToUtf32BE(in []byte) ([]byte, error) {
    r := bytes.NewReader(in) //gb18030
    // gb18030 -> utf8
    s := transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	// utf8 -> utf32
    d := transform.NewReader(s,
          utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM).NewEncoder())

    out, err := ioutil.ReadAll(d)
    if err != nil {
        return nil, err
    }
    return out, nil
}

func main() {
    src, err := catFile("gb18030.txt")
    if err != nil {
        fmt.Println("open file error:", err)
        return
    }

    // 从gb18030到utf-16be
    dst, err := gb18030ToUtf16BE(src)
    if err != nil {
        fmt.Println("convert error:", err)
        return
    }

    fmt.Printf("UTF-16BE(no BOM)编码: ")
    for _, v := range dst {
        fmt.Printf("0x%X ", v)
    }
    fmt.Printf("\n")

    // 从gb18030到utf-32be
    dst1, err := gb18030ToUtf32BE(src)
    if err != nil {
        fmt.Println("convert error:", err)
        return
    }

    fmt.Printf("UTF-32BE(no BOM)编码: ")
    for _, v := range dst1 {
        fmt.Printf("0x%X ", v)
    }
    fmt.Printf("\n")
}
```
