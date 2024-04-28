# CSI image populator

## Usage:

### Build plugin
```
$ make
```

### Start driver
```
$ sudo ./bin/imagepopulatorplugin --endpoint /tmp/csi.sock --nodeid CSINode -v=5
```

### Test using csc

Get ```csc``` tool from https://github.com/rexray/gocsi/tree/master/csc

#### Get plugin info

```
$ csc identity plugin-info --endpoint /tmp/csi.sock
"image.csi.k8s.io"  "0.1.0"
```

#### (UNIMPLEMENTED) Create a volume

```
$ csc controller new --endpoint /tmp/csi.sock --cap 1,block CSIVolumeName
CSIVolumeID
```

#### (UNIMPLEMENTED) Delete a volume

```
$ csc controller del --endpoint /tmp/csi.sock CSIVolumeID
CSIVolumeID
```

#### (UNIMPLEMENTED) Validate volume capabilities

```
$ csc controller validate-volume-capabilities --endpoint /tmp/csi.sock --cap 1,block CSIVolumeID
CSIVolumeID  true
```

#### NodePublish a volume

```
$ csc node publish --endpoint /tmp/csi.sock --cap 1,1 --target-path /mnt/mypath --vol-context image=alpine,changeDir=/bin CSIVolumeID
CSIVolumeID
```

#### NodeUnpublish a volume

```
$ csc node unpublish --endpoint /tmp/csi.sock --target-path /mnt/mypath CSIVolumeID
CSIVolumeID
```

#### Get NodeInfo

```
$ csc node get-info --endpoint /tmp/csi.sock
CSINode
```
