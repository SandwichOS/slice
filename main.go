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

		// Parse file

		var packageMetadata slicepackage.Package

		err = json.Unmarshal(data, &packageMetadata)

		if err != nil {
			panic(err)
		}

		fmt.Println("Installing package: " + packageMetadata.Name + " (Architecture: " + packageMetadata.Architecture + ")...")
	default:
		fmt.Println("???")
	}
}
