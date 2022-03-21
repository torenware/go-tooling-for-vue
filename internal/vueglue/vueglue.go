package vueglue

import (
	"embed"
	"errors"
	"io/fs"
	"path/filepath"
)

const (
	AssetsDir = "dist/assets"
)

type VueGlue struct {
	MainModule   string
	VendorModule string
	CSSModule    string
	DistFS       *embed.FS
}

func NewVueGlue(dist *embed.FS) (*VueGlue, error) {
	var glue VueGlue
	prefix := "/dist"
	removeChars := len(prefix) - 1

	if !fs.ValidPath(AssetsDir) {
		return nil, errors.New("vite dist directory not found")
	}
	glue.DistFS = dist

	jsmodGlob, err := fs.Glob(dist, filepath.Join(AssetsDir, "index.*.js"))
	if err != nil {
		return nil, err
	}
	if len(jsmodGlob) == 0 {
		return nil, errors.New("jsmod file not in file system")
	}
	glue.MainModule = jsmodGlob[0][removeChars:]

	jsVendorGlob, err := fs.Glob(dist, filepath.Join(AssetsDir, "vendor.*.js"))
	if err != nil {
		return nil, err
	}
	if len(jsVendorGlob) == 0 {
		return nil, errors.New("vendor bundle not in file system")
	}
	glue.VendorModule = jsVendorGlob[0][removeChars:]

	cssGlob, err := fs.Glob(dist, filepath.Join(AssetsDir, "index.*.css"))
	if err != nil {
		return nil, err
	}
	if len(cssGlob) == 0 {
		return nil, errors.New("css bundle not in file system")
	}
	glue.CSSModule = cssGlob[0][removeChars:]

	return &glue, nil

}
