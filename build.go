package main

import (
	"bytes"
	"github.com/pelletier/go-toml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Package struct {
	Repo     string
	Path     string
	Packages []string
	Url      string
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
  toml.Unmarshal(content, &ret)
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
}
