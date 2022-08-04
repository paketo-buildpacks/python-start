package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	pythonstart "github.com/paketo-buildpacks/python-start"
)

func main() {
	packit.Run(
		pythonstart.Detect(),
		pythonstart.Build(scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))),
	)
}
