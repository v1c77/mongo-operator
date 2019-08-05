package objsyncer

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	"github.smartx.com/mongo-operator/pkg/scheme/mongoCluster"
	"github.smartx.com/mongo-operator/pkg/staging/syncer"
)

// NewMongoServiceSyncer returns a new sync.
// Interface for reconciling Mongo headless service
func NewMongoConfigMap(mc *dbv1alpha1.MongoCluster, c client.Client,
	scheme *runtime.Scheme) syncer.Interface {
	cm := mongoCluster.GenerateConfigMap(mc)
	return syncer.NewObjectSyncer("MongoConfigMap", mc, cm, c,
		scheme, noFunc)
}
