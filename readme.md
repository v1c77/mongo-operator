
初始化
```
export workdir=$GOPATH/src/github.smartx.com/mongo-operator
mkdir -p $workpath
cd $workpath
export GO111MODULE=on
operator-sdk new mongo-operator
cd mongo-operator
```

资源初始化
```
operator-sdk add api --api-version=app.smartx.com/v1alpha1 --kind=mongo

operator-sdk add  controller --api-version=app.smartx.com/v1alpha1 --kind=mongo

```
