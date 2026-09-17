package pythonstart

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/fs"
)

// BuildPlanMetadata is the buildpack specific data included in build plan
// requirements.
type BuildPlanMetadata struct {

	// Launch flag requests the given requirement be made available during the
	// launch phase of the buildpack lifecycle.
	Launch bool `toml:"launch"`
	Build  bool `toml:"build"`
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
		envFile, err := fs.Exists(filepath.Join(context.WorkingDir, "environment.yml"))
		if err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed trying to stat environment.yml: %w", err)
		}

		pixiEnvFile, err := fs.Exists(filepath.Join(context.WorkingDir, "pixi.lock"))
		if err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed trying to stat pixi.lock: %w", err)
		}

		requirementsFile, err := fs.Exists(filepath.Join(context.WorkingDir, "requirements.txt"))
		if err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed trying to stat requirements.txt: %w", err)
		}

		lockFile, err := fs.Exists(filepath.Join(context.WorkingDir, "package-list.txt"))
		if err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed trying to stat package-list.txt: %w", err)
		}

		uvLockFile, err := fs.Exists(filepath.Join(context.WorkingDir, "uv.lock"))
		if err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed trying to stat uv.lock: %w", err)
		}

		pipenvLockFile, err := fs.Exists(filepath.Join(context.WorkingDir, "Pipfile.lock"))
		if err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed trying to stat Pipfile.lock: %w", err)
		}

		pyprojectTOMLFile, err := fs.Exists(filepath.Join(context.WorkingDir, "pyproject.toml"))
		if err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed trying to stat pyproject.toml: %w", err)
		}

		pythonFiles, err := filepath.Glob(filepath.Join(context.WorkingDir, "*.py"))
		if err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed trying to find *.py files: %w", err)
		}

		if !envFile &&
			!pixiEnvFile &&
			!requirementsFile &&
			!lockFile &&
			!uvLockFile &&
			!pipenvLockFile &&
			!pyprojectTOMLFile &&
			len(pythonFiles) < 1 {
			return packit.DetectResult{}, packit.Fail.WithMessage("No *.py, environment.yml, pixi.lock, requirements.txt, uv.lock, Pipfile.lock, pyproject.toml, or package-list.txt found")
		}

		simplePlan := packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				NewLaunchRequirement(CPython),
			},
		}

		pipPlan := packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				NewLaunchRequirement(CPython),
				NewLaunchRequirement(SitePackages),
			},
		}

		condaPlan := packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				NewLaunchRequirement(CondaEnv),
			},
		}

		pixiPlan := packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				NewLaunchRequirement(PixiEnv),
			},
		}

		uvPlan := packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				NewLaunchRequirement(UvEnv),
			},
		}

		pipenvPlan := packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				NewLaunchRequirement(CPython),
				NewLaunchRequirement(SitePackages),
				NewLaunchRequirement(PipEnv),
			},
		}

		poetryInstallPlan := packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				NewLaunchRequirement(CPython),
				NewLaunchRequirement(Poetry),
				NewLaunchRequirement(PoetryVenv),
			},
		}

		plans := []packit.BuildPlan{pipPlan, pipenvPlan, condaPlan, pixiPlan, poetryInstallPlan, simplePlan}

		// The current build plan from the python buildpack will make an uv project be detected as
		// a poetry project due to the pyproject.toml presence hence we workaround the issue by only
		// requiring uv.
		// Note: this is temporary.
		if uvLockFile {
			plans = []packit.BuildPlan{uvPlan}
		}

		if shouldReload, err := isEnvVarTrue(LiveReloadEnvName); err != nil {
			return packit.DetectResult{}, err
		} else if shouldReload {
			for i := range plans {
				plans[i].Requires = append(plans[i].Requires,
					NewLaunchRequirement(WatchExec))
			}
		}

		if shouldUsePackageManagers, err := isEnvVarTrue(PackageManagersEnvName); err != nil {
			return packit.DetectResult{}, err
		} else if shouldUsePackageManagers {
			for i := range plans {
				// Simple plan does not use package-managers
				if len(plans) > 1 && i == len(plans)-1 {
					continue
				}
				plans[i].Requires = append(plans[i].Requires,
					NewBuildRequirement(PackageManagersPlanEntry))
			}
		}

		if shouldLaunchWithTini, err := isEnvVarTrue(LaunchWithTiniEnvName); err != nil {
			return packit.DetectResult{}, err
		} else if shouldLaunchWithTini {
			for i := range plans {
				plans[i].Requires = append(plans[i].Requires,
					NewLaunchRequirement(Tini))
			}
		}

		return packit.DetectResult{
			Plan: or(plans...),
		}, nil
	}
}

func NewBuildRequirement(name string) packit.BuildPlanRequirement {
	return packit.BuildPlanRequirement{
		Name: name,
		Metadata: BuildPlanMetadata{
			Build: true,
		},
	}
}

func NewLaunchRequirement(name string) packit.BuildPlanRequirement {
	return packit.BuildPlanRequirement{
		Name: name,
		Metadata: BuildPlanMetadata{
			Launch: true,
		},
	}
}

func isEnvVarTrue(varName string) (bool, error) {
	if value, found := os.LookupEnv(varName); found {
		enable, err := strconv.ParseBool(value)
		if err != nil {
			return false, fmt.Errorf("failed to parse %s value %s: %w", varName, value, err)
		}
		return enable, nil
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
