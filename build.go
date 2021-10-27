package pythonstart

import (
	"fmt"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
)

// Build will return a packit.BuildFunc that will be invoked during the build
// phase of the buildpack lifecycle.
//
// Build assigns the image a launch process to run the Python REPL. If
// BP_LIVE_RELOAD_ENABLED=true in the build environment, it will make the
// default process reloadable; the process will restart whenever the contents
// of the app's working directory changes in the app container.
func Build(logger scribe.Logger) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		command := "python"
		processes := []packit.Process{
			{
				Type:    "web",
				Command: command,
			},
		}

		shouldReload, err := checkLiveReloadEnabled()
		if err != nil {
			return packit.BuildResult{}, err
		}

		if shouldReload {
			processes = []packit.Process{
				{
					Type:    "web",
					Command: fmt.Sprintf(`watchexec --restart --watch %s "%s"`, context.WorkingDir, command),
				},
				{
					Type:    "no-reload",
					Command: command,
				},
			}
		}

		logger.Process("Assigning launch process")
		for _, process := range processes {
			logger.Subprocess("%s: %s", process.Type, process.Command)
		}

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: processes,
			},
		}, nil
	}
}
