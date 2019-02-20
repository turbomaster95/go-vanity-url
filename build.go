package main

import (
	"os"
	"fmt"
	"bytes"
	"html/template"
	"path/filepath"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Imports struct {
  Package []Package
}

type Package struct {
	Repo     string `toml:"repo"`
	Path     string `toml:"path"`
	Packages []string `toml:"packages"`
	Url      string `toml:"url"`
}

func (p Package) GetRoutes() (ret []string) {
	ret = append(ret, p.Path)
	for _, pkg := range p.Packages {
		ret = append(ret, p.Path+"/"+pkg)
	}
	return
}

func loadConfig() (ret []Package, err error) {
	content, err := ioutil.ReadFile("packages.toml")
	if err != nil {
		return
	}
	
	var im Imports
	if _, err = toml.Decode(string(content), &im); err != nil {
		return
	}
	ret = im.Package
	return
}

func loadRoutes(packages []Package) (routes map[string]Package) {
	routes = map[string]Package{}

	for _, pkg := range packages {
		for _, route := range pkg.GetRoutes() {
			routes[route] = pkg
		}
	}
	return
}


func main() {
	packages, err := loadConfig()
	if err != nil {
		fmt.Errorf("%s\n", err)
		return
	}
	routes := loadRoutes(packages)
	t, err := template.ParseFiles("template.tmpl")
	if err != nil {
		fmt.Errorf("%s\n", err)
		return
	}
	for route, pkg := range routes {
		// generate file and write out
		pkg.Path = "src.techknowlogick.com"+pkg.Path
		pkgDir := filepath.Join("public", route)
		if err := os.MkdirAll(pkgDir, 0755); err != nil {
			fmt.Errorf("%s\n", err)
			return
		}
		var page bytes.Buffer
		if err := t.Execute(&page, pkg); err != nil {
			fmt.Errorf("%s\n", err)
			return
		}

		pkgFile := filepath.Join(pkgDir, "index.html")
		if err := ioutil.WriteFile(pkgFile, page.Bytes(), 0755); err != nil {
			fmt.Errorf("%s\n", err)
			return
		}
	}
}
