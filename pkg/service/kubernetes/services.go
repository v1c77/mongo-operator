package kubernetes

import (
	"k8s.io/client-go/kubernetes"
)

// Service is the K8s service entrypoint.
type Services interface {
	StatefulSet
}

type services struct {
	StatefulSet
}

// New returns a new Kubernetes service.
func New(kubeClient kubernetes.Interface) Services {
	return &services{
		StatefulSet: NewStatefulSetService(kubeClient),
	}
}