package pythonstart

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

// Build will return a packit.BuildFunc that will be invoked during the build
// phase of the buildpack lifecycle.
//
// Build assigns the image a launch process to run the Python REPL.
func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		enableTini, err := isEnvVarTrue(LaunchWithTiniEnvName)
		if err != nil {
			return packit.BuildResult{}, err
		}

		var originalProcess packit.Process
		if enableTini {
			originalProcess = packit.Process{
				Type:    Web,
				Command: Tini,
				Args:    []string{"-g", "--", Python},
				Default: true,
				Direct:  true,
			}
			logger.Process("Using tini for process launching")
		} else {
			originalProcess = packit.Process{
				Type:    Web,
				Command: Python,
				Default: true,
				Direct:  true,
			}
		}

		processes := []packit.Process{originalProcess}

		logger.LaunchProcesses(processes)

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: processes,
			},
		}, nil
	}
}
