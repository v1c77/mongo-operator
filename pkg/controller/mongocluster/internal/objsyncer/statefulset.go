package objsyncer



import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.smartx.com/mongo-operator/pkg/staging/syncer"
	"github.smartx.com/mongo-operator/pkg/scheme/mongoCluster"
	 dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
)

// NewMongoStatefulSetSyncer returns a new sync.
// Interface for reconciling Mongo StatefulSet
func NewMongoStatefulSetSyncer(mc *dbv1alpha1.MongoCluster, c client.Client,
	scheme *runtime.Scheme) syncer.Interface {
	statefulSet := mongoCluster.GenerateMCStatefulSet(mc, controllerLabels)
	return syncer.NewObjectSyncer("MongoStatefulSet", mc, statefulSet, c,
		scheme, noFunc)
}