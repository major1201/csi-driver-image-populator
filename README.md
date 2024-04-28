# csi-driver-image-populator

CSI driver that uses a container image as a volume.

## How it works:

Currently the driver makes use of buildah to download the container image if it is not already available, launch a new instance of it named after the volumeHandle, and mount it.

In the future, integration with CRI would be desirable so the driver could ask via CRI that the Container Runtime perform these activities in a generic way.

## Usage:

**This is a prototype driver. Do not use for production**

It also requires features that are still in development.

### Build imageplugin

```
$ make container
```

### Installing into Kubernetes

```
deploy/deploy-image.sh
```

### Example Usage in Kubernetes

```
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.13-alpine
    ports:
    - containerPort: 80
    volumeMount:
    - name: data
      mountPath: /usr/share/nginx/html
  volumes:
  - name: data
    csi:
      driver: image.csi.k8s.io
      volumeAttributes:
        image: alpine
        changeDir: /bin
```

### Start Image driver manually

```
$ sudo ./bin/imageplugin --endpoint /tmp/csi.sock --nodeid CSINode -v=5
```
