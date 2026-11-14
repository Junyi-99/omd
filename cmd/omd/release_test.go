package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tarGz(t *testing.T, name, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0755, Size: int64(len(content)), Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(tw, content); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// releaseServer serves one asset and API metadata recording digest for it.
func releaseServer(t *testing.T, archive []byte, digest string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/repos/o/tool/releases/tags/v1":
			fmt.Fprintf(w, `{"assets":[{"name":"tool-amd64.tar.gz","digest":%q}]}`, digest)
		case r.URL.Path == "/o/tool/releases/download/v1/tool-amd64.tar.gz":
			w.Write(archive)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestInstallReleaseVerifiesDigest(t *testing.T) {
	archive := tarGz(t, "tool", "#!/bin/sh\n")
	sum := sha256.Sum256(archive)
	good := "sha256:" + hex.EncodeToString(sum[:])
	rel := &Release{Repo: "o/tool", Asset: "tool-{arch}.tar.gz", Arch: map[string]string{"amd64": "amd64"}, Bins: []string{"tool"}}
	p := Package{ID: "tool", Name: "Tool"}

	for _, tc := range []struct {
		name, digest, wantErr string
	}{
		{"match", good, ""},
		{"mismatch", "sha256:" + strings.Repeat("0", 64), "checksum mismatch"},
		{"no digest", "", "no SHA-256"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := releaseServer(t, archive, tc.digest)
			m := testManager(t, "linux", nil)
			m.github, m.api = srv.URL, srv.URL
			for _, dir := range []string{"bin", "dep"} {
				if err := os.MkdirAll(filepath.Join(m.root, dir), 0755); err != nil {
					t.Fatal(err)
				}
			}
			err := m.installRelease(p, rel, "v1", io.Discard)
			_, statErr := os.Lstat(filepath.Join(m.root, "bin", "tool"))
			if tc.wantErr == "" {
				if err != nil || statErr != nil {
					t.Fatalf("install: %v, link: %v", err, statErr)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("err = %v, want %q", err, tc.wantErr)
			}
			if statErr == nil {
				t.Fatal("an unverified release was linked into bin/")
			}
		})
	}
}

func TestLinkLegacyXDG(t *testing.T) {
	m := testManager(t, "linux", nil)
	config := filepath.Join(m.home, ".config")
	for _, d := range []string{"gh", "nvim"} {
		if err := os.MkdirAll(filepath.Join(config, d), 0755); err != nil {
			t.Fatal(err)
		}
	}
	// The repo tracks its own nvim config: that one wins and is reported.
	if err := os.MkdirAll(filepath.Join(m.root, "etc", "nvim"), 0755); err != nil {
		t.Fatal(err)
	}

	linked, conflicts, err := m.linkLegacyXDG()
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := os.Readlink(filepath.Join(m.root, "etc", "gh")); got != filepath.Join(config, "gh") {
		t.Fatalf("etc/gh -> %q", got)
	}
	if len(linked) != 1 || len(conflicts) != 1 || !strings.Contains(conflicts[0], "nvim") {
		t.Fatalf("linked %v, conflicts %v", linked, conflicts)
	}

	// Idempotent: a second run links nothing new.
	linked, _, err = m.linkLegacyXDG()
	if err != nil || len(linked) != 0 {
		t.Fatalf("second run: linked %v, err %v", linked, err)
	}
}
