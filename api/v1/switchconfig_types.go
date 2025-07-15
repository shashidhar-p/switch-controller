/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SwitchConfigSpec defines the desired state of SwitchConfig
type SwitchConfigSpec struct {
	// SwitchIP is the address of the switch
	// +kubebuilder:validation:Pattern=`^([0-9]{1,3}\.){3}[0-9]{1,3}$`
	// +kubebuilder:validation:Required
	SwitchIP string `json:"switchIP"`

	// SSHUser is the username used for SSH authentication
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	SSHUser string `json:"sshUser"`

	// SSHPassword is the password for SSH login (use Secret in production)
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	SSHPassword string `json:"sshPassword"`

	// Config is the switch configuration commands or script
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Config string `json:"config"`
}

// SwitchConfigStatus defines the observed state of SwitchConfig.
type SwitchConfigStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SwitchConfig is the Schema for the switchconfigs API
type SwitchConfig struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of SwitchConfig
	// +required
	Spec SwitchConfigSpec `json:"spec"`

	// status defines the observed state of SwitchConfig
	// +optional
	Status SwitchConfigStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// SwitchConfigList contains a list of SwitchConfig
type SwitchConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwitchConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SwitchConfig{}, &SwitchConfigList{})
}
