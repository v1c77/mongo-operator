package failover

import (
	"fmt"
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

	// TODO(yuhua) some pre-check
	// check 网络 是否连通
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
	// - check pod, service, network, mongo replicaSet..

	return nil
}

// checkAndHealMongoReplicaSet
func (f *MongoClusterFailover) checkAndHealMongoReplSet(
	mc *dbv1alpha1.MongoCluster) error {

	podStatus := f.checker.GetMongoPodsStatus(mc)
	healthNode := make(map[string]podReplicaStatus)
	newNode := make(map[string]podReplicaStatus)
	issueNode := make(map[string]podReplicaStatus)
	for url, podStatus := range podStatus {
		// init mongo replicaset.
		if podStatus.Err != nil {
			if strings.Contains(podStatus.Err.Error(),
				"no replset config has been received") {
				newNode[url] = podStatus
			} else {
				// + pod  mongo daemon initialing.
				// + pod restart
				// + mongo rs status recovering ??
				// + TODO other issue status.
				issueNode[url] = podStatus
			}
		} else {
			healthNode[url] = podStatus
		}
	}
	logger.Debug("all node info",
		"healthNode", GetMapStringKeys(healthNode),
		"newNode", GetMapStringKeys(newNode),
		"issueNode", GetMapStringKeys(issueNode))

	// ---- 1. init  cluster
	// TODO to make sure the cluster config contains odd members .
	if len(newNode) == len(podStatus) && len(newNode) == int(mc.Spec.Mongo.
		Replicas) {
		nodes, err := getMapStringKeys(newNode)
		if err != nil {
			return err
		}

		logger.Debug("init cluster",
			"master", nodes[0],
			"members", nodes[1:])
		return f.healer.MongoReplSetInitiate(mc, nodes[0], nodes[1:]...)
	}

	// ---- 2. handle scala, newly added node
	// TODO to make sure the cluster config contains odd members .
	//if len(newNode) > 0  && len(newNode) % 2 == 0 {
	if len(newNode) > 0 {
		// TODO get not recover mode health pod
		var master string
		for url, node := range healthNode {
			if node.IsMaster != nil && node.IsMaster.IsMaster {
				master = url
				break
			}
		}
		if len(master) > 0 { // primary must exist.

			nodes, err := getMapStringKeys(newNode)
			if err != err {
				return err
			}
			logger.Debug("add new node",
				"master", master,
				"new node", newNode)
			return f.healer.MongoReplSetAdd(mc, master, nodes...)
		} else {
			return fmt.Errorf("no health master node")
		}
	}
	// TODO reduce node.
	// TODO try to heal other mongo pod.
	// fixme if we should handle long time recover mode.
	logger.Debug("TODO situation ignored.",
		"healthNode", GetMapStringKeys(healthNode),
		"newNode", GetMapStringKeys(newNode),
		"issueNode", GetMapStringKeys(issueNode))
	return nil
}
