package dir

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/webdav"
)

type filesystem struct {
	root *rootFile
	m    map[string]webdav.Dir
}

func newFilesystem(dirs []string) (*filesystem, error) {
	m := make(map[string]webdav.Dir)

	for _, d := range dirs {
		fi, err := os.Stat(d)
		if err != nil {
			return nil, fmt.Errorf("failed to stat '%s': %v", d, err)
		}

		if !fi.IsDir() {
			return nil, fmt.Errorf("'%s' is not directory", d)
		}

		name := filepath.Base(d)
		m[name] = webdav.Dir(d)
	}

	return &filesystem{
		root: &rootFile{
			fi: &rootFileInfo{
				modTime: time.Now(),
			},
			m: m,
		},
		m: m,
	}, nil
}

func (s *filesystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	fmt.Println(name)
	return fmt.Errorf("not supported operation")
}

func (s *filesystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if name == "/" {
		return s.root, nil
	} else {
		for prefix, d := range s.m {
			fList := filepath.SplitList(name)

			for _, f := range fList {
				switch {
				case f == fmt.Sprintf("/%s", prefix):
					return d.OpenFile(ctx, "/", flag, perm)
				case strings.HasPrefix(f, fmt.Sprintf("/%s/", prefix)):
					name := strings.TrimPrefix(f, fmt.Sprintf("/%s", prefix))
					return d.OpenFile(ctx, name, flag, perm)
				}
			}
		}

		return nil, os.ErrNotExist
	}
}

func (s *filesystem) RemoveAll(ctx context.Context, name string) error {
	fmt.Println(name)
	return fmt.Errorf("not suppported operation")
}

func (s *filesystem) Rename(ctx context.Context, oldName, newName string) error {
	fmt.Println(oldName, newName)
	return fmt.Errorf("not suppported")
}

func (s *filesystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	if name == "/" {
		return s.root.fi, nil
	} else {
		for prefix, d := range s.m {
			fList := filepath.SplitList(name)

			for _, f := range fList {
				switch {
				case f == fmt.Sprintf("/%s", prefix):
					return d.Stat(ctx, "/")
				case strings.HasPrefix(f, fmt.Sprintf("/%s/", prefix)):
					name := strings.TrimPrefix(f, fmt.Sprintf("/%s", prefix))
					return d.Stat(ctx, name)
				}
			}
		}

		return nil, os.ErrNotExist
	}
}
