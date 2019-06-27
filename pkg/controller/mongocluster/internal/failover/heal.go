package failover

import (
	k8s "github.smartx.com/mongo-operator/pkg/service/kubernetes"
	"github.smartx.com/mongo-operator/pkg/service/mongo/replicaset"
	"github.smartx.com/mongo-operator/pkg/service/mongo"
	"github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	"github.smartx.com/mongo-operator/pkg/constants"
	"github.com/globalsign/mgo"
)

type MongoClusterFailoverHealer struct {
	k8sService k8s.Services
}

func NewMongoClusterFailoverHealer(k8sService k8s.Services) *MongoClusterFailoverHealer {
	return &MongoClusterFailoverHealer{
		k8sService: k8sService,
	}
}

// MongoReplSetInitiate create mongo replset cluster.
func (h *MongoClusterFailoverHealer) MongoReplSetInitiate(
	mc *v1alpha1.MongoCluster, master string, members ...string) error {

	var replSetLabels = map[string]string{
		"app.kubernetes.io/managed-by": "mongo-operator",
		"cluster": mc.Name,
	}

	mongoClient := mongo.NewClient(master)
	mgoSession, err := mongoClient.DialDirect()
	defer mgoSession.Close()
	if err != nil {
		return err
	}
	if err := replicaset.Initiate(mgoSession, master,
		constants.MongoReplSetName,
		replSetLabels); err != nil {
		return err
	}
	// add pods to cluster.
	if err := mongoReplsetAddMemebers(mgoSession,
		replSetLabels, members...); err !=  nil {
		return err
	}

	return nil
}

func(h *MongoClusterFailoverHealer) MongoReplSetAdd(mc *v1alpha1.
	MongoCluster, tags map[string]string, members ...string) error {
	return nil
}

func mongoReplsetAddMemebers(session *mgo.Session, tags map[string]string,
	member ...string) error {
	s := session.Clone()
	defer s.Close()

	members := make([]replicaset.Member, 0, len(member))
	for _, m:= range member {
		members = append(members, replicaset.Member{
			Address: m,
			Tags: tags,
		})
	}
	replicaset.Add(s, members...)
	return nil
}
