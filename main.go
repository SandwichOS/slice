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

		fmt.Println("Installing " + packageMetadata.Name + " (Architecture: " + packageMetadata.Architecture + ")...")
	default:
		fmt.Println("???")
	}
}
