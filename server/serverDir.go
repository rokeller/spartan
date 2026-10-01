package server

import (
	"net/http"
	"path"
)

type serverDir struct {
	http.Dir
}

func newServerDir(dir http.Dir) http.FileSystem {
	return &serverDir{Dir: dir}
}

func (d *serverDir) Open(name string) (result http.File, err error) {
	f, err := d.Dir.Open(name)
	if err != nil {
		return
	}

	fi, err := f.Stat()
	if err != nil {
		return
	}
	if fi.IsDir() {
		return d.openForDir(name, f)
	}
	return f, nil
}

func (d *serverDir) openForDir(name string, f http.File) (http.File, error) {
	index := path.Join(name, "index.html")
	ff, err := d.Dir.Open(index)
	if err == nil {
		defer ff.Close()
		// Return the originally opened http.File because that's what the caller
		// expects here.
		return f, nil
	} else {
		// There isn't an 'index.html' file in the directory or there are other
		// errors with it. In such cases the http.FileServer can respond with a
		// listing of the contents of the directory, which we do NOT want. So
		// let's return the error back to the caller.
		return nil, err
	}
}
