# kuo
A kubernetes plugin that operates multiple contexts.

## Installation

```bash 
$ go build -o kubectl-kuo ./cmd/kuo
$  where kubectl                                                                                                      INT ✘  04:31:43 
/usr/local/bin/kubectl

$ mv kubectl-kuo /usr/local/bin/kubectl-kuo
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
