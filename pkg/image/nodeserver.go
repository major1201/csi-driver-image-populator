/*
Copyright 2017 The Kubernetes Authors.

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

package image

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"github.com/major1201/csi-driver-image-populator/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/mount-utils"
)

const (
	volumeContextImage     = "image"
	volumeContextChangeDir = "changeDir"
)

var (
	ErrTimeout = fmt.Errorf("Timeout")
)

func (d *driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (resp *csi.NodePublishVolumeResponse, err error) {
	// Check arguments
	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume capability missing in request")
	}
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}

	image := req.GetVolumeContext()[volumeContextImage]
	changeDir := filepath.Clean(req.GetVolumeContext()[volumeContextChangeDir])
	if strings.HasPrefix(changeDir, "..") {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Invaid change_dir found: %s", changeDir))
	}

	err = d.setupVolume(req.GetVolumeId(), image)
	if err != nil {
		return nil, err
	}
	defer func() {
		// do cleanup if failed on the following steps
		if err != nil {
			glog.Infof("cleanup volume: %s", req.GetVolumeId())
			_ = d.unsetupVolume(req.GetVolumeId())
		}
	}()

	targetPath := req.GetTargetPath()
	notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(targetPath, 0750); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			notMnt = true
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if !notMnt {
		return &csi.NodePublishVolumeResponse{}, nil
	}

	fsType := req.GetVolumeCapability().GetMount().GetFsType()
	readOnly := req.GetReadonly()
	volumeId := req.GetVolumeId()
	attrib := req.GetVolumeContext()
	mountFlags := req.GetVolumeCapability().GetMount().GetMountFlags()

	glog.V(4).Infof("target=%v, fstype=%v, readonly=%v, volumeId=%v, attributes=%v,  mountflags=%v",
		targetPath, fsType, readOnly, volumeId, attrib, mountFlags)

	options := []string{"bind"}
	if readOnly {
		options = append(options, "ro")
	}

	args := []string{"mount", volumeId}
	output, err := d.runCmd(args)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to mount, err=%s", err.Error()))
	}
	provisionRoot := strings.TrimSpace(string(output[:]))
	glog.V(4).Infof("container mount point at %s\n", provisionRoot)

	mounter := mount.New("")
	path := filepath.Join(provisionRoot, changeDir)
	if !utils.IsDir(path) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Dir not found in image fs: %s", changeDir))
	}

	if err := mounter.Mount(path, targetPath, "", options); err != nil {
		return nil, err
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (d *driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	// Check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	targetPath := req.GetTargetPath()
	volumeId := req.GetVolumeId()

	// Unmounting the image
	err := mount.New("").Unmount(req.GetTargetPath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.V(4).Infof("image: volume %s/%s has been unmounted.", targetPath, volumeId)

	err = d.unsetupVolume(volumeId)
	if err != nil {
		return nil, err
	}
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (d *driver) setupVolume(volumeId string, image string) error {
	args := []string{"from", "--name", volumeId, "--pull", image}
	glog.Infof("pulling image for volume=%s: %s", volumeId, image)

	output, err := d.runCmd(args)
	provisionRoot := strings.TrimSpace(string(output[:]))
	glog.V(4).Infof("setup: container mount point at %s\n", provisionRoot)
	return err
}

func (d *driver) unsetupVolume(volumeId string) error {
	args := []string{"delete", volumeId}
	output, err := d.runCmd(args)
	provisionRoot := strings.TrimSpace(string(output[:]))
	glog.V(4).Infof("unsetup: container mount point at %s\n", provisionRoot)
	return err
}

func (d *driver) runCmd(args []string) ([]byte, error) {
	cmd := exec.Command(d.config.BuildahPath, args...)

	output, killed, err := utils.SafeExecWithCombinedOutput(cmd, d.config.ExecTimeout)
	if killed {
		err = ErrTimeout
	}
	return output, err
}

func (d *driver) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (d *driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return &csi.NodeStageVolumeResponse{}, nil
}

func (d *driver) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{}, nil
}

func (d *driver) NodeGetInfo(context.Context, *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId: d.config.NodeID,
	}, nil
}
