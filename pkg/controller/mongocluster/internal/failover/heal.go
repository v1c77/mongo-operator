package failover

import (
	"github.com/globalsign/mgo"
	"github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	k8s "github.smartx.com/mongo-operator/pkg/service/kubernetes"
	"github.smartx.com/mongo-operator/pkg/service/mongo"
	"github.smartx.com/mongo-operator/pkg/service/mongo/replicaset"
	"strings"
)

type MongoClusterFailoverHealer struct {
	k8sService k8s.Services
}

func NewMongoClusterFailoverHealer(k8sService k8s.Services) *MongoClusterFailoverHealer {
	return &MongoClusterFailoverHealer{
		k8sService: k8sService,
	}
}

func getReplsetTags(mc *v1alpha1.MongoCluster) map[string]string {

	return map[string]string{
		"app.kubernetes.io/managed-by": "mongo-operator",
		"cluster":                      mc.Name,
	}
}

// MongoReplSetInitiate create mongo replset cluster.
func (h *MongoClusterFailoverHealer) MongoReplSetInitiate(
	mc *v1alpha1.MongoCluster, master string, members ...string) error {

	replSetLabels := getReplsetTags(mc)
	mongoClient := mongo.NewClient(master)
	mgoSession, err := mongoClient.DialDirect()
	defer mgoSession.Close()
	if err != nil {
		return err
	}

	if err := replicaset.Initiate(mgoSession, master,
		mc.Spec.Mongo.ReplSet,
		replSetLabels); err != nil {
		return err
	}
	// add pods to cluster.
	logger.Debug("mongo init cluster: add members",
		"master", master,
		"members", members)
	if err := mongoReplsetAddMemebers(mgoSession,
		replSetLabels, members...); err != nil {
		return err
	}
	return nil
}

func (h *MongoClusterFailoverHealer) MongoReplSetAdd(mc *v1alpha1.
	MongoCluster, clusterAddr string, members ...string) error {
	replSetLabels := getReplsetTags(mc)
	mongoClient := mongo.NewClient(clusterAddr)
	mgoSession, err := mongoClient.Dial()
	defer mgoSession.Close()
	if err != nil {
		return err
	}
	return mongoReplsetAddMemebers(mgoSession, replSetLabels, members...)
}

func mongoReplsetAddMemebers(session *mgo.Session, tags map[string]string,
	newMembers ...string) error {
	s := session.Clone()
	defer s.Close()

	currentConfig, _ := replicaset.CurrentConfig(s)
	// check if member already in current cluster.

	// reconfig exist member
	var toDelete []string
	for _, currMemeber := range currentConfig.Members {
		for _, newMember := range newMembers {
			if strings.Contains(currMemeber.Address, newMember) {
				toDelete = append(toDelete, newMember)
			}
		}
	}
	if len(toDelete) > 0 {
		logger.Debug("try to re-config some node", "node", toDelete)
		replicaset.Remove(s, toDelete...)
	}

	members := make([]replicaset.Member, 0, len(newMembers))
	for _, m := range newMembers {
		members = append(members, replicaset.Member{
			Address: m,
			Tags:    tags,
		})
	}
	replicaset.Add(s, members...)
	return nil
}
