package failover

import (
	dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	k8s "github.smartx.com/mongo-operator/pkg/service/kubernetes"
	"github.smartx.com/mongo-operator/pkg/service/mongo"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("mongocluster.failover")

type MongoClusterFailover struct {
	checker *MongoClusterFailoverChecker
	//healer *MongoClusterFailoverHealer
}

func NewMongoClusterFailover(k8sService k8s.Services,
	mongoClient mongo.Client) * MongoClusterFailover {

	checker := NewMongoClusterFailoverChecker(k8sService, mongoClient)
	//healer := NewMongoClusterFailoverHealer(k8sService, mongoClient)
		return &MongoClusterFailover{
			checker: checker,
			//healer: healer,
		}
}


//  CheckAndHeal
func (f *MongoClusterFailover) CheckAndHeal(mc *dbv1alpha1.
	MongoCluster) error {

		fLogger := log.WithValues("Request.Namespace", mc.Namespace,
			"Request.Name", mc.Name)
		// TODO(yuhua) some pre-check

		// 获取 mongoCluster pod service
		dnsList := f.checker.GetMembersDNS(mc)
		fLogger.Info("Check pod Dns list", "dns", dnsList)
		// TODO 获取  mongo 串.

		err := f.checker.MemebersStatus(mc)
		if err!= nil {
			fLogger.Error(err, "Status not ok.")
		}
		//statusResp := f.checker.GetReplicaSetStatus()
		//fLogger.Info("GetReplicaSet status", "status", StatusResp.msg)
		// - check pod, service, network, mongo replicaset..

		return nil
}