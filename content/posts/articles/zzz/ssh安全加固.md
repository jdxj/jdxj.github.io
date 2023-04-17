---
title: "SSH 安全加固的一些措施"
date: 2023-04-17T14:39:37+08:00
tags:
  - ssh
---

[原文](https://zzz.buzz/zh/2016/04/18/hardening-sshd/)

可用的简单配置

```
# 禁用旧版本的有漏洞的 SSH 协议
# (旧版本的 sshd 默认允许新旧版本(2,1)的协议，新版本的 sshd 默认只接受新版本(2)的协议）
Protocol 2

# 禁止使用密码登入（默认允许）
# 注意：需先配置好有 root 权限的用户的公私钥对登入后再禁止密码登入，
#       避免自己无法登入服务器，或使用低权限用户登入服务器后无法取得 root 权限。
PasswordAuthentication no

# 禁止空密码账户登入（默认禁止）
PermitEmptyPasswords no

# 禁止 root 账户通过密码登入（默认允许）
# （root 账户仍可以通过公私钥对登入，如果配置了的话）
PermitRootLogin without-password
# （新版本的 sshd 也支持使用更符合直觉的名字 prohibit-password）
#PermitRootLogin prohibit-password
# 或：禁止 root 账户通过 SSH 登入（默认允许）
#PermitRootLogin no
```
