package web

import (
	"net/http"
	"path/filepath"
)

type NeuteredFileSystem struct {
	Fs http.FileSystem
}

func (nfs NeuteredFileSystem) Open(path string) (http.File, error) {
	file, err := nfs.Fs.Open(path)

	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return file, nil
	}

	indexPath := filepath.Join(path, "index.html")

	index, err := nfs.Fs.Open(indexPath)
	if err != nil {
		file.Close()

		return nil, err
	}
	index.Close()

	return file, nil
}
