package pythonstart

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/paketo-buildpacks/packit/v2"
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
// If this buildpack detects files that indicate your app is a Python project,
// it will pass detection. It will require "cpython" OR "cpython" and
// "site-packages" OR "conda-environment" as launch-time build plan
// requirements, depending on whether it detects files indicating the use of
// different package managers.
//
// If BP_LIVE_RELOAD_ENABLED=true in the build environment, it will
// additionally require "watchexec" at launch-time
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

		simplePlan := packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				{
					Name: "cpython",
					Metadata: BuildPlanMetadata{
						Launch: true,
					},
				},
			},
		}

		pipPlan := packit.BuildPlan{
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
		}

		condaPlan := packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				{
					Name: "conda-environment",
					Metadata: BuildPlanMetadata{
						Launch: true,
					},
				},
			},
		}

		poetryInstallPlan := packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				{
					Name: "cpython",
					Metadata: BuildPlanMetadata{
						Launch: true,
					},
				},
				{
					Name: "poetry-venv",
					Metadata: BuildPlanMetadata{
						Launch: true,
					},
				},
			},
		}

		plans := []packit.BuildPlan{simplePlan, pipPlan, condaPlan, poetryInstallPlan}

		shouldReload, err := checkLiveReloadEnabled()
		if err != nil {
			return packit.DetectResult{}, err
		}

		if shouldReload {
			for i := range plans {
				plans[i].Requires = append(plans[i].Requires, packit.BuildPlanRequirement{
					Name: "watchexec",
					Metadata: BuildPlanMetadata{
						Launch: true,
					},
				})
			}
		}
		return packit.DetectResult{
			Plan: or(plans...),
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

func checkLiveReloadEnabled() (bool, error) {
	if reload, ok := os.LookupEnv("BP_LIVE_RELOAD_ENABLED"); ok {
		shouldEnableReload, err := strconv.ParseBool(reload)
		if err != nil {
			return false, fmt.Errorf("failed to parse BP_LIVE_RELOAD_ENABLED value %s: %w", reload, err)
		}
		return shouldEnableReload, nil
	}
	return false, nil
}

func or(plans ...packit.BuildPlan) packit.BuildPlan {
	if len(plans) < 1 {
		return packit.BuildPlan{}
	}
	combinedPlan := plans[0]

	for i := range plans {
		if i == 0 {
			continue
		}
		combinedPlan.Or = append(combinedPlan.Or, plans[i])
	}
	return combinedPlan
}
