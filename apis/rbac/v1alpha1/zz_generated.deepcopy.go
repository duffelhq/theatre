//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/rbac/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DirectoryRoleBinding) DeepCopyInto(out *DirectoryRoleBinding) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DirectoryRoleBinding.
func (in *DirectoryRoleBinding) DeepCopy() *DirectoryRoleBinding {
	if in == nil {
		return nil
	}
	out := new(DirectoryRoleBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DirectoryRoleBinding) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DirectoryRoleBindingList) DeepCopyInto(out *DirectoryRoleBindingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DirectoryRoleBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DirectoryRoleBindingList.
func (in *DirectoryRoleBindingList) DeepCopy() *DirectoryRoleBindingList {
	if in == nil {
		return nil
	}
	out := new(DirectoryRoleBindingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DirectoryRoleBindingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DirectoryRoleBindingSpec) DeepCopyInto(out *DirectoryRoleBindingSpec) {
	*out = *in
	if in.Subjects != nil {
		in, out := &in.Subjects, &out.Subjects
		*out = make([]v1.Subject, len(*in))
		copy(*out, *in)
	}
	out.RoleRef = in.RoleRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DirectoryRoleBindingSpec.
func (in *DirectoryRoleBindingSpec) DeepCopy() *DirectoryRoleBindingSpec {
	if in == nil {
		return nil
	}
	out := new(DirectoryRoleBindingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DirectoryRoleBindingStatus) DeepCopyInto(out *DirectoryRoleBindingStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DirectoryRoleBindingStatus.
func (in *DirectoryRoleBindingStatus) DeepCopy() *DirectoryRoleBindingStatus {
	if in == nil {
		return nil
	}
	out := new(DirectoryRoleBindingStatus)
	in.DeepCopyInto(out)
	return out
}
