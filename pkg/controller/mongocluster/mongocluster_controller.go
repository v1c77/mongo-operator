package mongocluster

import (
	"context"

	dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	"github.smartx.com/mongo-operator/pkg/controller/mongocluster/internal/objsyncer"
	"github.smartx.com/mongo-operator/pkg/staging/syncer"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_mongocluster")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MongoCluster Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMongoCluster{
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		recorder: mgr.GetRecorder("mongo-operator")}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mongocluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MongoCluster
	err = c.Watch(&source.Kind{Type: &dbv1alpha1.MongoCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(vici) change Watch for changes to the resources that owned by the
	// primary resource

	subResources := []runtime.Object{
		&corev1.Service{},
		&corev1.ConfigMap{},   // TODO ConfigMap send mongoCluster config.
		&appsv1.StatefulSet{}, // TODO STS hold all mongodb pod.
		&appsv1.Deployment{},  // TODO ==DEPRECATED==
		&corev1.Pod{},         // TODO pod add.
	}

	for _, subResource := range subResources {
		err = c.Watch(&source.Kind{Type: subResource},
			&handler.EnqueueRequestForOwner{
				IsController: true,
				OwnerType:    &dbv1alpha1.MongoCluster{},
			})
		if err != nil {
			return err
		}
	}

	return nil
}

// blank assignment to verify that ReconcileMongoCluster implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMongoCluster{}

// ReconcileMongoCluster reconciles a MongoCluster object
type ReconcileMongoCluster struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a MongoCluster object and makes changes based on the state read
// and what is in the MongoCluster.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMongoCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MongoCluster")

	// Fetch the MongoCluster instance
	mc := &dbv1alpha1.MongoCluster{}
	err := r.client.Get(context.TODO(), request.NamespacedName, mc)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	r.scheme.Default(mc)
	mc.SetDefaults()

	// use stage/syncer to manage resource changes.
	// each type of resources managed by MC has its own syncer
	// TODO(vici) mongo statefulset. + headless mongo-service.
	syncers := []syncer.Interface{
		//TODO(vici) service syncer...
		objsyncer.NewMongoStatefulSetSyncer(mc, r.client, r.scheme),
	}

	if err = r.sync(syncers); err != nil {
		return reconcile.Result{}, err
	}

	// TODO(vici) something about mongo ops failover.... add node. delete node.
	//if err = r.failover.CheckAndHeal(redis); err != nil {
	//	return reconcile.Result{}, err
	//}

	return reconcile.Result{}, nil

	//mcStatefulSet := mongoCluster.GenerateMCStatefulSet(mc)
	//
	//// Set MongoCluster instance as the owner and controller
	//if err := controllerutil.SetControllerReference(mc, pod,
	//	r.scheme); err != nil {
	//	return reconcile.Result{}, err
	//}
	//
	//// Check if this Pod already exists
	//found := &corev1.Pod{}
	//err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	//if err != nil && errors.IsNotFound(err) {
	//	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	//	err = r.client.Create(context.TODO(), pod)
	//	if err != nil {
	//		return reconcile.Result{}, err
	//	}
	//
	//	// Pod created successfully - don't requeue
	//	return reconcile.Result{}, nil
	//} else if err != nil {
	//	return reconcile.Result{}, err
	//}
	return reconcile.Result{}, nil
}

func (r *ReconcileMongoCluster) sync(syncers []syncer.Interface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s, r.recorder); err != nil {
			return err
		}
	}
	return nil
}
