package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type dirHideSymLink struct {
	afero.Fs
}

type fileHideSymlink struct {
	afero.File
	fs afero.Fs
}

func (h dirHideSymLink) Create(name string) (afero.File, error) {
	f, err := h.Fs.Create(name)
	if err == nil {
		f = fileHideSymlink{File: f, fs: h}
	}
	return f, err
}

func (h dirHideSymLink) Open(name string) (afero.File, error) {
	f, err := h.Fs.Open(name)
	if err == nil {
		f = fileHideSymlink{File: f, fs: h}
	}
	return f, err
}

func (h dirHideSymLink) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := h.Fs.OpenFile(name, flag, perm)
	if err == nil {
		f = fileHideSymlink{File: f, fs: h}
	}
	return f, err
}

func (h fileHideSymlink) Readdir(count int) ([]fs.FileInfo, error) {
	fis, err := h.File.Readdir(count)
	if err != nil {
		return nil, err
	}
	for i, fi := range fis {
		if (fi.Mode() & fs.ModeSymlink) != 0 {
			// replace with resolved FileInfo from Stat()
			if fis[i], err = h.fs.Stat(filepath.Join(h.Name(), fi.Name())); err != nil {
				return nil, err
			}
		}
	}
	return fis, nil
}

func main() {
	ofs := afero.NewOsFs()
	bfs := afero.NewBasePathFs(dirHideSymLink{ofs}, "/home/thhal/go/src/github.com/thallgren/scratch")
	fis, err := afero.ReadDir(bfs, "fs")
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range fis {
		var fm os.FileMode = fi.Mode()
		fmt.Printf("%b %s\n", fm, fm)
	}
}
