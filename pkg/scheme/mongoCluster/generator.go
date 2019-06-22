package mongoCluster

import (
	//"fmt"

	//"github.com/lithammer/dedent"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	//policyv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/util/intstr"

	dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	"github.smartx.com/mongo-operator/pkg/constants"
	"github.smartx.com/mongo-operator/pkg/utils"
	//"github.com/go-openapi/spec"
)


func GenerateMCService(mc *dbv1alpha1.MongoCluster,
	labels map[string]string) *corev1.Service {
	name := utils.GetMCName(mc)
	namespace := mc.Namespace

	labels = utils.MergeLabels(labels, utils.GetLabels(constants.
		MCRoleName, mc.Name))

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					Port:     constants.MongoPort,
					Protocol: corev1.ProtocolTCP,
					Name:     constants.MongoPortName,
				},
			},
			Selector: labels,
		},
	}
}

// getMongoCommand generate mongo start command from `MongoCluster.Spec`
func getMongoCommand(mc *dbv1alpha1.MongoCluster) []string {

	if len(mc.Spec.Mongo.Command) > 0 {
		return mc.Spec.Mongo.Command
	}
	return []string{
		"mongod",
		"--wiredTigerCacheSizeGB",
		"0.25",
		"--bind_ip",
		"0.0.0.0",
		"--replSet",
		"rs0",
		"--smallfiles",
		"--noprealloc",
	}
}

// getMCResources
func getMCResources(spec dbv1alpha1.MongoClusterSpec) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: getRequests(spec.Mongo.Resources),
		Limits:   getLimits(spec.Mongo.Resources),
	}
}

func getLimits(resources dbv1alpha1.MongoResources) corev1.ResourceList {
	return generateResourceList(resources.Limits.CPU, resources.Limits.Memory)
}

func getRequests(resources dbv1alpha1.MongoResources) corev1.ResourceList {
	return generateResourceList(resources.Requests.CPU, resources.Requests.Memory)
}

func generateResourceList(cpu string, memory string) corev1.ResourceList {
	resources := corev1.ResourceList{}
	if cpu != "" {
		resources[corev1.ResourceCPU], _ = resource.ParseQuantity(cpu)
	}
	if memory != "" {
		resources[corev1.ResourceMemory], _ = resource.ParseQuantity(memory)
	}
	return resources
}

func getMCVolumeMounts(mc *dbv1alpha1.MongoCluster) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name: getMongoDataVolumeName(mc),
			MountPath: "/data/db",
		},
	}

	return volumeMounts
}

// getMongoVolumeName describe the mongo volume name.
func getMongoDataVolumeName(mc *dbv1alpha1.MongoCluster) string {
	return constants.MongoStorageVolumeName
}

// getMongoVolumes return all used volume like configMap, pv, secrets.
func getMongoVolumes(mc * dbv1alpha1.MongoCluster) []corev1.Volume {
	// TODO(vici) TODO...

	volumes := []corev1.Volume{}
	dataVolume := getMongoDataVolume(mc)
	volumes = append(volumes, *dataVolume)
	return volumes
}

func getMongoDataVolume(mc *dbv1alpha1.MongoCluster) *corev1.Volume {
	// TODO(vici) ....
	return nil
}

func getMongoVolumeClaimTemplates(mc *dbv1alpha1.MongoCluster) []corev1.
	PersistentVolumeClaim {
		return []corev1.PersistentVolumeClaim{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: getMongoDataVolumeName(mc),
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.ReadWriteOnce,
					},
					StorageClassName: &mc.Spec.Mongo.Storage.StorageClassName,
					Resources: mc.Spec.Mongo.Storage.Resources,
				},
			},
		}
}

func getGraceTime() *int64 {
	 graceTime := int64(constants.GraceTime)
	 return &graceTime
}

// GenerateMCStatefulSet generate a standard mongoCluster statefulset
func GenerateMCStatefulSet(mc *dbv1alpha1.MongoCluster,
	labels map[string]string) *appsv1.StatefulSet {
	name := utils.GetMCName(mc)
	namespace := mc.Namespace

	MongoCommand := getMongoCommand(mc)
	resources := getMCResources(mc.Spec)
	labels = utils.MergeLabels(labels, utils.GetLabels(constants.MCRoleName,
		mc.Name))
	volumeMounts := getMCVolumeMounts(mc)
	//volumes := getMongoVolumes(mc)

	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,  // {MC.Name}-mongo
			Namespace: namespace, // default
			Labels:    labels, // app={MC.Name};
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    &mc.Spec.Mongo.Replicas,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Tolerations: mc.Spec.Mongo.Tolerations,
					TerminationGracePeriodSeconds: getGraceTime(),
					Containers: []corev1.Container{
						{
							Name:            "mongo",
							Image:           mc.Spec.Mongo.Image,
							ImagePullPolicy: mc.Spec.Mongo.ImagePullPolicy,
							Ports: []corev1.ContainerPort{
								{
									Name:          "mongo",
									ContainerPort: 27017,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: volumeMounts,
							Command:      MongoCommand,
							//LivenessProbe: &corev1.Probe{
							//	InitialDelaySeconds: constants.GraceTime,
							//	TimeoutSeconds:      5,
							//	Handler: corev1.Handler{
							//		Exec: &corev1.ExecAction{
							//			Command: []string{
							//				"sh",
							//				"-c",
							//				"mongo --evel 'db.runCommand({ping:1})'",
							//			},
							//		},
							//	},
							//},
							Resources: resources,
							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{"/bin/sh", "-c",
										"echo  'TODO some preStop script.'"},
									},
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: getMongoVolumeClaimTemplates(mc),
		},
	}
	// persistentVolumeClaim required.

	return ss
}
