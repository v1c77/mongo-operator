package failover

import (
	k8s "github.smartx.com/mongo-operator/pkg/service/kubernetes"
	"github.smartx.com/mongo-operator/pkg/service/mongo/replicaset"
	"github.smartx.com/mongo-operator/pkg/service/mongo"
)

type MongoClusterFailoverHealer struct {
	k8sService k8s.Services
}

func NewMongoClusterFailoverHealer(k8sService k8s.Services) *MongoClusterFailoverHealer {
	return &MongoClusterFailoverHealer{
		k8sService: k8sService,
	}
}

func (h *MongoClusterFailoverHealer) MongoReplSetInitial(url string,
	tags map[string]string) error {

		mongo.NewClient(url)

		replicaset.Initiate()
}