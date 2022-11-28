//go:build windows

/*
Copyright 2020 The Kubernetes Authors.

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

package resizefs

import (
	"context"

	"github.com/kubernetes-csi/csi-proxy/v2/pkg/volume"

	"k8s.io/klog/v2"
	"k8s.io/mount-utils"
	mounter "sigs.k8s.io/gcp-compute-persistent-disk-csi-driver/pkg/mount-manager"
)

var _ Resizefs = &resizeFs{}

// ResizeFs Provides support for resizing file systems
type resizeFs struct {
	mounter *mount.SafeFormatAndMount
}

// NewResizeFs returns new instance of resizer
func NewResizeFs(mounter *mount.SafeFormatAndMount) *resizeFs {
	return &resizeFs{mounter: mounter}
}

// resize perform resize of file system
func (resizefs *resizeFs) Resize(devicePath string, deviceMountPath string) (bool, error) {
	klog.V(3).Infof("resizeFS.Resize - Expanding mounted volume %s", deviceMountPath)

	proxy := resizefs.mounter.Interface.(*mounter.CSIProxyMounterImpl)

	idRequest := &volume.GetVolumeIDFromTargetPathRequest{
		TargetPath: deviceMountPath,
	}
	idResponse, err := proxy.Volume.GetVolumeIDFromTargetPath(context.Background(), idRequest)
	if err != nil {
		return false, err
	}
	volumeID := idResponse.VolumeID

	request := &volume.ResizeVolumeRequest{
		VolumeID: volumeID,
	}
	_, err = proxy.Volume.ResizeVolume(context.Background(), request)
	if err != nil {
		return false, err
	}
	return true, nil
}
