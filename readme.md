# mongo-operator

to manage a replica set of mongo in k8s.

# build  && run

```
kubectl create -f deploy/crds/cache_v1alpha1_memcached_crd.yaml
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
