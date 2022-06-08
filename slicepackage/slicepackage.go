package slicepackage

type Package struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Architecture string   `json:"architecture"`
	Maintainer   string   `json:"maintainer"`
	Dependencies []string `json:"dependencies"`
	Description  string   `json:"description"`
}
