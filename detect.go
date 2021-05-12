package pythonstart

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
)

// BuildPlanMetadata is the buildpack specific data included in build plan
// requirements.
type BuildPlanMetadata struct {

	// Launch flag requests the given requirement be made available during the
	// launch phase of the buildpack lifecycle.
	Launch bool `toml:"launch"`
}

// Detect will return a packit.DetectFunc that will be invoked during the
// detect phase of the buildpack lifecycle.
//
// This buildpack will always pass detection.
// It will require "cpython" as launch-time build plan requirements.
func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		envFile, err := fileExists(filepath.Join(context.WorkingDir, "environment.yml"))
		if err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed trying to stat environment.yml: %w", err)
		}

		lockFile, err := fileExists(filepath.Join(context.WorkingDir, "package-list.txt"))
		if err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed trying to stat package-list.txt: %w", err)
		}

		pythonFiles, err := filepath.Glob(filepath.Join(context.WorkingDir, "*.py"))
		if err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed trying to find *.py files: %w", err)
		}

		if !envFile && !lockFile && len(pythonFiles) < 1 {
			return packit.DetectResult{}, packit.Fail.WithMessage("No *.py, environment.yml or package-list.txt found")
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "cpython",
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
				},
				Or: []packit.BuildPlan{
					{
						Provides: []packit.BuildPlanProvision{},
						Requires: []packit.BuildPlanRequirement{
							{
								Name: "cpython",
								Metadata: BuildPlanMetadata{
									Launch: true,
								},
							},
							{
								Name: "site-packages",
								Metadata: BuildPlanMetadata{
									Launch: true,
								},
							},
						},
					},
					{
						Provides: []packit.BuildPlanProvision{},
						Requires: []packit.BuildPlanRequirement{
							{
								Name: "conda-environment",
								Metadata: BuildPlanMetadata{
									Launch: true,
								},
							},
						},
					},
				},
			},
		}, nil
	}
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
