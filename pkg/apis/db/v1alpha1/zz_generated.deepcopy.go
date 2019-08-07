// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUAndMem) DeepCopyInto(out *CPUAndMem) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUAndMem.
func (in *CPUAndMem) DeepCopy() *CPUAndMem {
	if in == nil {
		return nil
	}
	out := new(CPUAndMem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoCluster) DeepCopyInto(out *MongoCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoCluster.
func (in *MongoCluster) DeepCopy() *MongoCluster {
	if in == nil {
		return nil
	}
	out := new(MongoCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MongoCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoClusterList) DeepCopyInto(out *MongoClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MongoCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoClusterList.
func (in *MongoClusterList) DeepCopy() *MongoClusterList {
	if in == nil {
		return nil
	}
	out := new(MongoClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MongoClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoClusterSpec) DeepCopyInto(out *MongoClusterSpec) {
	*out = *in
	in.Mongo.DeepCopyInto(&out.Mongo)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoClusterSpec.
func (in *MongoClusterSpec) DeepCopy() *MongoClusterSpec {
	if in == nil {
		return nil
	}
	out := new(MongoClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoClusterStatus) DeepCopyInto(out *MongoClusterStatus) {
	*out = *in
	if in.ObservedGeneration != nil {
		in, out := &in.ObservedGeneration, &out.ObservedGeneration
		*out = new(int64)
		**out = **in
	}
	if in.PodsFQDN != nil {
		in, out := &in.PodsFQDN, &out.PodsFQDN
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.HealthMembers != nil {
		in, out := &in.HealthMembers, &out.HealthMembers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.IssueMembers != nil {
		in, out := &in.IssueMembers, &out.IssueMembers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoClusterStatus.
func (in *MongoClusterStatus) DeepCopy() *MongoClusterStatus {
	if in == nil {
		return nil
	}
	out := new(MongoClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoResources) DeepCopyInto(out *MongoResources) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoResources.
func (in *MongoResources) DeepCopy() *MongoResources {
	if in == nil {
		return nil
	}
	out := new(MongoResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoSettings) DeepCopyInto(out *MongoSettings) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(MongoResources)
		**out = **in
	}
	in.Storage.DeepCopyInto(&out.Storage)
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoSettings.
func (in *MongoSettings) DeepCopy() *MongoSettings {
	if in == nil {
		return nil
	}
	out := new(MongoSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoStorage) DeepCopyInto(out *MongoStorage) {
	*out = *in
	in.Resources.DeepCopyInto(&out.Resources)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoStorage.
func (in *MongoStorage) DeepCopy() *MongoStorage {
	if in == nil {
		return nil
	}
	out := new(MongoStorage)
	in.DeepCopyInto(out)
	return out
}
