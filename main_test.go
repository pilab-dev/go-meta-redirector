package main

import (
	"net/http/httptest"
	"testing"
)

func TestMatchFallback(t *testing.T) {
	const pattern, target = "cloud/*", "https://github.com/pilab-dev/*"
	cases := []struct {
		req, git, root string
		ok             bool
	}{
		{"cloud/virtpb", "https://github.com/pilab-dev/virtpb.git", "cloud/virtpb", true},
		{"cloud/virtpb/v3", "https://github.com/pilab-dev/virtpb.git", "cloud/virtpb", true},
		{"cloud/virtpb/v3/pilab/common/v1", "https://github.com/pilab-dev/virtpb.git", "cloud/virtpb", true},
		{"cloud/", "", "", false},
		{"other/virtpb", "", "", false},
	}
	for _, c := range cases {
		git, root, ok := matchFallback(pattern, target, c.req)
		if git != c.git || root != c.root || ok != c.ok {
			t.Errorf("matchFallback(%q) = %q, %q, %v; want %q, %q, %v", c.req, git, root, ok, c.git, c.root, c.ok)
		}
	}
}

func TestHandlerAdvertisesModuleRoot(t *testing.T) {
	config = Config{Domains: map[string]DomainConfig{"go.pilab.hu": {
		Fallback: &FallbackConfig{Pattern: "cloud/*", Target: "https://github.com/pilab-dev/*"},
	}}}

	for _, path := range []string{"/cloud/virtpb", "/cloud/virtpb/v3", "/cloud/virtpb/v3/pilab/common/v1"} {
		req := httptest.NewRequest("GET", path+"?go-get=1", nil)
		req.Host = "go.pilab.hu"
		rec := httptest.NewRecorder()
		handler(rec, req)

		const want = `<html><head><meta name="go-import" content="go.pilab.hu/cloud/virtpb git https://github.com/pilab-dev/virtpb.git"></head></html>`
		if got := rec.Body.String(); got != want {
			t.Errorf("%s: got %s", path, got)
		}
	}
}
