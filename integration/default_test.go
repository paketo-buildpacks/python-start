package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testDefault(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually
		pack       occam.Pack
		docker     occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when building a default app", func() {
		var (
			image     occam.Image
			container occam.Container
			name      string
			source    string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("builds an oci image with python launch command", func() {
			var err error
			source, err = occam.Source(filepath.Join("testdata", "default_app"))
			Expect(err).NotTo(HaveOccurred())

			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithPullPolicy("never").
				WithBuildpacks(
					cpythonBuildpack,
					buildpack,
				).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			Expect(logs).To(ContainLines(
				MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, buildpackInfo.Buildpack.Name)),
				"  Assigning launch processes:",
				"    web (default): python",
			))

			container, err = docker.Container.Run.
				WithTTY().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				cLogs, err := docker.Container.Logs.Execute(container.ID)
				Expect(err).NotTo(HaveOccurred())
				return cLogs.String()
			}).Should(
				And(
					MatchRegexp(`Python 3\.\d+\.\d+`),
					ContainSubstring(`Type "help", "copyright", "credits" or "license" for more information.`),
				),
			)
		})

		it("builds an oci image with site-packages", func() {
			var err error
			source, err = occam.Source(filepath.Join("testdata", "packages_app"))
			Expect(err).NotTo(HaveOccurred())

			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithPullPolicy("never").
				WithBuildpacks(
					cpythonBuildpack,
					pipBuildpack,
					pipInstallBuildpack,
					buildpack,
				).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			Expect(logs).To(ContainLines(
				MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, buildpackInfo.Buildpack.Name)),
				"  Assigning launch processes:",
				"    web (default): python",
			))

			container, err = docker.Container.Run.
				WithTTY().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				cLogs, err := docker.Container.Logs.Execute(container.ID)
				Expect(err).NotTo(HaveOccurred())
				return cLogs.String()
			}).Should(
				And(
					MatchRegexp(`Python 3\.\d+\.\d+`),
					ContainSubstring(`Type "help", "copyright", "credits" or "license" for more information.`),
				),
			)
		})

		it("builds an oci image with site-packages and module", func() {
			var err error
			source, err = occam.Source(filepath.Join("testdata", "module_app"))
			Expect(err).NotTo(HaveOccurred())

			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithPullPolicy("never").
				WithBuildpacks(
					cpythonBuildpack,
					pipBuildpack,
					pipInstallBuildpack,
					buildpack,
				).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			Expect(logs).To(ContainLines(
				MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, buildpackInfo.Buildpack.Name)),
				"  Assigning launch processes:",
				"    web (default): python",
			))

			container, err = docker.Container.Run.
				WithTTY().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				cLogs, err := docker.Container.Logs.Execute(container.ID)
				Expect(err).NotTo(HaveOccurred())
				return cLogs.String()
			}).Should(
				And(
					MatchRegexp(`Python 3\.\d+\.\d+`),
					ContainSubstring(`Type "help", "copyright", "credits" or "license" for more information.`),
				),
			)
		})

		it("builds an oci image with conda-environment", func() {
			var err error
			source, err = occam.Source(filepath.Join("testdata", "conda_app"))
			Expect(err).NotTo(HaveOccurred())

			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithPullPolicy("never").
				WithBuildpacks(
					minicondaBuildpack,
					condaEnvUpdateBuildpack,
					buildpack,
				).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			Expect(logs).To(ContainLines(
				MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, buildpackInfo.Buildpack.Name)),
				"  Assigning launch processes:",
				"    web (default): python",
			))

			container, err = docker.Container.Run.
				WithTTY().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				cLogs, err := docker.Container.Logs.Execute(container.ID)
				Expect(err).NotTo(HaveOccurred())
				return cLogs.String()
			}).Should(
				And(
					MatchRegexp(`Python 2.7\.\d+`),
					ContainSubstring(`Type "help", "copyright", "credits" or "license" for more information.`),
				),
			)
		})

		context("when building an app with poetry (dependency management only)", func() {
			var container2 occam.Container
			it.After(func() {
				Expect(docker.Container.Remove.Execute(container2.ID)).To(Succeed())

			})

			it("builds an oci image with poetry on PATH", func() {
				var err error
				source, err = occam.Source(filepath.Join("testdata", "poetry"))
				Expect(err).NotTo(HaveOccurred())

				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithPullPolicy("never").
					WithBuildpacks(
						cpythonBuildpack,
						pipBuildpack,
						poetryBuildpack,
						poetryInstallBuildpack,
						buildpack,
					).
					Execute(name, source)
				Expect(err).NotTo(HaveOccurred(), logs.String())

				Expect(logs).To(ContainLines(
					MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, buildpackInfo.Buildpack.Name)),
					"  Assigning launch processes:",
					"    web (default): python",
				))

				container, err = docker.Container.Run.
					WithTTY().
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() string {
					cLogs, err := docker.Container.Logs.Execute(container.ID)
					Expect(err).NotTo(HaveOccurred())
					return cLogs.String()
				}).Should(
					And(
						MatchRegexp(`Python 3\.\d+\.\d+`),
						ContainSubstring(`Type "help", "copyright", "credits" or "license" for more information.`),
					),
				)

				container2, err = docker.Container.Run.
					WithTTY().
					WithEntrypoint("launcher").
					WithCommand("poetry --no-ansi --version"). // Use the no-ansi flag to disable color output - required for regex to match
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() string {
					cLogs, err := docker.Container.Logs.Execute(container2.ID)
					Expect(err).NotTo(HaveOccurred())
					return cLogs.String()
				}).Should(MatchRegexp(`Poetry \(version \d+\.\d+\.\d+\)`))
			})
		})
	})
}
