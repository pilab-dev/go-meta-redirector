package main

import "testing"

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
