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
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
)

var (
	version = "1.0.0"
)

type driver struct {
	*csi.UnimplementedIdentityServer
	*csi.UnimplementedGroupControllerServer
	*csi.UnimplementedNodeServer

	config Config
}

type Config struct {
	DriverName    string
	Endpoint      string
	NodeID        string
	VendorVersion string
	BuildahPath   string
	ExecTimeout   time.Duration
}

func NewDriver(driverName, nodeID, endpoint, buildahPath string, pullTimeout time.Duration) *driver {
	glog.Infof("Driver: %v version: %v", driverName, version)

	d := &driver{}
	d.config = Config{
		DriverName:    driverName,
		Endpoint:      endpoint,
		NodeID:        nodeID,
		VendorVersion: version,
		BuildahPath:   buildahPath,
		ExecTimeout:   pullTimeout,
	}
	return d
}

func (d *driver) Run() {
	s := NewNonBlockingGRPCServer()
	s.Start(d.config.Endpoint, d, d, d, d)
	s.Wait()
}
