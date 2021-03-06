package importer

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type customImporter struct {
	imported map[string]*types.Package
	base     types.Importer
}

func (i *customImporter) Import(path string) (*types.Package, error) {
	if pkg, ok := i.imported[path]; ok {
		return pkg, nil
	}
	pkg, err := i.fsPkg(path)
	if err != nil {
		return nil, err
	}
	i.imported[path] = pkg
	return pkg, nil
}

func gopathDir(pkg string) (string, error) {
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		absPath, err := filepath.Abs(path.Join(gopath, "src", pkg))
		if err != nil {
			return "", err
		}
		if dir, err := os.Stat(absPath); err == nil && dir.IsDir() {
			return absPath, nil
		}
	}
	return "", fmt.Errorf("%s not in $GOPATH", pkg)
}

func removeGopath(p string) string {
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		p = strings.Replace(p, path.Join(gopath, "src")+"/", "", 1)
	}
	return p
}

func (i *customImporter) fsPkg(pkg string) (*types.Package, error) {
	dir, err := gopathDir(pkg)
	if err != nil {
		return i.base.Import(pkg)
	}

	dirFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return i.base.Import(pkg)
	}

	fset := token.NewFileSet()
	var files []*ast.File
	for _, fileInfo := range dirFiles {
		if fileInfo.IsDir() {
			continue
		}
		n := fileInfo.Name()
		if path.Ext(fileInfo.Name()) != ".go" {
			continue
		}
		file := path.Join(dir, n)
		src, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		f, err := parser.ParseFile(fset, file, src, 0)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	conf := types.Config{Importer: i}
	p, err := conf.Check(pkg, fset, files, nil)

	if err != nil {
		p, err = i.base.Import(pkg)
	}
	return p, nil
}

// Default returns an importer that will try to import code from gopath before using go/importer.Default
func Default() types.Importer {
	return &customImporter{
		make(map[string]*types.Package),
		importer.Default(),
	}
}
