package docker

import (
	"net/http"
	"strings"
	"testing"

	"github.com/docker/distribution/manifest/ocischema"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

const ociManifestFixture = `{
	"schemaVersion": 2,
	"mediaType": "application/vnd.oci.image.manifest.v1+json",
	"config": {
		"mediaType": "application/vnd.oci.image.config.v1+json",
		"digest": "sha256:ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12",
		"size": 7023
	},
	"layers": [
		{
			"mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
			"digest": "sha256:ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12ab12",
			"size": 32654
		}
	]
}`

func TestParseOCIManifest(t *testing.T) {
	parsed, contentType, err := parseManifestPayload(v1.MediaTypeImageManifest, []byte(ociManifestFixture))
	if err != nil {
		t.Fatalf("OCI manifest rejected: %v", err)
	}
	if contentType != v1.MediaTypeImageManifest {
		t.Fatalf("content type = %q", contentType)
	}
	manifest, ok := parsed.(*ocischema.Manifest)
	if !ok {
		t.Fatalf("parsed type = %T, want *ocischema.Manifest", parsed)
	}
	if len(manifest.Layers) != 1 {
		t.Fatalf("layers = %d, want 1", len(manifest.Layers))
	}
	for _, digest := range append([]string{string(manifest.Config.Digest)}, manifest.Layers[0].Digest.String()) {
		trimmed := strings.TrimPrefix(digest, "sha256:")
		if len(trimmed) != 64 {
			t.Fatalf("digest %q is not a full 256-bit hex value", digest)
		}
		for i := 0; i < len(trimmed); i++ {
			character := trimmed[i]
			if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
				t.Fatalf("digest %q is not lowercase hex", digest)
			}
		}
	}
}

func TestParseManifestRejectsUnknown(t *testing.T) {
	_, _, err := parseManifestPayload("application/vnd.example.unknown", []byte(`{}`))
	if err == nil || !strings.Contains(err.Error(), "unknown content type") {
		t.Fatalf("unknown type must fail closed, got %v", err)
	}
	_, _, err = parseManifestPayload(v1.MediaTypeImageManifest, []byte(`{invalid`))
	if err == nil {
		t.Fatalf("malformed payload must fail")
	}
}

func TestRegistryClientVerifiesTLS(t *testing.T) {
	client := httpClient()
	transport, ok := client.Transport.(*http.Transport)
	if !ok || transport == nil {
		t.Fatalf("unexpected transport type %T", client.Transport)
	}
	if config := transport.TLSClientConfig; config != nil && config.InsecureSkipVerify {
		t.Fatalf("registry client must verify TLS")
	}
}
