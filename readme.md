# mongo-operator

to manage a replica set of mongo in k8s.

# build  && run

```
kubectl create -f deploy/crds/cache_v1alpha1_memcached_crd.yaml
```



# dev

1. 如果更改了 pkg/apis/*_types.go 相关 scheme 记得:

```bash
operator-sdk generate k8s --verbose

```

**国内运行 operator-sdk 需要代理以安装依赖。**

2. 更新代码逻辑后, image, 并更新 operator-部署 manifast 文件 image 相关内容。
```bash
operator-sdk build  dockerhub.smartx.com/mongo-operator:v0.02    
docker push dockerhub.smartx.com/mongo-operator:v0.02 
```



2. 快速验证

operator-sdk 提供了 k8s 外运行operator 的机制

需要开发者配置好相应的 ~/.kube/config 到开发机。 或者也可以主动

```bash
export OPERATOR_NAME=mongo-operator
operator-sdk up local --namespace=default

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
