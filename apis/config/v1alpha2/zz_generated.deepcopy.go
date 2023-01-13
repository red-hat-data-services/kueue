//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022 The Kubernetes Authors.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha2

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClientConnection) DeepCopyInto(out *ClientConnection) {
	*out = *in
	if in.QPS != nil {
		in, out := &in.QPS, &out.QPS
		*out = new(float32)
		**out = **in
	}
	if in.Burst != nil {
		in, out := &in.Burst, &out.Burst
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClientConnection.
func (in *ClientConnection) DeepCopy() *ClientConnection {
	if in == nil {
		return nil
	}
	out := new(ClientConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Configuration) DeepCopyInto(out *Configuration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Namespace != nil {
		in, out := &in.Namespace, &out.Namespace
		*out = new(string)
		**out = **in
	}
	in.ControllerManagerConfigurationSpec.DeepCopyInto(&out.ControllerManagerConfigurationSpec)
	if in.InternalCertManagement != nil {
		in, out := &in.InternalCertManagement, &out.InternalCertManagement
		*out = new(InternalCertManagement)
		(*in).DeepCopyInto(*out)
	}
	if in.WaitForPodsReady != nil {
		in, out := &in.WaitForPodsReady, &out.WaitForPodsReady
		*out = new(WaitForPodsReady)
		**out = **in
	}
	if in.ClientConnection != nil {
		in, out := &in.ClientConnection, &out.ClientConnection
		*out = new(ClientConnection)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Configuration.
func (in *Configuration) DeepCopy() *Configuration {
	if in == nil {
		return nil
	}
	out := new(Configuration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Configuration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InternalCertManagement) DeepCopyInto(out *InternalCertManagement) {
	*out = *in
	if in.Enable != nil {
		in, out := &in.Enable, &out.Enable
		*out = new(bool)
		**out = **in
	}
	if in.WebhookServiceName != nil {
		in, out := &in.WebhookServiceName, &out.WebhookServiceName
		*out = new(string)
		**out = **in
	}
	if in.WebhookSecretName != nil {
		in, out := &in.WebhookSecretName, &out.WebhookSecretName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InternalCertManagement.
func (in *InternalCertManagement) DeepCopy() *InternalCertManagement {
	if in == nil {
		return nil
	}
	out := new(InternalCertManagement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WaitForPodsReady) DeepCopyInto(out *WaitForPodsReady) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WaitForPodsReady.
func (in *WaitForPodsReady) DeepCopy() *WaitForPodsReady {
	if in == nil {
		return nil
	}
	out := new(WaitForPodsReady)
	in.DeepCopyInto(out)
	return out
}
