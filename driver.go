package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	//"time"

	"github.com/docker/go-plugins-helpers/volume"
	maya_api "github.com/openebs/maya"
)

type openEBSDriver struct {
	client       *maya_api.OpenEBSClient
	openEBSMount string
	m            *sync.Mutex
	maxFSChecks  int
	maxWaitTime  float64
}

// newOpenEBSDriver creates a plugin handler from an existing volume
// driver. This could be used, for instance, by the `local` volume driver built-in
// to Docker Engine and it would create a plugin from it that maps plugin API calls
// directly to any volume driver that satisfies the volume.Driver interface from
// Docker Engine.
func newOpenEBSDriver(apiURL string, username string, password string, openEBSMount string, maxFSChecks int, maxWaitTime float64) openEBSDriver {
	driver := openEBSDriver{
		client:       maya_api.NewOpenEBSClient(apiURL, username, password),
		openEBSMount: openEBSMount,
		m:            &sync.Mutex{},
		maxFSChecks:  maxFSChecks,
		maxWaitTime:  maxWaitTime,
	}

	return driver
}

// Create request results in API call to Openebs registry
// that creates requested volume if possible.
func (driver openEBSDriver) Create(request volume.Request) volume.Response {
	log.Printf("Creating volume %s\n", request.Name)
	driver.m.Lock()
	defer driver.m.Unlock()

	user, group := "root", "root"

	if usr, ok := request.Options["user"]; ok {
		user = usr
	}

	if grp, ok := request.Options["group"]; ok {
		group = grp
	}

	if _, err := driver.client.CreateVolume(&maya_api.CreateVolumeRequest{
		Name:        request.Name,
		RootUserID:  user,
		RootGroupID: group,
	}); err != nil {
		log.Println(err)

		if !strings.Contains(err.Error(), "ENTITY_EXISTS_ALREADY/POSIX_ERROR_NONE") {
			return volume.Response{Err: err.Error()}
		}
	}

	mPoint := filepath.Join(driver.openEBSMount, request.Name)
	log.Printf("Validate mounting volume %s on %s\n", request.Name, mPoint)
	if err := driver.checkMountPoint(mPoint); err != nil {
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{Err: ""}
}

// List all volumes and their respective mount points.

func (driver openEBSDriver) List(request volume.Request) volume.Response {
	driver.m.Lock()
	defer driver.m.Unlock()

	var vols []*volume.Volume
	files, err := ioutil.ReadDir(driver.openEBSMount)
	if err != nil {
		log.Println(err)
		return volume.Response{Err: err.Error()}
	}

	for _, entry := range files {
		if entry.IsDir() {
			vols = append(vols, &volume.Volume{Name: entry.Name(), Mountpoint: filepath.Join(driver.openEBSMount, entry.Name())})
		}
	}

	return volume.Response{Volumes: vols}
}
// Get request the information about specified volume
// and returns the name, mountpoint & status.
func (driver openEBSDriver) Get(request volume.Request) volume.Response {
	driver.m.Lock()
	defer driver.m.Unlock()

	mPoint := filepath.Join(driver.openEBSMount, request.Name)

	if fi, err := os.Lstat(mPoint); err != nil || !fi.IsDir() {
		log.Println(err)
		return volume.Response{Err: fmt.Sprintf("%v not mounted", mPoint)}
	}

	return volume.Response{Volume: &volume.Volume{Name: request.Name, Mountpoint: mPoint}}
}

// Remove is called to delete a volume.
func (driver openEBSDriver) Remove(request volume.Request) volume.Response {
	log.Printf("Removing volume %s\n", request.Name)
	driver.m.Lock()
	defer driver.m.Unlock()

	if err := driver.client.DeleteVolumeByName(request.Name, ""); err != nil {
		log.Println(err)
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{Err: ""}
}

// Path is called to get the path of a volume mounted on the host.

func (driver openEBSDriver) Path(request volume.Request) volume.Response {
	return volume.Response{Mountpoint: filepath.Join(driver.openEBSMount, request.Name)}
}

// Mount bind the volume to a container specified by the Path.
func (driver openEBSDriver) Mount(request volume.MountRequest) volume.Response {
	driver.m.Lock()
	defer driver.m.Unlock()
	mPoint := filepath.Join(driver.openEBSMount, request.Name)
	log.Printf("Mounting volume %s on %s\n", request.Name, mPoint)
	return volume.Response{Err: "", Mountpoint: mPoint}
}

// Unmount is called to stop the container from using the volume
// and it is probably safe to unmount.
func (driver openEBSDriver) Unmount(request volume.UnmountRequest) volume.Response {
	return volume.Response{}
}

// Capabilities indicates if a volume has to be created multiple times or only once.
func (driver openEBSDriver) Capabilities(request volume.Request) volume.Response {
	return volume.Response{Capabilities: volume.Capability{Scope: "global"}}
}
