package objsyncer

import (
	"k8s.io/apimachinery/pkg/runtime"
)

var controllerLabels = map[string]string{
	"app.kubernetes.io/managed-by": "mongo-operator",
}

var noFunc = func(existing runtime.Object) error {
	return nil
}
