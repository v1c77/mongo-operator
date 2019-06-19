package utils

import (

	"fmt"

	dbv1alpha1 "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1"
	"github.smartx.com/mongo-operator/pkg/constants"
)

func GetMCName(mc *dbv1alpha1.MongoCluster) string {
return generateName(constants.MongoName, mc.Name)
}


func generateName(typeName, metaName string) string {
	return fmt.Sprintf("%s-%s", typeName, metaName)
}