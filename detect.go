package pythonstart

import "github.com/paketo-buildpacks/packit"

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
