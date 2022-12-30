---
title: "Redis security"
date: 2022-12-30T15:54:38+08:00
draft: false
---

# Security model

ACLs

# Network security

- 避免向公网暴露

# Protected mode

# Authentication

- 推荐使用 ACLs
- 遗留的方式是使用 requirepass 配置
  - 一定要是很长的密码, 防止暴力破解
- redis 流量没有加密, auth 命令像其他命令一样, 可能被窃听

# TLS support

# Disallowing specific commands

- 推荐使用 ACLs

# Attacks triggered by malicious inputs from external clients

# String escaping and NoSQL injection

Redis 协议中没有字符串转义的概念.

# Code security

CONFIG 命令允许改变工作目录, 这个安全问题可能破坏系统, 或者运行不信任的代码
- 不要用 root 运行 redis