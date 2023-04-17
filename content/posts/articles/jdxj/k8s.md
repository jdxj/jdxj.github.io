---
title: "安装K8s"
date: 2023-04-17T13:38:15+08:00
tags:
  - kubernetes
---

# 安装Docker

- [Install Docker Engine on Debian](https://docs.docker.com/engine/install/debian/)

配置containerd

```shell
$ containerd config default | tee /etc/containerd/config.toml >/dev/null 2>&1
$ vim /etc/containerd/config.toml
SystemdCgroup = true
$ systemctl restart containerd
```

# 安装K8s

关闭交换

```shell
$ swapoff -a
$ vim /etc/fstab
```

官方安装教程

1. [容器运行时](https://kubernetes.io/zh-cn/docs/setup/production-environment/container-runtimes/)
2. [安装 kubeadm](https://kubernetes.io/zh-cn/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)
3. [使用 kubeadm 创建集群](https://kubernetes.io/zh-cn/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)

网络插件weave

- [Installation](https://www.weave.works/docs/net/latest/kubernetes/kube-addon/#install)

允许控制平面调度pod

```shell
$ kubectl taint nodes --all node-role.kubernetes.io/control-plane-
```

bash补全

- [启动 kubectl 自动补全功能](https://kubernetes.io/zh-cn/docs/tasks/tools/install-kubectl-linux/#enable-kubectl-autocompletion)

# 参考

- [如何用 Kubeadm 在 Debian 11 上安装 Kubernetes 集群](https://www.51cto.com/article/740996.html)
- [master节点部署Pod处于Pending状态](https://cloud.tencent.com/developer/article/1992628)
