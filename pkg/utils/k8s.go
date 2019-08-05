package utils

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func GetServiceFQDN(s *corev1.Service) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local",
		s.Name, s.Namespace)
}

func GetStatefulsetPodNames(s *appsv1.StatefulSet) []string {
	var podNames []string
	for idx := 0; idx < int(*s.Spec.Replicas); idx++ {
		podNames = append(podNames, getMemberHostName(idx, s.Name))
	}
	return podNames
}

func getMemberHostName(idx int, name string) string {
	return fmt.Sprintf("%s-%v", name, idx)
}
