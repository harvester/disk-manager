/*
Copyright 2024 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/harvester/node-disk-manager/pkg/apis/harvesterhci.io/v1beta1"
	"github.com/rancher/wrangler/v3/pkg/generic"
)

// LVMVolumeGroupController interface for managing LVMVolumeGroup resources.
type LVMVolumeGroupController interface {
	generic.ControllerInterface[*v1beta1.LVMVolumeGroup, *v1beta1.LVMVolumeGroupList]
}

// LVMVolumeGroupClient interface for managing LVMVolumeGroup resources in Kubernetes.
type LVMVolumeGroupClient interface {
	generic.ClientInterface[*v1beta1.LVMVolumeGroup, *v1beta1.LVMVolumeGroupList]
}

// LVMVolumeGroupCache interface for retrieving LVMVolumeGroup resources in memory.
type LVMVolumeGroupCache interface {
	generic.CacheInterface[*v1beta1.LVMVolumeGroup]
}
