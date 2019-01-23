// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/rbac/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Console) DeepCopyInto(out *Console) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Console.
func (in *Console) DeepCopy() *Console {
	if in == nil {
		return nil
	}
	out := new(Console)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Console) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsoleList) DeepCopyInto(out *ConsoleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Console, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConsoleList.
func (in *ConsoleList) DeepCopy() *ConsoleList {
	if in == nil {
		return nil
	}
	out := new(ConsoleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ConsoleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsoleSpec) DeepCopyInto(out *ConsoleSpec) {
	*out = *in
	out.ConsoleTemplateRef = in.ConsoleTemplateRef
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConsoleSpec.
func (in *ConsoleSpec) DeepCopy() *ConsoleSpec {
	if in == nil {
		return nil
	}
	out := new(ConsoleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsoleStatus) DeepCopyInto(out *ConsoleStatus) {
	*out = *in
	if in.ExpiryTime != nil {
		in, out := &in.ExpiryTime, &out.ExpiryTime
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConsoleStatus.
func (in *ConsoleStatus) DeepCopy() *ConsoleStatus {
	if in == nil {
		return nil
	}
	out := new(ConsoleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsoleTemplate) DeepCopyInto(out *ConsoleTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConsoleTemplate.
func (in *ConsoleTemplate) DeepCopy() *ConsoleTemplate {
	if in == nil {
		return nil
	}
	out := new(ConsoleTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ConsoleTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsoleTemplateList) DeepCopyInto(out *ConsoleTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ConsoleTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConsoleTemplateList.
func (in *ConsoleTemplateList) DeepCopy() *ConsoleTemplateList {
	if in == nil {
		return nil
	}
	out := new(ConsoleTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ConsoleTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsoleTemplateSpec) DeepCopyInto(out *ConsoleTemplateSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
	if in.AdditionalAttachSubjects != nil {
		in, out := &in.AdditionalAttachSubjects, &out.AdditionalAttachSubjects
		*out = make([]v1.Subject, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConsoleTemplateSpec.
func (in *ConsoleTemplateSpec) DeepCopy() *ConsoleTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(ConsoleTemplateSpec)
	in.DeepCopyInto(out)
	return out
}
