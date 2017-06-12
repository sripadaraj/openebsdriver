package main

import (
	"flag"
	"log"
	"os"

	"github.com/docker/go-plugins-helpers/volume"
)

const openEBSID string = "openEBS"

var (
	version  string
	revision string
)

func main() {
	openEBSMountPath := flag.String("path", "/run/docker/openEBS/mnt", "Path where OpenEBS Volume is mounted on the host")
	openEBSMountOptions := flag.String("options", "-o user_xattr", "Fuse options to be used when OPenEBS is mounted")

	openEBSUser := flag.String("user", "root", "User to connect to the OpenEBS-Maya API server")
	openEBSPassword := flag.String("password", "openEBS", "Password for the user to connect to the OpenEBS-Maya API server")
	openEBSAPIURL := flag.String("api", "http://localhost:7860", "URL to the API server(s) in the form http(s)://host[:port][,host:port] or SRV record name")
	openEBSRegistry := flag.String("registry", "localhost:7861", "URL to the registry server(s) in the form of host[:port][,host:port] or SRV record name")

	group := flag.String("group", "root", "Group to create the unix socket")
	maxWaitTime := flag.Float64("max-wait-time", 30, "Maximimum wait time for filesystem checks to complete when a Volume is created before returning an error")
	maxFSChecks := flag.Int("max-fs-checks", 5, "Maximimum number of filesystem checks when a Volume is created before returning an error")
	showVersion := flag.Bool("version", false, "Shows version string")
	flag.Parse()

	if *showVersion {
		log.Printf("Version: %s - Revision: %s\n", version, revision)
		return
	}

	if err := validateAPIURL(*openEBSAPIURL); err != nil {
		log.Fatalln(err)
	}

	if err := os.MkdirAll(*openEBSMountPath, 0555); err != nil {
		log.Println(err.Error())
	}

	if !isMounted(*openEBSMountPath) {
		log.Printf("Mounting OpenEBS namespace in %s", *openEBSMountPath)
		mountAll(*openEBSMountOptions, *openEBSRegistry, *openEBSMountPath)
	}

	nDriver := newOpenEBSDriver(*openEBSAPIURL, *openEBSUser, *openEBSPassword, *openEBSMountPath, *maxFSChecks, *maxWaitTime)
	handler := volume.NewHandler(nDriver)

	log.Println(handler.ServeUnix(*group, openEBSID))
}
