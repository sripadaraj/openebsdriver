package openebs

import (
	"github.com/docker/docker/volume"
	volumeplugin "github.com/docker/go-plugins-helpers/volume"
)

type openebsDriver struct {
	d volume.Driver
}

// NewHandlerFromVolumeDriver creates a plugin handler from an existing volume
// driver. This could be used, for instance, by the `local` volume driver built-in
// to Docker Engine and it would create a plugin from it that maps plugin API calls
// directly to any volume driver that satifies the volume.Driver interface from
// Docker Engine.
func NewHandlerFromVolumeDriver(d volume.Driver) *volumeplugin.Handler {
	return volumeplugin.NewHandler(&openebsDriver{d})
}

// Create request results in API call to Openebs registry 
// that creates requested volume if possible.
func (d *openebsDriver) Create(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	_, err := d.d.Create(req.Name, req.Options)
	if err != nil {
		res.Err = err.Error()
	}
	return res
}

// List all volumes and their respective mount points.
func (d *openebsDriver) List(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	ls, err := d.d.List()
	if err != nil {
		res.Err = err.Error()
		return res
	}
	vols := make([]*volumeplugin.Volume, len(ls))

	for i, v := range ls {
		vol := &volumeplugin.Volume{
			Name:       v.Name(),
			Mountpoint: v.Path(),
		}
		vols[i] = vol
	}
	res.Volumes = vols
	return res
}

// Get request the information about specified volume
// and returns the name, mountpoint & status.
func (d *openebsDriver) Get(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	res.Volume = &volumeplugin.Volume{
		Name:       v.Name(),
		Mountpoint: v.Path(),
		Status:     v.Status(),
	}
	return res
}

// Remove is called to delete a volume.
func (d *openebsDriver) Remove(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	if err := d.d.Remove(v); err != nil {
		res.Err = err.Error()
	}
	return res
}

// Path is called to get the path of a volume mounted on the host.
func (d *openebsDriver) Path(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	res.Mountpoint = v.Path()
	return res
}

// Mount bind the volume to a container specified by the Path.
func (d *openebsDriver) Mount(req volumeplugin.MountRequest) volumeplugin.Response {
	var res volumeplugin.Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	pth, err := v.Mount(req.ID)
	if err != nil {
		res.Err = err.Error()
	}
	res.Mountpoint = pth
	return res
}

// Unmount is called to stop the container from using the volume
// and it is probably safe to unmount.
func (d *openebsDriver) Unmount(req volumeplugin.UnmountRequest) volumeplugin.Response {
	var res volumeplugin.Response
	v, err := d.d.Get(req.Name)
	if err != nil {
		res.Err = err.Error()
		return res
	}
	if err := v.Unmount(req.ID); err != nil {
		res.Err = err.Error()
	}
	return res
}

// Capabilities indicates if a volume has to be created multiple times or only once.
func (d *openebsDriver) Capabilities(req volumeplugin.Request) volumeplugin.Response {
	var res volumeplugin.Response
	res.Capabilities = volumeplugin.Capability{Scope: d.d.Scope()}
	return res
}
