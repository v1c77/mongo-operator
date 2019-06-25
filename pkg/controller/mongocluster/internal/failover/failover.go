package failover

import (
	dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	k8s "github.smartx.com/mongo-operator/pkg/service/kubernetes"
	"github.smartx.com/mongo-operator/pkg/service/mongo"
	"github.smartx.com/mongo-operator/pkg/utils"
	"strings"
)

var logger = utils.NewLogger("mongocluster.failover")

type MongoClusterFailover struct {
	checker *MongoClusterFailoverChecker
	//healer *MongoClusterFailoverHealer
}

func NewMongoClusterFailover(k8sService k8s.Services) *MongoClusterFailover {

	checker := NewMongoClusterFailoverChecker(k8sService)
	//healer := NewMongoClusterFailoverHealer(k8sService, mongoClient)
	return &MongoClusterFailover{
		checker: checker,
		//healer: healer,
	}
}

//  CheckAndHeal
func (f *MongoClusterFailover) CheckAndHeal(mc *dbv1alpha1.
	MongoCluster) error {

	//fLogger := logger.L().WithValues("Request.Namespace", mc.Namespace,
	//	"Request.Name", mc.Name)
	// TODO(yuhua) some pre-check
	// TODO(yuhua) check pods status. pod 需要全部启动并且获取到 ip.
	// pod 数量需要与 spec.Replicas 一致。

	// ================== mongo status check/

	// TODO mongo replica
	// 1。 not -init , 创建之，

	// 2。 other(pod  断电等意外重启，集群健康时， 进到相应节点 reconfig 之)
	err := f.checkAndHealMongoReplicaSet(mc)
	if err != nil {
		return err
	}

	// TODO 其他状态待定。

	// - check pod, service, network, mongo replicaset..

	return nil
}

// checkAndHealMongoReplicaSet
func (f *MongoClusterFailover) checkAndHealMongoReplicaSet(
	mc *dbv1alpha1.MongoCluster) error {

	// 获取 mongoCluster pod service
	dnsList := f.checker.GetMembersDNS(mc)

	url := dnsList[0]
	logger.Info("Check pod Dns list", "url", url)
	mongoClient := mongo.NewClient(url)

	mgoSession, err := mongoClient.DialDirect()
	if err != nil {
		return err
	}

	rStatus, err := f.checker.CheckReplicaSetStatus(mgoSession)
	if err != nil {
		if strings.Contains(err.Error(),
			"no replset config has been received") {
			// do init
			logger.Info("do mongo cluster initial.")
			// make
			return nil
		} else {
			return err
		}
	}

	logger.Info("get status", "status", rStatus)
	return nil
}
