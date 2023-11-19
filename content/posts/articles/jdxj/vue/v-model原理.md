---
title: "v-model原理"
date: 2023-10-28T11:51:48+08:00
tags:
  - vue
---

是一个语法糖

```html
<template>
    <div id="app">
        <input v-model="msg" type="text">
        <input :value="msg" @input="msg = $event.target.value" type="text">
    </div>
</template>
```

不同的标签有不同的事件

```html
<select :value="value" @change="selectCity">
```
