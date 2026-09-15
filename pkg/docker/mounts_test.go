package docker

import (
	"testing"

	"github.com/docker/docker/api/types/mount"
)

func TestClassifyMount(t *testing.T) {
	cases := []struct {
		source    string
		mountType mount.Type
		mkdir     bool
	}{
		{"/data/app", mount.TypeBind, true},
		{"/", mount.TypeBind, true},
		{"mydata", mount.TypeVolume, false},
		{"app-data-1", mount.TypeVolume, false},
		{"", mount.TypeVolume, false},
		{"relative/path", mount.TypeVolume, false},
		{"./local", mount.TypeVolume, false},
		{"../escape", mount.TypeVolume, false},
	}
	for _, tc := range cases {
		mountType, mkdir := ClassifyMount(tc.source)
		if mountType != tc.mountType || mkdir != tc.mkdir {
			t.Fatalf("source %q: got (%v, %v), want (%v, %v)",
				tc.source, mountType, mkdir, tc.mountType, tc.mkdir)
		}
	}
}
