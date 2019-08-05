package mongoCluster

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	"github.smartx.com/mongo-operator/pkg/constants"
	"github.smartx.com/mongo-operator/pkg/utils"
	"fmt"
	"strings"
)

func GetMcServiceName(mc *dbv1alpha1.MongoCluster) string {
	return utils.GetMCName(mc)
}

func GenerateMCService(mc *dbv1alpha1.MongoCluster,
	labels map[string]string) *corev1.Service {
	name := GetMcServiceName(mc)
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

func GetPodsFQDN(mc *dbv1alpha1.MongoCluster) []string {
	podNames := utils.GetStatefulsetPodNames(GenerateMCStatefulSet(mc,
		map[string]string{}))

	seviceFQDN := utils.GetServiceFQDN(GenerateMCService(mc,
		map[string]string{}))

	podsFQDN := make([]string, 0, len(podNames))
	for _, podName := range podNames {
		podsFQDN = append(podsFQDN, fmt.Sprintf("%s.%s", podName, seviceFQDN))
	}
	return podsFQDN
}

func GetMCConfigMapName(mc *dbv1alpha1.MongoCluster) string {
	return fmt.Sprintf("%s-conf", utils.GetMCName(mc))
}
func GenerateConfigMap(mc *dbv1alpha1.MongoCluster) *corev1.ConfigMap {
	podsFQDN := GetPodsFQDN(mc)
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetMCConfigMapName(mc),
			Namespace: mc.Namespace,
		},
		Data: map[string]string{
			"cluster.mongo": strings.Join(podsFQDN, ","),
		},
	}

}

// getMongoCommand generate mongo start command from `MongoCluster.Spec`
func getMongoCommand(mc *dbv1alpha1.MongoCluster) []string {
	mc.SetDefaults() // idempotent
	commands := []string{"mongod"}

	// wiredTigerCacheSize
	commands = append(commands,
		"--wiredTigerCacheSizeGB", mc.Spec.Mongo.WiredTigerCacheSize)

	// bindIp
	commands = append(commands,
		"--bind_ip", mc.Spec.Mongo.BindIp)

	// replSet
	commands = append(commands,
		"--replSet", mc.Spec.Mongo.ReplSet)

	if mc.Spec.Mongo.SmallFiles {
		commands = append(commands, "--smallfiles")
	}

	if mc.Spec.Mongo.Noprealloc {
		commands = append(commands, "--noprealloc")
	}

	return commands
}

// getMCResources
func getMCResources(spec dbv1alpha1.MongoClusterSpec) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: getRequests(spec.Mongo.Resources),
		Limits:   getLimits(spec.Mongo.Resources),
	}
}

func getLimits(resources *dbv1alpha1.MongoResources) corev1.ResourceList {
	return generateResourceList(resources.Limits.CPU, resources.Limits.Memory)
}

func getRequests(resources *dbv1alpha1.MongoResources) corev1.ResourceList {
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
			Name:      getMongoDataVolumeName(mc),
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
func getMongoVolumes(mc *dbv1alpha1.MongoCluster) []corev1.Volume {
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
				Resources:        mc.Spec.Mongo.Storage.Resources,
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
	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,      // {MC.Name}-mongo
			Namespace: namespace, // default
			Labels:    labels,    // app={MC.Name};
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    &mc.Spec.Mongo.Replicas,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: constants.UpdatePolicy,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Tolerations:                   mc.Spec.Mongo.Tolerations,
					TerminationGracePeriodSeconds: getGraceTime(),
					Containers: []corev1.Container{
						{
							Name:            constants.MongoName,
							Image:           mc.Spec.Mongo.Image,
							ImagePullPolicy: mc.Spec.Mongo.ImagePullPolicy,
							Ports: []corev1.ContainerPort{
								{
									Name:          constants.MongoName,
									ContainerPort: 27017,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: volumeMounts,
							Command:      MongoCommand,
							// TODO(yuhua): LivenessProbe
							Resources: resources,
							// TODO(yuhua): Lifecycle
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
