package v1alpha1

// set default for MC spec.

const (
	defaultMongoReplica = 3
	defaultMongoImage   = "mongo"
	//defaultImagePullPolicy = ???
	//defaultExporterImage = "prometheus-mongo-exporter"

)

func (in *MongoCluster) SetDefaults() {

	if in.Spec.Mongo.Replicas == 0 {
		in.Spec.Mongo.Replicas = defaultMongoReplica
	}

	if len(in.Spec.Mongo.Image) == 0 {
		in.Spec.Mongo.Image = defaultMongoImage
	}

	//if len(mc.Spec.Mongo.ImagePullPolicy) == 0 {
	//	mc.Spec.Mongo.ImagePullPolicy = defaultImagePullPolicy
	//}

	//if len(mc.Spec.Mongo.Exporter.Image) == 0 {
	//	mc.Spec.Mongo.Exporter.Image = defaultExporterImage
	//}

	//if len(mc.Spec.Mongo.Exporter.ImagePullPolicy) == 0 {
	//	mc.Spec.Mongo.Exporter.ImagePullPolicy = defaultImagePullPolicy
	//}
}
