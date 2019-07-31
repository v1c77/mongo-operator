package v1alpha1

// set default for MC spec.
import corev1 "k8s.io/api/core/v1"

const (
	DefaultMongoReplica = 3
	DefaultMongoImage   = "mongo:4.0.11"
	DefaultImagePullPolicy = corev1.PullIfNotPresent
	DefaultWiredTigerCacheSize = "0.25"
	DefaultBindIp = "0.0.0.0"
	DefaultReplSet = "rs0"
	DefaultRequestCpu = "1"
	DefaultRequestRam = "3G"

)

func (in *MongoCluster) SetDefaults() {

	if in.Spec.Mongo.Replicas == 0 {
		in.Spec.Mongo.Replicas = DefaultMongoReplica
	}

	if len(in.Spec.Mongo.Image) == 0 {
		in.Spec.Mongo.Image = DefaultMongoImage
	}

	if len(in.Spec.Mongo.ImagePullPolicy) == 0 {
		in.Spec.Mongo.ImagePullPolicy = DefaultImagePullPolicy
	}

	if len(in.Spec.Mongo.ReplSet) == 0 {
		in.Spec.Mongo.ReplSet = DefaultReplSet
	}

	if len(in.Spec.Mongo.WiredTigerCacheSize) == 0 {
		in.Spec.Mongo.WiredTigerCacheSize = DefaultWiredTigerCacheSize
	}

	if len(in.Spec.Mongo.BindIp) == 0 {
		in.Spec.Mongo.BindIp = DefaultBindIp
	}

	if in.Spec.Mongo.Resources == nil {
		in.Spec.Mongo.Resources = &MongoResources{Requests:CPUAndMem{
			CPU: DefaultRequestCpu,
			Memory: DefaultRequestRam,
		}}
	}
}
