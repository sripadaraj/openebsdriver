# OpenEBS Volume Plugin

![Travis Ci](https://travis-ci.org/maheshreddy7797/openebsdriver.svg?branch=master) ![Go Report Card](https://goreportcard.com/badge/github.com/maheshreddy7797/openebsdriver)![Go Report](https://codecov.io/gh/maheshreddy7797/openebsdriver)

## VolumeDriver for Openebs

The VolumeDriver capability basically gives plugins control over the volumes life cycle. A plugin registers itself as a VolumeDriver plugin and when the host requires a volume with a specific name for that Driver. The plugin provides a Mountpoint for that volume on the host machine.

VolumeDriver plugins can be used for things like distributed filesystems and stateful volumes.





 
======
Usage:
1) Clone this repository
```
git clone https://github.com/maheshreddy7797/openebsdriver && cd openebsdriver
```
2) Copy nvd.json.example to /etc/nvd/nvd.json and change values according to your NexentaStor setup
```
mkdir /etc/openebsdriver
cp config.json /etc/openebsdriver/config.json
```

