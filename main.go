package main

import (
	"flag"
	"fmt"
	"github.com/docker/go-plugins-helpers/volume"
	"path/filepath"
)

const (
	pluginName = "ebs"
)

var (
	root = flag.String("root", volume.DefaultDockerRootDirectory, "Docker volumes root directory")
)

type ebsDriver struct {
	root string
}

func (e ebsDriver) Create(r volume.Request) volume.Response {
	fmt.Printf("Create %v\n", r)

	return volume.Response{}
}

func (e ebsDriver) Remove(r volume.Request) volume.Response {
	fmt.Printf("Remove %v\n", r)

	return volume.Response{}
}

func (e ebsDriver) Path(r volume.Request) volume.Response {
	fmt.Printf("Path %v\n", r)

	return volume.Response{Mountpoint: filepath.Join(e.root, r.Name)}
}

func (e ebsDriver) Mount(r volume.Request) volume.Response {
	p := filepath.Join(e.root, r.Name)
	fmt.Printf("Mount %s\n", p)

	return volume.Response{Mountpoint: p}
}

func (e ebsDriver) Unmount(r volume.Request) volume.Response {
	p := filepath.Join(e.root, r.Name)
	fmt.Printf("Unmount %s\n", p)

	return volume.Response{}
}

func main() {
	d := ebsDriver{*root}
	h := volume.NewHandler(d)
	h.ServeUnix("root", pluginName)
	fmt.Printf("Starting listening on unix socket with name %s (usually in /run/docker/plugins/%s)", pluginName, pluginName)
}
