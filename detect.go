package pythonstart

import "github.com/paketo-buildpacks/packit"

type BuildPlanMetadata struct {
	Launch bool `toml:"launch"`
}

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
			},
		}, nil
	}
}
