package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
	"os"
	"path/filepath"
	"strings"
	"time"
	"xiam.li/meta"
)

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "print the version and exit")
	flag.BoolVar(&showVersion, "v", false, "print the version and exit")
	flag.Parse()
	if showVersion {
		fmt.Printf("%s %s (%s), built on %s\n", filepath.Base(os.Args[0]), meta.VersionOr("v0.0.0-dev+dirty"), meta.ShortSHAOr(strings.Repeat("x", 40)), meta.DateOr(time.Now()))
		return
	}

	var (
		flags flag.FlagSet
	)
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if f.Generate {
				generateFile(gen, f)
			}
		}
		return nil
	})
}
