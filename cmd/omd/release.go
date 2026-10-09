package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Release installs a prebuilt archive from a project's GitHub Releases.
// Asset may use {tag}, {version} (tag without a leading v) and {arch}, where
// {arch} is the project's own name for runtime.GOARCH.
type Release struct {
	Repo  string
	Asset string
	Arch  map[string]string
	// Executables inside the unpacked archive; each is linked into bin/.
	Bins []string
}

const versionFile = ".omd-version"

var errUnsupportedArch = errors.New("no prebuilt release for this CPU architecture")

func (r Release) assetName(tag, goarch string) (string, error) {
	arch, ok := r.Arch[goarch]
	if !ok {
		return "", fmt.Errorf("%w (%s)", errUnsupportedArch, goarch)
	}
	return strings.NewReplacer(
		"{tag}", tag,
		"{version}", strings.TrimPrefix(tag, "v"),
		"{arch}", arch,
	).Replace(r.Asset), nil
}

func (m *manager) releaseFor(p Package) *Release {
	return p.installerFor(m.platform).Release
}

func (m *manager) releaseDir(p Package) string { return filepath.Join(m.root, "dep", p.ID) }

// installedRelease reports the version omd installed. A binary left in bin/ by
// an older shell installer has no recorded version and reports "unknown".
func (m *manager) installedRelease(p Package) (version string, managed bool) {
	rel := m.releaseFor(p)
	if rel == nil {
		return "", false
	}
	if data, err := os.ReadFile(filepath.Join(m.releaseDir(p), versionFile)); err == nil {
		return strings.TrimSpace(string(data)), true
	}
	for _, bin := range rel.Bins {
		info, err := os.Lstat(filepath.Join(m.root, "bin", path.Base(bin)))
		if err == nil && info.Mode().IsRegular() {
			return "unknown", true
		}
	}
	return "", false
}

// latestRelease reads the tag from GitHub's /releases/latest redirect, which,
// unlike the REST API, is not rate limited for anonymous clients.
func (m *manager) latestRelease(ctx context.Context, rel *Release) (string, error) {
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, m.github+"/"+rel.Repo+"/releases/latest", nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	location := resp.Header.Get("Location")
	_, tag, ok := strings.Cut(location, "/releases/tag/")
	if !ok || tag == "" {
		return "", fmt.Errorf("%s: no latest release (HTTP %s)", rel.Repo, resp.Status)
	}
	return url.PathUnescape(tag)
}

// checkRelease returns the installed and latest versions of an omd-managed release.
func (m *manager) checkRelease(ctx context.Context, p Package) (installed, latest string, err error) {
	installed, managed := m.installedRelease(p)
	if !managed {
		return "", "", nil
	}
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	latest, err = m.latestRelease(ctx, m.releaseFor(p))
	return installed, latest, err
}

func (m *manager) installRelease(p Package, rel *Release, tag string, out io.Writer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	if tag == "" {
		var err error
		if tag, err = m.latestRelease(ctx, rel); err != nil {
			return err
		}
	}
	asset, err := rel.assetName(tag, m.arch)
	if err != nil {
		return err
	}
	dep := filepath.Join(m.root, "dep")
	staging, err := os.MkdirTemp(dep, ".omd-"+p.ID+"-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)

	link := m.github + "/" + rel.Repo + "/releases/download/" + url.PathEscape(tag) + "/" + url.PathEscape(asset)
	fmt.Fprintf(out, "Downloading %s\n", link)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: HTTP %s", asset, resp.Status)
	}
	if err := extractTarGz(resp.Body, staging); err != nil {
		return fmt.Errorf("unpack %s: %w", asset, err)
	}
	unpacked, err := stripSingleDir(staging)
	if err != nil {
		return err
	}
	for _, bin := range rel.Bins {
		if info, err := os.Stat(filepath.Join(unpacked, bin)); err != nil || info.Mode()&0111 == 0 {
			return fmt.Errorf("%s does not contain executable %s", asset, bin)
		}
	}
	if err := os.WriteFile(filepath.Join(unpacked, versionFile), []byte(tag+"\n"), 0644); err != nil {
		return err
	}

	// Swap the new tree in, keeping the old one until the rename succeeds.
	target := m.releaseDir(p)
	old := staging + ".old"
	if err := os.Rename(target, old); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.Rename(unpacked, target); err != nil {
		if restoreErr := os.Rename(old, target); restoreErr != nil && !errors.Is(restoreErr, os.ErrNotExist) {
			err = errors.Join(err, restoreErr)
		}
		return err
	}
	os.RemoveAll(old)

	for _, bin := range rel.Bins {
		link := filepath.Join(m.root, "bin", path.Base(bin))
		if err := os.Remove(link); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := os.Symlink(filepath.Join(target, bin), link); err != nil {
			return err
		}
	}
	fmt.Fprintf(out, "Installed %s %s (%s)\n", p.Name, tag, asset)
	return nil
}

func extractTarGz(r io.Reader, dir string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		name := path.Clean(strings.TrimPrefix(hdr.Name, "./"))
		if name == "." {
			continue
		}
		if !filepath.IsLocal(name) {
			return fmt.Errorf("unsafe path in archive: %s", hdr.Name)
		}
		dest := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(dest, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, hdr.FileInfo().Mode().Perm()|0600)
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(f, tr)
			if err := errors.Join(copyErr, f.Close()); err != nil {
				return err
			}
		case tar.TypeSymlink:
			resolved := path.Join(path.Dir(name), hdr.Linkname)
			if path.IsAbs(hdr.Linkname) || !filepath.IsLocal(resolved) {
				return fmt.Errorf("unsafe symlink in archive: %s -> %s", hdr.Name, hdr.Linkname)
			}
			if err := os.Symlink(hdr.Linkname, dest); err != nil {
				return err
			}
		}
	}
}

// stripSingleDir returns the archive's only top-level directory, if it has one.
func stripSingleDir(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	if len(entries) == 1 && entries[0].IsDir() {
		return filepath.Join(dir, entries[0].Name()), nil
	}
	return dir, nil
}
