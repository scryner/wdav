package dir

import (
	"context"
	"fmt"
	"io/fs"
	"time"

	"golang.org/x/net/webdav"
)

type rootFileInfo struct {
	modTime time.Time
}

func (r *rootFileInfo) Name() string {
	return "/"
}

func (r *rootFileInfo) Size() int64 {
	return 0
}

func (r *rootFileInfo) Mode() fs.FileMode {
	return fs.ModeDir
}

func (r *rootFileInfo) ModTime() time.Time {
	return r.modTime
}

func (r *rootFileInfo) IsDir() bool {
	return true
}

func (r *rootFileInfo) Sys() interface{} {
	panic("implement me")
}

type rootFile struct {
	fi *rootFileInfo
	m map[string]webdav.Dir
}

func (r *rootFile) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("unsupported operation 'write'")
}

func (r *rootFile) Close() error {
	return nil
}

func (r *rootFile) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("unsupported operation 'read'")
}

func (r *rootFile) Seek(offset int64, whence int) (int64, error) {
	return 0, fmt.Errorf("unsupported opration 'seek'")
}

func (r *rootFile) Readdir(count int) ([]fs.FileInfo, error) {
	var fis []fs.FileInfo

	for _, d := range r.m {
		fi, err := d.Stat(context.Background(), "/")
		if err != nil {
			return nil, err
		}

		fis = append(fis, fi)
	}

	return fis, nil
}

func (r *rootFile) Stat() (fs.FileInfo, error) {
	return r.fi, nil
}
