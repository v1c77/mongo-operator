package failover

import (
	"fmt"
	"github.com/globalsign/mgo"
	dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	k8s "github.smartx.com/mongo-operator/pkg/service/kubernetes"
	"github.smartx.com/mongo-operator/pkg/service/mongo/replicaset"
	"github.smartx.com/mongo-operator/pkg/utils"
	"github.smartx.com/mongo-operator/pkg/service/mongo"
	"strings"
)

type MongoClusterFailoverChecker struct {
	k8sService k8s.Services
}

func NewMongoClusterFailoverChecker(k8sService k8s.Services) *MongoClusterFailoverChecker {
	return &MongoClusterFailoverChecker{
		k8sService: k8sService,
	}
}

func (c *MongoClusterFailoverChecker) CheckReplSetStatus(session *mgo.
	Session) (
	*replicaset.Status, error) {
	monotonicSession := session.Clone()
	defer monotonicSession.Close()
	monotonicSession.SetMode(mgo.Monotonic, true)

	return replicaset.CurrentStatus(monotonicSession)
}

func (c *MongoClusterFailoverChecker) GetMembersDNS(mc *dbv1alpha1.
	MongoCluster) []string {
	var dnsList []string
	replicaCount := int(mc.Spec.Mongo.Replicas)
	clusterName := utils.GetMCName(mc)
	namespace := mc.Namespace
	for idx := 0; idx < replicaCount; idx++ {
		dnsList = append(dnsList, getMemberHostName(idx, clusterName, namespace))
	}
	return dnsList
}

type podReplicaStatus struct {
	Status *replicaset.Status
	Err error
	IsReplica bool
}

// checkMongoPodsStatus check all alive mongo instance status
func (c *MongoClusterFailoverChecker) GetMongoPodsStatus(mc *dbv1alpha1.
MongoCluster) map[string]podReplicaStatus {
	dnsList :=  c.GetMembersDNS(mc)
	var podsMap = map[string]podReplicaStatus{}
	for _, url := range dnsList {
		mongoClient := mongo.NewClient(url)
		mgoSession, err := mongoClient.DialDirect()
		if err != nil {
			// mongod not started or network error pods
			podsMap[url] = podReplicaStatus{
				Status: nil,
				Err: err,
				IsReplica: false, //ignore this type pods until Dial connected.x
			}
			continue
		}
		status, err := c.CheckReplSetStatus(mgoSession)
		if err != nil {
			if strings.Contains(err.Error(),
				"no replset config has been received") {
				podsMap[url] = podReplicaStatus{
					Status: status,
					Err: err,
					IsReplica: false,
				}
				continue
			}
		}
		podsMap[url] = podReplicaStatus{
			Status: status,
			Err: err,
			IsReplica: true,
		}

		mgoSession.Close()
	}
	return podsMap
}


//func (c *MongoClusterFailoverChecker) MemebersStatus(mc *dbv1alpha1.
//	MongoCluster) error {
//		// TODO refactor it !!!
//		dnsList := c.GetMembersDNS(mc)
//		url := fmt.Sprintf("%s:%v", dnsList[0], constants.MongoPort)
//		//url := "10.76.203.198:27017"
//		logger.Info("use url", "url", url)
//		Session, err := mgo.Dial(url)
//		if err != nil {
//			log.Info("can not get client", "error", err)
//			return err
//		}
//		mgo.Monotonic
//		err = Session.Ping()
//		if err != nil {
//			log.Info("mongod not work or network error", "error", err)
//			return err
//		}
//		log.Info("check random pod status done.")
//		return nil
//}

func getMemberHostName(idx int, clusterName, namespace string) string {
	return fmt.Sprintf("%s-%v.%s.%s.svc.cluster.local", clusterName, idx,
		clusterName, namespace)
}
