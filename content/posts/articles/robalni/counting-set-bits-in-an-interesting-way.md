---
title: "Counting Set Bits in an Interesting Way"
date: 2023-06-19T16:41:53+08:00
tags:
  - algorithm
---

[原文](https://www.robalni.org/posts/20220428-counting-set-bits-in-an-interesting-way.txt)

```c
unsigned popcnt(unsigned x) {
    unsigned diff = x;
    while (x) {
        x >>= 1;
        diff -= x;
    }
    return diff;
}
```
