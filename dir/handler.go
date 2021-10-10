package dir

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

func NewHandler(prefix string, dirs []string) (*webdav.Handler, error) {
	if len(dirs) < 1 {
		return nil, fmt.Errorf("insufficient directories")
	}

	// logger func
	loggerFunc := func(request *http.Request, err error) {
		if err != nil {
			log.Println(err)
		}
	}

	// lock system
	ls := webdav.NewMemLS()

	// filesystem
	var fs webdav.FileSystem

	if len(dirs) == 1 {
		// stat directories
		d := dirs[0]
		fi, err := os.Stat(d)
		if err != nil {
			return nil, fmt.Errorf("failed to stat on '%s': %v", d, err)
		}

		if !fi.IsDir() {
			return nil, fmt.Errorf("'%s' is not directory", d)
		}

		fs = webdav.Dir(d)

	} else {
		var err error
		fs, err = newFilesystem(dirs)
		if err != nil {
			return nil, fmt.Errorf("failed to make filesystem: %v", err)
		}
	}

	// return
	return &webdav.Handler{
		Prefix:     fmt.Sprintf("/%s", prefix),
		FileSystem: fs,
		LockSystem: ls,
		Logger:     loggerFunc,
	}, nil
}
