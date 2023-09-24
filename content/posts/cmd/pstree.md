---
title: "pstree"
date: 2023-09-24T10:45:57+08:00
summary: 用树状形式显示所有进程之间的关系
tags:
  - optimize
---

```bash
# -a 表示输出命令行选项
# p表PID
# s表示指定进程的父进程
$ pstree -aps 3084
systemd,1
  └─dockerd,15006 -H fd://
      └─docker-containe,15024 --config /var/run/docker/containerd/containerd.toml
          └─docker-containe,3991 -namespace moby -workdir...
              └─app,4009
                  └─(app,3084)
```