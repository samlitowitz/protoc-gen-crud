package main

import (
	"flag"
	"fmt"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	versionFlag = flag.Bool("version", false, "print protoc-gen-go-crud Version")
)

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version %v, commit %v, built at %v\n", Version, Commit, Date)
		os.Exit(0)
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {

		return nil
	})
}
