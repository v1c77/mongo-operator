// +build !ignore_autogenerated

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1.MongoCluster":       schema_pkg_apis_db_v1alpha1_MongoCluster(ref),
		"github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1.MongoClusterSpec":   schema_pkg_apis_db_v1alpha1_MongoClusterSpec(ref),
		"github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1.MongoClusterStatus": schema_pkg_apis_db_v1alpha1_MongoClusterStatus(ref),
	}
}

func schema_pkg_apis_db_v1alpha1_MongoCluster(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "MongoCluster is the Schema for the mongoclusters API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1.MongoClusterSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1.MongoClusterStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1.MongoClusterSpec", "github.smartx.com/mongo-operator/pkg/apis/db/v1alpha1.MongoClusterStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_db_v1alpha1_MongoClusterSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "MongoClusterSpec defines the desired state of MongoCluster",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_db_v1alpha1_MongoClusterStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "MongoClusterStatus defines the observed state of MongoCluster",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
