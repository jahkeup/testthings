package skeletonfs_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/jahkeup/testthings/skeletonfs"
)

func TestSkeletonFS(t *testing.T) {
	t.Run("readonly files", func(t *testing.T) {
		skel := &fstest.MapFS{
			"foo/bar.readonly": &fstest.MapFile{
				Data: []byte{},
				Mode: 0444,
			},
		}

		installDir := t.TempDir()
		err := skeletonfs.SkeletonFS(skel).Install(installDir)
		assert.NoError(t, err)
		dumpFSPaths(t, os.DirFS(installDir))
		assert.NoError(t, os.RemoveAll(installDir), "should be able to remove the tree")
	})
	t.Run("happy path", func(t *testing.T) {
		skel := fstest.MapFS{
			"foo/etc/baz.conf": &fstest.MapFile{
				Data: []byte(`some data`),
				Mode: 0, // should be minimumFilePerm
			},
			"foo/share/baz.bin": &fstest.MapFile{
				Data: []byte(`some data`),
				Mode: 0444,
			},
			"foo/bin/baz.sh": &fstest.MapFile{
				Data: []byte(`#!/usr/bin/env sh\nexit 0\n`),
				Mode: 0750,
			},
		}

		installDir := t.TempDir()
		err := skeletonfs.SkeletonFS(skel).Install(installDir)
		assert.NoError(t, err)

		fsys := os.DirFS(installDir)
		dumpFSPaths(t, fsys)

		var paths []string
		for path, spec := range skel {
			paths = append(paths, path)
			info, err := fs.Stat(fsys, path)
			if assert.NoError(t, err, "should be able to stat file") {
				actualPerm := info.Mode().Perm()
				assert.True(t, actualPerm&skeletonfs.MinimumFilePerm == skeletonfs.MinimumFilePerm)

				adjustedPerm := spec.Mode.Perm() | skeletonfs.MinimumFilePerm
				assert.Equalf(t, adjustedPerm.String(), actualPerm.String(), "for path: %q", path)
			}
		}

		// test through standard fsys access
		assert.NoError(t, fstest.TestFS(fsys, paths...), "should have skeleton contents")
		assert.NoError(t, os.RemoveAll(installDir), "should be able to remove the tree")
	})

}

func dumpFSPaths(t testing.TB, fsys fs.FS) {
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk: %w", walkErr)
		}
		p := path
		if d.IsDir() {
			p = p + "/"
		}
		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("file info: %w", err)
		}

		t.Logf("path: %q\t(%v, %d bytes)", p, info.Mode(), info.Size())

		return nil
	})

	assert.NoError(t, err, "should be able to walk tree")
}
