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
	)

	it.Before(func() {
		var err error
		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		Expect(os.WriteFile(filepath.Join(workingDir, "x.py"), []byte{}, os.ModePerm)).To(Succeed())

		detect = pythonstart.Detect()
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("detection phase", func() {
		it("detects", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "cpython",
						Metadata: pythonstart.BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "site-packages",
						Metadata: pythonstart.BuildPlanMetadata{
							Launch: true,
						},
					},
				},
				Or: []packit.BuildPlan{
					{
						Provides: []packit.BuildPlanProvision{},
						Requires: []packit.BuildPlanRequirement{
							{
								Name: "conda-environment",
								Metadata: pythonstart.BuildPlanMetadata{
									Launch: true,
								},
							},
						},
					},
					{
						Provides: []packit.BuildPlanProvision{},
						Requires: []packit.BuildPlanRequirement{
							{
								Name: "cpython",
								Metadata: pythonstart.BuildPlanMetadata{
									Launch: true,
								},
							},
							{
								Name: "poetry",
								Metadata: pythonstart.BuildPlanMetadata{
									Launch: true,
								},
							},
							{
								Name: "poetry-venv",
								Metadata: pythonstart.BuildPlanMetadata{
									Launch: true,
								},
							},
						},
					},
					{
						Provides: []packit.BuildPlanProvision{},
						Requires: []packit.BuildPlanRequirement{
							{
								Name: "cpython",
								Metadata: pythonstart.BuildPlanMetadata{
									Launch: true,
								},
							},
						},
					},
				},
			}))
		})

		context("when BP_LIVE_RELOAD_ENABLED=true in the build environment", func() {
			it.Before(func() {
				os.Setenv("BP_LIVE_RELOAD_ENABLED", "true")
			})

			it.After(func() {
				os.Unsetenv("BP_LIVE_RELOAD_ENABLED")
			})

			it("requires watchexec at launch", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan).To(Equal(packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{},
					Requires: []packit.BuildPlanRequirement{
						{
							Name: "cpython",
							Metadata: pythonstart.BuildPlanMetadata{
								Launch: true,
							},
						},
						{
							Name: "site-packages",
							Metadata: pythonstart.BuildPlanMetadata{
								Launch: true,
							},
						},
						{
							Name: "watchexec",
							Metadata: pythonstart.BuildPlanMetadata{
								Launch: true,
							},
						},
					},
					Or: []packit.BuildPlan{
						{
							Provides: []packit.BuildPlanProvision{},
							Requires: []packit.BuildPlanRequirement{
								{
									Name: "conda-environment",
									Metadata: pythonstart.BuildPlanMetadata{
										Launch: true,
									},
								},
								{
									Name: "watchexec",
									Metadata: pythonstart.BuildPlanMetadata{
										Launch: true,
									},
								},
							},
						},
						{
							Provides: []packit.BuildPlanProvision{},
							Requires: []packit.BuildPlanRequirement{
								{
									Name: "cpython",
									Metadata: pythonstart.BuildPlanMetadata{
										Launch: true,
									},
								},
								{
									Name: "poetry",
									Metadata: pythonstart.BuildPlanMetadata{
										Launch: true,
									},
								},
								{
									Name: "poetry-venv",
									Metadata: pythonstart.BuildPlanMetadata{
										Launch: true,
									},
								},
								{
									Name: "watchexec",
									Metadata: pythonstart.BuildPlanMetadata{
										Launch: true,
									},
								},
							},
						},
						{
							Provides: []packit.BuildPlanProvision{},
							Requires: []packit.BuildPlanRequirement{
								{
									Name: "cpython",
									Metadata: pythonstart.BuildPlanMetadata{
										Launch: true,
									},
								},
								{
									Name: "watchexec",
									Metadata: pythonstart.BuildPlanMetadata{
										Launch: true,
									},
								},
							},
						},
					},
				}))
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

		context("When no python related files are present", func() {
			it.Before(func() {
				Expect(os.RemoveAll(filepath.Join(workingDir, "x.py"))).To(Succeed())
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("No *.py, environment.yml, pyproject.toml, or package-list.txt found")))
			})
		})
	})

	context("failure cases", func() {
		context("when BP_LIVE_RELOAD_ENABLED is set to an invalid value", func() {
			it.Before(func() {
				os.Setenv("BP_LIVE_RELOAD_ENABLED", "not-a-bool")
			})

			it.After(func() {
				os.Unsetenv("BP_LIVE_RELOAD_ENABLED")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("failed to parse BP_LIVE_RELOAD_ENABLED value not-a-bool")))
			})
		})

	})
}
