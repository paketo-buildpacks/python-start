package pythonstart

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

// Build will return a packit.BuildFunc that will be invoked during the build
// phase of the buildpack lifecycle.
//
// Build assigns the image a launch process to run the Python REPL.
func Build(logger scribe.Logger) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logger.Process("Assigning launch process")
		command := "python"
		logger.Subprocess("web: %s", command)

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{
					{
						Type:    "web",
						Command: command,
						Default: true,
					},
				},
			},
		}, nil
	}
}
