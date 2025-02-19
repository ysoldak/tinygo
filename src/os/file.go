// Portions copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package os implements a subset of the Go "os" package. See
// https://godoc.org/os for details.
//
// Note that the current implementation is blocking. This limitation should be
// removed in a future version.
package os

import (
	"errors"
	"io"
	"syscall"
)

// Seek whence values.
//
// Deprecated: Use io.SeekStart, io.SeekCurrent, and io.SeekEnd.
const (
	SEEK_SET int = io.SeekStart
	SEEK_CUR int = io.SeekCurrent
	SEEK_END int = io.SeekEnd
)

// Mkdir creates a directory. If the operation fails, it will return an error of
// type *PathError.
func Mkdir(path string, perm FileMode) error {
	fs, suffix := findMount(path)
	if fs == nil {
		return &PathError{"mkdir", path, ErrNotExist}
	}
	err := fs.Mkdir(suffix, perm)
	if err != nil {
		return &PathError{"mkdir", path, err}
	}
	return nil
}

// Many functions in package syscall return a count of -1 instead of 0.
// Using fixCount(call()) instead of call() corrects the count.
func fixCount(n int, err error) (int, error) {
	if n < 0 {
		n = 0
	}
	return n, err
}

// Remove removes a file or (empty) directory. If the operation fails, it will
// return an error of type *PathError.
func Remove(path string) error {
	fs, suffix := findMount(path)
	if fs == nil {
		return &PathError{"remove", path, ErrNotExist}
	}
	err := fs.Remove(suffix)
	if err != nil {
		return err
	}
	return nil
}

// RemoveAll is a stub, it is not implemented.
func RemoveAll(path string) error {
	return ErrNotImplemented
}

// File represents an open file descriptor.
type File struct {
	handle FileHandle
	name   string
}

// Name returns the name of the file with which it was opened.
func (f *File) Name() string {
	return f.name
}

// OpenFile opens the named file. If the operation fails, the returned error
// will be of type *PathError.
func OpenFile(name string, flag int, perm FileMode) (*File, error) {
	fs, suffix := findMount(name)
	if fs == nil {
		return nil, &PathError{"open", name, ErrNotExist}
	}
	handle, err := fs.OpenFile(suffix, flag, perm)
	if err != nil {
		return nil, &PathError{"open", name, err}
	}
	return &File{name: name, handle: handle}, nil
}

// Open opens the file named for reading.
func Open(name string) (*File, error) {
	return OpenFile(name, O_RDONLY, 0)
}

// Create creates the named file, overwriting it if it already exists.
func Create(name string) (*File, error) {
	return OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
}

// Read reads up to len(b) bytes from the File. It returns the number of bytes
// read and any error encountered. At end of file, Read returns 0, io.EOF.
func (f *File) Read(b []byte) (n int, err error) {
	n, err = f.handle.Read(b)
	// TODO: want to always wrap, like upstream, but ReadFile() compares against exactly io.EOF?
	if err != nil && err != io.EOF {
		err = &PathError{"read", f.name, err}
	}
	return
}

var errNegativeOffset = errors.New("negative offset")

// ReadAt reads up to len(b) bytes from the File at the given absolute offset.
// It returns the number of bytes read and any error encountered, possible io.EOF.
// At end of file, Read returns 0, io.EOF.
func (f *File) ReadAt(b []byte, offset int64) (n int, err error) {
	if offset < 0 {
		return 0, &PathError{Op: "readat", Path: f.name, Err: errNegativeOffset}
	}

	for len(b) > 0 {
		m, e := f.handle.ReadAt(b, offset)
		if e != nil {
			// TODO: want to always wrap, like upstream, but TestReadAtEOF compares against exactly io.EOF?
			if e != io.EOF {
				err = &PathError{"readat", f.name, e}
			} else {
				err = e
			}
			break
		}
		n += m
		b = b[m:]
		offset += int64(m)
	}

	return
}

// Write writes len(b) bytes to the File. It returns the number of bytes written
// and an error, if any. Write returns a non-nil error when n != len(b).
func (f *File) Write(b []byte) (n int, err error) {
	n, err = f.handle.Write(b)
	if err != nil {
		err = &PathError{"write", f.name, err}
	}
	return
}

// WriteString is like Write, but writes the contents of string s rather than a
// slice of bytes.
func (f *File) WriteString(s string) (n int, err error) {
	return f.Write([]byte(s))
}

func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	return 0, ErrNotImplemented
}

// Close closes the File, rendering it unusable for I/O.
func (f *File) Close() (err error) {
	err = f.handle.Close()
	if err != nil {
		err = &PathError{"close", f.name, err}
	}
	return
}

// Readdir is a stub, not yet implemented
func (f *File) Readdir(n int) ([]FileInfo, error) {
	return nil, &PathError{"readdir", f.name, ErrNotImplemented}
}

// Readdirnames is a stub, not yet implemented
func (f *File) Readdirnames(n int) (names []string, err error) {
	return nil, &PathError{"readdirnames", f.name, ErrNotImplemented}
}

func (f *File) SyscallConn() (syscall.RawConn, error) {
	return nil, ErrNotImplemented
}

// Fd returns the file handle referencing the open file.
func (f *File) Fd() uintptr {
	panic("unimplemented: os.file.Fd()")
}

// Truncate is a stub, not yet implemented
func (f *File) Truncate(size int64) error {
	return &PathError{"truncate", f.name, ErrNotImplemented}
}

// PathError records an error and the operation and file path that caused it.
// TODO: PathError moved to io/fs in go 1.16 and left an alias in os/errors.go.
// Do the same once we drop support for go 1.15.
type PathError struct {
	Op   string
	Path string
	Err  error
}

func (e *PathError) Error() string {
	return e.Op + " " + e.Path + ": " + e.Err.Error()
}

func (e *PathError) Unwrap() error {
	return e.Err
}

// LinkError records an error during a link or symlink or rename system call and
// the paths that caused it.
type LinkError struct {
	Op  string
	Old string
	New string
	Err error
}

func (e *LinkError) Error() string {
	return e.Op + " " + e.Old + " " + e.New + ": " + e.Err.Error()
}

func (e *LinkError) Unwrap() error {
	return e.Err
}

const (
	O_RDONLY int = syscall.O_RDONLY
	O_WRONLY int = syscall.O_WRONLY
	O_RDWR   int = syscall.O_RDWR
	O_APPEND int = syscall.O_APPEND
	O_CREATE int = syscall.O_CREAT
	O_EXCL   int = syscall.O_EXCL
	O_SYNC   int = syscall.O_SYNC
	O_TRUNC  int = syscall.O_TRUNC
)

func Getwd() (string, error) {
	return syscall.Getwd()
}

// TempDir returns the default directory to use for temporary files.
//
// On Unix systems, it returns $TMPDIR if non-empty, else /tmp.
// On Windows, it uses GetTempPath, returning the first non-empty
// value from %TMP%, %TEMP%, %USERPROFILE%, or the Windows directory.
//
// The directory is neither guaranteed to exist nor have accessible
// permissions.
func TempDir() string {
	return tempDir()
}
