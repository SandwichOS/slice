package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/SandwichOS/slice/slicepackage"
)

func main() {
	fmt.Println("~-~-Slice!~-~-")

	switch os.Args[1] {
	case "create":
		// Read metadata file

		data, err := ioutil.ReadFile(os.Args[2] + "/metadata.json")

		if err != nil {
			panic(err)
		}

		// Parse metadata file

		var packageMetadata slicepackage.Package

		err = json.Unmarshal(data, &packageMetadata)

		if err != nil {
			panic(err)
		}

		fmt.Println("Creating package: " + packageMetadata.Name + " (Architecture: " + packageMetadata.Architecture + ")...")

		data, err = slicepackage.CreatePackageTarball(os.Args[2])

		if err != nil {
			panic(err)
		}

		ioutil.WriteFile(os.Args[3], slicepackage.CompressData(data), 775)

	case "install":
		// Read file

		data, err := ioutil.ReadFile(os.Args[2])

		if err != nil {
			panic(err)
		}

		decompressedData := slicepackage.DecompressData(data)

		packageMetadata, err := slicepackage.GetPackageMetadata(decompressedData)

		if err != nil {
			panic(err)
		}

		installDirectory, ok := os.LookupEnv("SLICE_DESTDIR")

		if !ok {
			installDirectory = "/"
		}

		if installDirectory == "/" {
			fmt.Println("Installing package: " + packageMetadata.Name + " (Architecture: " + packageMetadata.Architecture + ")...")
		} else {
			fmt.Println("Installing package: " + packageMetadata.Name + " (Architecture: " + packageMetadata.Architecture + ") to " + installDirectory + "...")
		}

		slicepackage.ExtractPackageTarball(decompressedData, installDirectory)
	case "info":
		// Read file

		data, err := ioutil.ReadFile(os.Args[2])

		if err != nil {
			panic(err)
		}

		decompressedData := slicepackage.DecompressData(data)

		packageMetadata, err := slicepackage.GetPackageMetadata(decompressedData)

		if err != nil {
			panic(err)
		}

		fmt.Println("Package name:", packageMetadata.Name)
		fmt.Println("Package version:", packageMetadata.Version)
		fmt.Println("Package architecture:", packageMetadata.Architecture)
		fmt.Println("Package Maintainer:", packageMetadata.Maintainer)
		fmt.Println("Package Description:", packageMetadata.Description)
	case "remove":
		// Read file

		data, err := ioutil.ReadFile(os.Args[2])

		if err != nil {
			panic(err)
		}

		decompressedData := slicepackage.DecompressData(data)

		packageMetadata, err := slicepackage.GetPackageMetadata(decompressedData)

		if err != nil {
			panic(err)
		}

		installDirectory, ok := os.LookupEnv("SLICE_DESTDIR")

		if !ok {
			installDirectory = "/"
		}

		if installDirectory == "/" {
			fmt.Println("Removing package: " + packageMetadata.Name + " (Architecture: " + packageMetadata.Architecture + ")...")
		} else {
			fmt.Println("Removing package: " + packageMetadata.Name + " (Architecture: " + packageMetadata.Architecture + ") at " + installDirectory + "...")
		}

		slicepackage.RemovePackage(decompressedData, installDirectory)
	default:
		fmt.Println("???")
	}
}
