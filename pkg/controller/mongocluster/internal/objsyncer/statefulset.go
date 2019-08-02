package objsyncer

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	"github.smartx.com/mongo-operator/pkg/scheme/mongoCluster"
	"github.smartx.com/mongo-operator/pkg/staging/syncer"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"github.smartx.com/mongo-operator/pkg/constants"
	"fmt"
)

// NewMongoStatefulSetSyncer returns a new sync.
// Interface for reconciling Mongo StatefulSet
func NewMongoStatefulSetSyncer(mc *dbv1alpha1.MongoCluster, c client.Client,
	scheme *runtime.Scheme) syncer.Interface {
	news := mongoCluster.GenerateMCStatefulSet(mc, controllerLabels)
	return syncer.NewObjectSyncer("MongoStatefulSet", mc, news.DeepCopy(), c, scheme,
		func(exist runtime.Object) error {
			// TODO (vici) more meticulous
			// user can only
			// - scala cluster
			out := exist.(*appsv1.StatefulSet) // exist mongo statefulSet
			syncReplica(news, out)
			syncMongoPod(news, out)
			return nil
		})
}

type mongoPod struct {
	 fsts *appsv1.StatefulSet
	 cIdx int
	 p corev1.Container
}

func NewMongoPodFromSts(sts *appsv1.StatefulSet) *mongoPod {
	return &mongoPod{
		fsts: sts,
	}
}

func (p *mongoPod) getMongoPod()  error {
	for idx, c := range p.fsts.Spec.Template.Spec.Containers {
		if c.Name == constants.MongoName {
			p.cIdx = idx
			p.p = c
			return nil
		}
	}
	return fmt.Errorf("can not find mongo pod")
}

func (p *mongoPod) Pod() *corev1.Container {

	if len(p.p.Name) == 0 {
		p.getMongoPod()
	}
	return &p.p

}

func (p *mongoPod) Commit() {
	if len(p.p.Name) == 0 {
		return
	}
	p.fsts.Spec.Template.Spec.Containers[p.cIdx] = p.p
}
func syncMongoPod(new *appsv1.StatefulSet, exist *appsv1.StatefulSet) {
	newMongoPod := NewMongoPodFromSts(new)
	oldMongoPod := NewMongoPodFromSts(exist)

	// images
	if !reflect.DeepEqual(newMongoPod.Pod().Image, oldMongoPod.Pod().Image) {
		oldMongoPod.Pod().Image = newMongoPod.Pod().Image
	}

	// Resource
	if !reflect.DeepEqual(newMongoPod.Pod().Resources, oldMongoPod.Pod().Resources) {
		oldMongoPod.Pod().Resources = newMongoPod.Pod().Resources
	}
	// cmd
	if !reflect.DeepEqual(newMongoPod.Pod().Command, oldMongoPod.Pod().Command) {
		oldMongoPod.Pod().Command = newMongoPod.Pod().Command
	}
	// command
	if !reflect.DeepEqual(newMongoPod.Pod().Command, oldMongoPod.Pod().Command) {
		oldMongoPod.Pod().Command = newMongoPod.Pod().Command
	}

	oldMongoPod.Commit()

}

func syncReplica(news *appsv1.StatefulSet, exist *appsv1.StatefulSet) {
	if !reflect.DeepEqual(exist.Spec.Replicas, news.Spec.Replicas) {
		// scala.
		logger.Info("scala mongo cluster",
			"before", exist.Spec.Replicas,
			"new", news.Spec.Replicas)
		exist.Spec.Replicas = news.Spec.Replicas
	}
}


