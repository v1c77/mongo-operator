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
func NewMongoServiceSyncer(mc *dbv1alpha1.MongoCluster, c client.Client,
	scheme *runtime.Scheme) syncer.Interface {
	statefulSet := mongoCluster.GenerateMCService(mc, controllerLabels)
	return syncer.NewObjectSyncer("MongoService", mc, statefulSet, c,
		scheme, noFunc)
}
