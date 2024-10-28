package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
)

var (
	dataPath     = flag.String("datapath", filepath.Join("./", "data"), "Path to your custom 'data' directory")
	datName      = flag.String("datname", "geosite.dat", "Name of the generated dat file")
	outputPath   = flag.String("outputpath", "./publish", "Output path to the generated files")
)

func main() {
	flag.Parse()

	dir := GetDataDir()
	listInfoMap := make(ListInfoMap)

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if err := listInfoMap.Marshal(path); err != nil {
			return err
		}
		return nil
	}); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}

	if err := listInfoMap.FlattenAndGenUniqueDomainList(); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}

	// Read rules
	RulesInFile := make(map[fileName]map[attribute]bool)

	// Generate dlc.dat
	if geositeList := listInfoMap.ToProto(RulesInFile); geositeList != nil {
		protoBytes, err := proto.Marshal(geositeList)
		if err != nil {
			fmt.Println("Failed:", err)
			os.Exit(1)
		}
		if err := os.MkdirAll(*outputPath, 0755); err != nil {
			fmt.Println("Failed:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(filepath.Join(*outputPath, *datName), protoBytes, 0644); err != nil {
			fmt.Println("Failed:", err)
			os.Exit(1)
		} else {
			fmt.Printf("%s has been generated successfully in '%s'.\n", *datName, *outputPath)
		}
	}

}
