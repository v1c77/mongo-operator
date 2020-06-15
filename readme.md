# mongo-operator

to manage a replica set of mongo in k8s.

# requirement

-  a pvc provisioner: like a [local storage provisioner.](https://github.com/kubernetes-sigs/sig-storage-local-static-provisioner)
- dev requirement: operator-sdk.


# build  && run

```bash
kubectl create -f deploy/crds/db_v1alpha1_mongocluster_crd.yaml

kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/operator.yaml
```

# dev.cn

> operator-sdk 版本 0.8.1

1. 如果更改了 pkg/apis/*_types.go 相关 scheme，需要更新k8s api:

```bash
make gen-k8s
```

> **国内运行 operator-sdk 需要代理以安装依赖。**

2. 更新 operator 代码逻辑后，需要重新打包上传.
```bash
make build # or `make build VERSION=0.0.4`  
```

2. 快速验证

operator-sdk 提供了 k8s 外运行 operator 的机制

条件：
-  `~/.kube/config`  存在并正确配置
- 主机网络可感知 coreos, 正确配置 `/etc/resolv.conf` 

```bash
make run
```

--------------------

备忘：

```
export workdir=$GOPATH/src/github.smartx.com/mongo-operator
mkdir -p $workpath
cd $workpath
export GO111MODULE=on
operator-sdk new mongo-operator
cd mongo-operator
```
```
operator-sdk add api --api-version=db.smartx.com/v1alpha1 --kind=MongoCluster
operator-sdk generate k8s --verbose  
operator-sdk add controller --api-version=db.smartx.com/v1alpha1 --kind=MongoCluster 
go mod vendor

```
