/*
Copyright 2026.

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

package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/qdrant/go-client/qdrant"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	// CollectionSpec defines the desired state of Collection
	CollectionSpec struct {
		// +kubebuilder:validation:Minimum=1
		VectorSize int `json:"vectorSize"`
		// +kubebuilder:validation:Enum=Cosine;Euclid;Dot;Manhattan
		Distance string `json:"distance"` // Cosine, Euclid, Dot, Manhattan
	}

	// CollectionStatus defines the observed state of Collection.
	CollectionStatus struct {
		Conditions []metav1.Condition `json:"conditions,omitempty"`
		Status     string             `json:"status,omitempty"` // Green, Yellow, Red
	}

	// +kubebuilder:object:root=true
	// +kubebuilder:subresource:status

	// Collection is the Schema for the collections API
	Collection struct {
		metav1.TypeMeta `json:",inline"`

		// metadata is a standard object metadata
		// +optional
		metav1.ObjectMeta `json:"metadata,omitzero"`

		// spec defines the desired state of Collection
		// +required
		Spec CollectionSpec `json:"spec"`

		// status defines the observed state of Collection
		// +optional
		Status CollectionStatus `json:"status,omitzero"`
	}

	// +kubebuilder:object:root=true

	// CollectionList contains a list of Collection
	CollectionList struct {
		metav1.TypeMeta `json:",inline"`
		metav1.ListMeta `json:"metadata,omitzero"`
		Items           []Collection `json:"items"`
	}
)

func init() {
	SchemeBuilder.Register(&Collection{}, &CollectionList{})
}

func (s *CollectionSpec) ToQdrantParams(collectionName string) (*qdrant.CreateCollection, error) {
	distance, err := parseDistance(s.Distance)
	if err != nil {
		return nil, err
	}

	return &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(s.VectorSize),
			Distance: distance,
		}),
	}, nil
}

func parseDistance(d string) (qdrant.Distance, error) {
	switch strings.ToLower(d) {
	case "cosine":
		return qdrant.Distance_Cosine, nil
	case "euclid":
		return qdrant.Distance_Euclid, nil
	case "dot":
		return qdrant.Distance_Dot, nil
	case "manhattan":
		return qdrant.Distance_Manhattan, nil
	default:
		return 0, fmt.Errorf("unknown distance %q, must be one of: Cosine, Euclid, Dot, Manhattan", d)
	}
}
