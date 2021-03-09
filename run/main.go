package main

import (
	"github.com/paketo-buildpacks/packit"
	pythonstart "github.com/paketo-community/python-start"
)

func main() {
	packit.Run(pythonstart.Detect(), pythonstart.Build())
}
