package failover

import (
	dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	k8s "github.smartx.com/mongo-operator/pkg/service/kubernetes"
	"github.smartx.com/mongo-operator/pkg/utils"
	"strings"
)

var logger = utils.NewLogger("mongocluster.failover")

type MongoClusterFailover struct {
	checker *MongoClusterFailoverChecker
	healer  *MongoClusterFailoverHealer
}

func NewMongoClusterFailover(k8sService k8s.Services) *MongoClusterFailover {

	checker := NewMongoClusterFailoverChecker(k8sService)
	healer := NewMongoClusterFailoverHealer(k8sService)
	return &MongoClusterFailover{
		checker: checker,
		healer:  healer,
	}
}

//  CheckAndHeal
func (f *MongoClusterFailover) CheckAndHeal(mc *dbv1alpha1.
	MongoCluster) error {

	//fLogger := logger.L().WithValues("Request.Namespace", mc.Namespace,
	//	"Request.Name", mc.Name)
	// TODO(yuhua) some pre-check
	// TODO(yuhua) check pods status. pod 需要全部启动并且获取到 ip.
	// TODO check if all the pods health. and then do the init.
	// pod 数量需要与 spec.Replicas 一致。

	// ================== mongo status check/

	// TODO mongo replica
	// 1。 not -init , 创建之，

	// 2。 other(pod  断电等意外重启，集群健康时， 进到相应节点 reconfig 之)
	err := f.checkAndHealMongoReplSet(mc)
	if err != nil {
		return err
	}

	// TODO 其他状态待定。

	// - check pod, service, network, mongo replicaset..

	return nil
}

// checkAndHealMongoReplicaSet
func (f *MongoClusterFailover) checkAndHealMongoReplSet(
	mc *dbv1alpha1.MongoCluster) error {

	podStatus := f.checker.GetMongoPodsStatus(mc)
	var healthNode []string
	var newNode []string
	var issueNode []string
	for url, podStatus := range podStatus {
		// init mongo replicaset.
		if podStatus.Err != nil {
			if strings.Contains(podStatus.Err.Error(),
				"no replset config has been received") {
				newNode = append(newNode, url)
			} else {
				// other issue like pod restart.
				issueNode = append(issueNode, url)
			}
		} else {
			healthNode = append(healthNode, url)
		}
	}

	logger.Info("", "healthNode", healthNode, "newNode", newNode,
		"issueNode", issueNode)
	// do initiate in new node.
	if len(newNode) == len(podStatus) && len(newNode) == int(mc.Spec.Mongo.
		Replicas) {
		return f.healer.MongoReplSetInitiate(mc, newNode[0], newNode[1:]...)
	}

	if len(newNode) > 0 {
		// TODO get not recover mode health pod
		if err := f.healer.MongoReplSetAdd(mc, healthNode[0],
			newNode...); err != nil {
			return err
		}
	}
	//TODO try to heal other mongo pod.
	// fixme if we should handle long time recover mode.

	return nil
}
