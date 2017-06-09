# openebsdriver
![Go Report Card](https://goreportcard.com/badge/github.com/maheshreddy7797/openebsdriver)

## VolumeDriver for Openebs

The VolumeDriver capability basically gives plugins control over the volumes life cycle. A plugin registers itself as a VolumeDriver plugin and when the host requires a volume with a specific name for that Driver. The plugin provides a Mountpoint for that volume on the host machine.

VolumeDriver plugins can be used for things like distributed filesystems and stateful volumes.

![A volume response process ](https://github.com/sripadaraj/openebsdriver/blob/master/images/Chart_Docker-Volume-Plugin-Architecture.png)


OpenEBS Plugin 
======================================

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

## Examples

### Create a volume

```
$ docker volume create --driver openebsdriver --name <volumename>
# Set user and group of the volume
$ docker volume create --driver openebsdriver --name <volumename> --opt user=docker --opt group=docker
```

### Delete a volume

```
$ docker volume rm <volumename>
```

### List all volumes

```
$ docker volume ls
```

### Attach volume to container

```
$ docker run --volume-driver=openebsdriver -v <volumename>:/vol busybox sh -c 'echo "Hello World" > /vol/hello.txt'
```

## Development

### Build

Get the code

```
$ go get -u github.com/openebsdriver/docker-volume
```

#### Dependency Management

For the dependency management we use [golang dep](https://github.com/golang/dep)

#### Linux

```
$ go build -ldflags "-s -w" -o bin/docker-openebsdriver-plugin .
```

#### OSX/MacOS

```
$ GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/docker-openebsdriver-plugin .
```

#### Docker

```
$ docker run --rm -v "$GOPATH":/work -e "GOPATH=/work" -w /work/src/github.com/quobyte/docker-volume golang:1.8 go build -v -ldflags "-s -w" -o bin/openebsdriver-docker-plugin
```

