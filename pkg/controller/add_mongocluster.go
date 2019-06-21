package controller

import (
	"github.smartx.com/mongo-operator/pkg/controller/mongocluster"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, mongocluster.Add)
}
