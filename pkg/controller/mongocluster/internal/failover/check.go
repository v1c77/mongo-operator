package failover

import (
	"github.smartx.com/mongo-operator/pkg/service/mongo"
	k8s "github.smartx.com/mongo-operator/pkg/service/kubernetes"
	dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	"fmt"
	"github.smartx.com/mongo-operator/pkg/utils"
	"github.com/globalsign/mgo"
)

type MongoClusterFailoverChecker struct {
	k8sService k8s.Services
	MongoClient mongo.Client
}

func NewMongoClusterFailoverChecker(k8sService k8s.Services,
	mongoClient mongo.Client) *MongoClusterFailoverChecker {
	return &MongoClusterFailoverChecker{
		k8sService:  k8sService,
		MongoClient: mongoClient,
	}
}



func (c *MongoClusterFailoverChecker) GetReplicaSetStatus()  {
}

func (c *MongoClusterFailoverChecker) GetMembersDNS(mc *dbv1alpha1.
	MongoCluster) []string {
	var dnsList []string
	replicaCount := int(mc.Spec.Mongo.Replicas)
	clusterName := utils.GetMCName(mc)
	namespace := mc.Namespace
	for idx:= 0; idx < replicaCount; idx++ {
		dnsList = append(dnsList, getMemberHostName(idx, clusterName, namespace))
	}
	return dnsList
}

func (c *MongoClusterFailoverChecker) MemebersStatus(mc *dbv1alpha1.
	MongoCluster) error {
		// TODO refactor it !!!
		dnsList := c.GetMembersDNS(mc)
		url := dnsList[0]
		Session, err := mgo.Dial(url)
		if err != nil {
			log.Info("can not get client", "error", err)
			return err
		}
		err = Session.Ping()
		if err != nil {
			log.Info("mongod not work or network error", "error", err)
			return err
		}
		log.Info("check random pod status done.")
		return nil
}


func getMemberHostName(idx int, clusterName, namespace string) string {
	return fmt.Sprintf("%s-%v.%s.%s.svc.cluster.local", clusterName, idx,
		clusterName, namespace)
}
