package docker

import (
	"path/filepath"

	"github.com/docker/docker/api/types/mount"
)

// ClassifyMount decides how one volume source is attached. Absolute sources
// are host binds and the caller must ensure the directory exists; anything
// else is a daemon-managed named volume and must never be created on the
// host filesystem. Treating a name as a bind silently redirects daemon
// volume management into an arbitrary host directory.
func ClassifyMount(source string) (mount.Type, bool) {
	if filepath.IsAbs(source) {
		return mount.TypeBind, true
	}
	return mount.TypeVolume, false
}
