package testthings

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	// MinimumFilePerm is the minimum allowed permission bits that skeleton will
	// create during install.
	MinimumFilePerm os.FileMode = 0440
	// MinimumDirPerm is the minimum allowed permission bits that skeleton will
	// create during install.
	MinimumDirPerm os.FileMode = 0700
)

// SkeletonInstallError is returned when the skeleton cannot install the file
// for some path.
type SkeletonInstallError struct {
	Path string
	Err  error

	fileMode fs.FileMode
}

// Error implements error.
func (skel SkeletonInstallError) Error() string {
	return fmt.Sprintf("skel path %q: %v", skel.Path, skel.Err)
}

// Unwrap returns the wrapped error.
func (skel SkeletonInstallError) Unwrap() error {
	return skel.Err
}

var _ error = (*SkeletonInstallError)(nil)

// SkeletonFS is used to create new copies of the fsys in test case directories.
//
// Skeleton installation has no support for nodes other than directories and
// regular files.
func SkeletonFS(fsys fs.FS) skeletonFS {
	return skeletonFS{fsys}
}

type skeletonFS struct {
	skeleton fs.FS
}

// Install the skeleton into the provided directory.
func (skel skeletonFS) Install(dir string) error {
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("dir: %w", err)
	}

	return fs.WalkDir(skel.skeleton, ".",
		ignoreNodeTypeErrors(
			skeletonInstaller(skel.skeleton, dir)))
}

// InstallOrFail will install the skeleton into the provided directory. Or.. it
// fails the test run.
func (skel skeletonFS) InstallOrFail(testingT Terminator, dir string) {
	err := skel.Install(dir)
	if err != nil {
		testingT.Fatal(fmt.Sprintf("skeleton install: %v", err))
	}
}

func ignoreNodeTypeErrors(fn fs.WalkDirFunc) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, walkErr error) error {
		err := fn(path, d, walkErr)
		var ierr SkeletonInstallError
		if errors.As(err, &ierr) {
			if ierr.fileMode == 0 {
				return err
			}
			if !ierr.fileMode.IsRegular() &&
				!ierr.fileMode.IsDir() {
				// we don't operate on these files, so just supress the errors
				return nil
			}
		}
		return err
	}
}

func skeletonInstaller(skelFS fs.FS, dir string) fs.WalkDirFunc {
	installPath := func(p ...string) string {
		elms := append([]string{dir}, p...)
		return filepath.Join(elms...)
	}

	return func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return SkeletonInstallError{Path: path, Err: fmt.Errorf("walk: %w", walkErr)}
		}
		info, err := d.Info()
		if err != nil {
			return SkeletonInstallError{Path: path, Err: fmt.Errorf("info: %w", err)}
		}

		if info.IsDir() {
			err := os.MkdirAll(installPath(path), info.Mode()|MinimumDirPerm)
			if err != nil {
				return SkeletonInstallError{Path: path, Err: err, fileMode: info.Mode()}
			}
			return nil
		}

		// TODO: maybe handle symlinks?
		//
		// seems like it could be useful, but might breed complex test
		// filesystem setup (which isn't something I want this library to
		// encourage). That can be dealt with by callers who can prepare their
		// directories however they want (merging or modifying the target tree).
		if isSymlink(info.Mode()) {
			return SkeletonInstallError{Path: path, Err: errors.New("unsupported file type: symlink"), fileMode: info.Mode()}
		}

		// Catch the rest - deal only with regular files.
		if !info.Mode().IsRegular() {
			return SkeletonInstallError{Path: path, Err: fmt.Errorf("unsupported file type: %v", info.Mode()), fileMode: info.Mode()}
		}

		skelF, err := skelFS.Open(path)
		if err != nil {
			return SkeletonInstallError{Path: path, Err: fmt.Errorf("open skel file: %w", err), fileMode: info.Mode()}
		}
		defer skelF.Close()

		// NOTE: will fail on conflicts
		installF, err := os.OpenFile(installPath(path), os.O_WRONLY|os.O_EXCL|os.O_CREATE, info.Mode()|MinimumFilePerm)
		if err != nil {
			return SkeletonInstallError{Path: path, Err: fmt.Errorf("install file: %w", err)}
		}
		defer installF.Close()

		_, err = io.Copy(installF, skelF)
		if err != nil {
			return SkeletonInstallError{Path: path, Err: err}
		}
		return installF.Close()
	}
}

func isSymlink(info fs.FileMode) bool {
	return info&fs.ModeSymlink == fs.ModeSymlink
}
