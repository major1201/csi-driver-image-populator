/*
Copyright 2019 The Kubernetes Authors.

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

package main

import (
	"flag"
	"time"

	"github.com/major1201/csi-driver-image-populator/pkg/image"
)

func init() {
	flag.Set("logtostderr", "true")
}

var (
	endpoint    = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	driverName  = flag.String("drivername", "image.csi.k8s.io", "name of the driver")
	nodeID      = flag.String("nodeid", "", "node id")
	buildahPath = flag.String("buildah_path", "/bin/buildah", "buildah path, default: /bin/buildah")
	pullTimeout = flag.Duration("pull_timeout", 5*time.Minute, "image pull timeout, default: 5m")
)

func main() {
	flag.Parse()

	handle()
}

func handle() {
	driver := image.NewDriver(*driverName, *nodeID, *endpoint, *buildahPath, *pullTimeout)
	driver.Run()
}
