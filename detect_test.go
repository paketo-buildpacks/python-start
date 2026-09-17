package pythonstart_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit/v2"
	pythonstart "github.com/paketo-buildpacks/python-start"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir string
		detect     packit.DetectFunc

		expectedFullPlan packit.BuildPlan
		expectedUvPlan   packit.BuildPlan
	)

	it.Before(func() {
		var err error
		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		Expect(os.WriteFile(filepath.Join(workingDir, "x.py"), []byte{}, os.ModePerm)).To(Succeed())

		detect = pythonstart.Detect()

		expectedFullPlan = packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				pythonstart.NewLaunchRequirement(pythonstart.CPython),
				pythonstart.NewLaunchRequirement(pythonstart.SitePackages),
			},
			Or: []packit.BuildPlan{
				{
					Provides: []packit.BuildPlanProvision{},
					Requires: []packit.BuildPlanRequirement{
						pythonstart.NewLaunchRequirement(pythonstart.CPython),
						pythonstart.NewLaunchRequirement(pythonstart.SitePackages),
						pythonstart.NewLaunchRequirement(pythonstart.PipEnv),
					},
				},
				{
					Provides: []packit.BuildPlanProvision{},
					Requires: []packit.BuildPlanRequirement{
						pythonstart.NewLaunchRequirement(pythonstart.CondaEnv),
					},
				},
				{
					Provides: []packit.BuildPlanProvision{},
					Requires: []packit.BuildPlanRequirement{
						pythonstart.NewLaunchRequirement(pythonstart.PixiEnv),
					},
				},
				{
					Provides: []packit.BuildPlanProvision{},
					Requires: []packit.BuildPlanRequirement{
						pythonstart.NewLaunchRequirement(pythonstart.CPython),
						pythonstart.NewLaunchRequirement(pythonstart.Poetry),
						pythonstart.NewLaunchRequirement(pythonstart.PoetryVenv),
					},
				},
				{
					Provides: []packit.BuildPlanProvision{},
					Requires: []packit.BuildPlanRequirement{
						pythonstart.NewLaunchRequirement(pythonstart.CPython),
					},
				},
			},
		}
		expectedUvPlan = packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				pythonstart.NewLaunchRequirement(pythonstart.UvEnv),
			},
		}
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("detection phase", func() {
		context("without uv.lock file", func() {
			it("detects", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan).To(Equal(expectedFullPlan))
			})
		})

		context("with uv.lock file", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(workingDir, "uv.lock"), []byte{}, os.ModePerm)).To(Succeed())
			})

			it("detects", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan).To(Equal(expectedUvPlan))
			})
		})

		context("when BP_LIVE_RELOAD_ENABLED=true in the build environment", func() {
			it.Before(func() {
				t.Setenv(pythonstart.LiveReloadEnvName, "true")
			})

			context("without uv.lock", func() {
				it("requires watchexec at launch", func() {
					result, err := detect(packit.DetectContext{
						WorkingDir: workingDir,
					})
					watchexecBuildPlan := packit.BuildPlan{}
					watchexecBuildPlan.Provides = expectedFullPlan.Provides
					watchexecBuildPlan.Requires = append(expectedFullPlan.Requires, pythonstart.NewLaunchRequirement(pythonstart.WatchExec))
					watchexecBuildPlan.Or = make([]packit.BuildPlan, len(expectedFullPlan.Or))
					for i := range expectedFullPlan.Or {
						watchexecBuildPlan.Or[i].Provides = expectedFullPlan.Or[i].Provides
						watchexecBuildPlan.Or[i].Requires = append(expectedFullPlan.Or[i].Requires,
							pythonstart.NewLaunchRequirement(pythonstart.WatchExec))
					}
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Plan).To(Equal(watchexecBuildPlan))
				})
			})

			context("with uv.lock file", func() {
				it.Before(func() {
					Expect(os.WriteFile(filepath.Join(workingDir, "uv.lock"), []byte{}, os.ModePerm)).To(Succeed())
				})

				it("detects", func() {
					result, err := detect(packit.DetectContext{
						WorkingDir: workingDir,
					})
					watchexecBuildPlan := packit.BuildPlan{}
					watchexecBuildPlan.Provides = expectedUvPlan.Provides
					watchexecBuildPlan.Requires = append(expectedUvPlan.Requires, pythonstart.NewLaunchRequirement(pythonstart.WatchExec))
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Plan).To(Equal(watchexecBuildPlan))
				})
			})
		})

		context("when BP_LAUNCH_WITH_TINI is true", func() {
			it.Before(func() {
				t.Setenv(pythonstart.LaunchWithTiniEnvName, "true")
				Expect(os.WriteFile(filepath.Join(workingDir, "uv.lock"), []byte{}, os.ModePerm)).To(Succeed())
			})

			it("requires tini at launch time", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan.Requires).To(Equal([]packit.BuildPlanRequirement{
					pythonstart.NewLaunchRequirement(pythonstart.UvEnv),
					pythonstart.NewLaunchRequirement(pythonstart.Tini),
				}))
			})
		})

		context("when BP_LAUNCH_WITH_TINI is malformed", func() {
			it.Before(func() {
				t.Setenv(pythonstart.LaunchWithTiniEnvName, "not-a-bool")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("failed to parse BP_LAUNCH_WITH_TINI value not-a-bool")))
			})
		})

		context("when BP_ENABLE_PACKAGE_MANAGERS=true in the build environment", func() {
			it.Before(func() {
				t.Setenv(pythonstart.PackageManagersEnvName, "true")
			})

			context("without uv.lock", func() {
				it("requires package managers at launch", func() {
					result, err := detect(packit.DetectContext{
						WorkingDir: workingDir,
					})
					pmBuildPlan := packit.BuildPlan{}
					pmBuildPlan.Provides = expectedFullPlan.Provides
					pmBuildPlan.Requires = append(expectedFullPlan.Requires, pythonstart.NewBuildRequirement(pythonstart.PackageManagersPlanEntry))
					pmBuildPlan.Or = make([]packit.BuildPlan, len(expectedFullPlan.Or))
					for i := range expectedFullPlan.Or {
						pmBuildPlan.Or[i].Provides = expectedFullPlan.Or[i].Provides
						if i == len(expectedFullPlan.Or) - 1 {
							pmBuildPlan.Or[i].Requires = expectedFullPlan.Or[i].Requires
						} else {

						pmBuildPlan.Or[i].Requires = append(expectedFullPlan.Or[i].Requires,
							pythonstart.NewBuildRequirement(pythonstart.PackageManagersPlanEntry))
						}
					}
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Plan).To(Equal(pmBuildPlan))
				})
			})

			context("with uv.lock file", func() {
				it.Before(func() {
					Expect(os.WriteFile(filepath.Join(workingDir, "uv.lock"), []byte{}, os.ModePerm)).To(Succeed())

				})

				it("detects", func() {
					result, err := detect(packit.DetectContext{
						WorkingDir: workingDir,
					})
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Plan).To(Equal(packit.BuildPlan{
						Provides: []packit.BuildPlanProvision{},
						Requires: []packit.BuildPlanRequirement{
							{
								Name: pythonstart.UvEnv,
								Metadata: pythonstart.BuildPlanMetadata{
									Launch: true,
								},
							},
							{
								Name: pythonstart.PackageManagersPlanEntry,
								Metadata: pythonstart.BuildPlanMetadata{
									Build: true,
								},
							},
						},
					}))
				})
			})
		})

		context("When only an environment.yml file is present", func() {
			it.Before(func() {
				Expect(os.RemoveAll(filepath.Join(workingDir, "x.py"))).To(Succeed())
				Expect(os.WriteFile(filepath.Join(workingDir, "environment.yml"), []byte{}, os.ModePerm)).To(Succeed())
			})

			it("passes detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
			})
		})

		context("When only a requirements.txt file is present", func() {
			it.Before(func() {
				Expect(os.RemoveAll(filepath.Join(workingDir, "x.py"))).To(Succeed())
				Expect(os.WriteFile(filepath.Join(workingDir, "requirements.txt"), []byte{}, os.ModePerm)).To(Succeed())
			})

			it("passes detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
			})
		})

		context("When only a package-list.txt file is present", func() {
			it.Before(func() {
				Expect(os.RemoveAll(filepath.Join(workingDir, "x.py"))).To(Succeed())
				Expect(os.WriteFile(filepath.Join(workingDir, "package-list.txt"), []byte{}, os.ModePerm)).To(Succeed())
			})

			it("passes detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
			})
		})

		context("When only a pyproject.toml file is present", func() {
			it.Before(func() {
				Expect(os.RemoveAll(filepath.Join(workingDir, "x.py"))).To(Succeed())
				Expect(os.WriteFile(filepath.Join(workingDir, "pyproject.toml"), []byte{}, os.ModePerm)).To(Succeed())
			})

			it("passes detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
			})
		})

		context("When only a pixi.lock file is present", func() {
			it.Before(func() {
				Expect(os.RemoveAll(filepath.Join(workingDir, "x.py"))).To(Succeed())
				Expect(os.WriteFile(filepath.Join(workingDir, "pixi.lock"), []byte{}, os.ModePerm)).To(Succeed())
			})

			it("passes detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
			})
		})

		context("When only a uv.lock file is present", func() {
			it.Before(func() {
				Expect(os.RemoveAll(filepath.Join(workingDir, "x.py"))).To(Succeed())
				Expect(os.WriteFile(filepath.Join(workingDir, "uv.lock"), []byte{}, os.ModePerm)).To(Succeed())
			})

			it("passes detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
			})
		})

		context("When no python related files are present", func() {
			it.Before(func() {
				Expect(os.RemoveAll(filepath.Join(workingDir, "x.py"))).To(Succeed())
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("No *.py, environment.yml, pixi.lock, requirements.txt, uv.lock, Pipfile.lock, pyproject.toml, or package-list.txt found")))
			})
		})
	})

	context("failure cases", func() {
		context("when BP_LIVE_RELOAD_ENABLED is set to an invalid value", func() {
			it.Before(func() {
				t.Setenv(pythonstart.LiveReloadEnvName, "not-a-bool")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("failed to parse BP_LIVE_RELOAD_ENABLED value not-a-bool")))
			})
		})

		context("when BP_ENABLE_PACKAGE_MANAGERS is set to an invalid value", func() {
			it.Before(func() {
				t.Setenv(pythonstart.PackageManagersEnvName, "not-a-bool")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("failed to parse BP_ENABLE_PACKAGE_MANAGERS value not-a-bool")))
			})
		})

	})
}
