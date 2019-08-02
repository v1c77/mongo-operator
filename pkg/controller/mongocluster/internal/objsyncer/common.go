package objsyncer

import (
	"github.smartx.com/mongo-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/runtime"
)

var logger = utils.NewLogger("syncer")

var controllerLabels = map[string]string{
	"app.kubernetes.io/managed-by": "mongo-operator",
}

var noFunc = func(existing runtime.Object) error {
	return nil
}
