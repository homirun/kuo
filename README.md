# kuo
[![CI](https://github.com/homirun/kuo/actions/workflows/test.yaml/badge.svg)](https://github.com/homirun/kuo/actions/workflows/test.yaml)

A kubernetes plugin that operates multiple contexts.

## Installation

```bash
$ make install
```

## Usage
### Setup .kuoconfig
Register contexts to be operated simultaneously. 

```bash
$ kubectl kuo set cluster1 cluster2

set .kuoconfig: [cluster1 cluster2]
```

### Execute kubectl in multiple contexts
```bash
$ kubectl kuo get node

======== cluster1 ========
NAME    STATUS   ROLES           AGE   VERSION
chino   Ready    control-plane   25d   v1.25.4+k0s
maya    Ready    <none>          25d   v1.25.4+k0s
megu    Ready    <none>          25d   v1.25.4+k0s
======== cluster2 ========
NAME       STATUS   ROLES           AGE    VERSION
minikube   Ready    control-plane   122m   v1.25.3
```

### Execute kubectl subcommand with configured contexts
Execute in the form `kubectl kuo [kubectl-subcommand] [flags]`.
```bash
$ kubectl kuo get node -o wide

======== cluster1 ========
NAME    STATUS   ROLES           AGE   VERSION       INTERNAL-IP      EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION      CONTAINER-RUNTIME
chino   Ready    control-plane   26d   v1.25.4+k0s   10.x.y.124       <none>        Ubuntu 22.04.1 LTS   xxxxxxxxxxxxxxxxx   containerd://x.y.z
maya    Ready    <none>          26d   v1.25.4+k0s   10.x.y.88        <none>        Ubuntu 22.04.1 LTS   xxxxxxxxxxxxxxxxx   containerd://x.y.z
megu    Ready    <none>          26d   v1.25.4+k0s   10.x.y.112       <none>        Ubuntu 22.04.1 LTS   xxxxxxxxxxxxxxxxx   containerd://x.y.z
======== cluster2 ========
NAME       STATUS   ROLES           AGE   VERSION   INTERNAL-IP    EXTERNAL-IP   OS-IMAGE               KERNEL-VERSION   CONTAINER-RUNTIME
minikube   Ready    control-plane   22h   v1.25.3   192.x.y.2      <none>        Buildroot 2021.02.12   x.y.z            docker://x.y.z

```
