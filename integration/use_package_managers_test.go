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

func testUsePythonPackageManager(t *testing.T, context spec.G, it spec.S) {
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
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		context("successful", func() {
			it.After(func() {
				Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
				Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
				Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			})
			it("builds an oci image with conda-environment", func() {
				var err error
				source, err = sourceWithCode(filepath.Join("testdata", "conda_app"))
				Expect(err).NotTo(HaveOccurred())

				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithPullPolicy("never").
					WithBuildpacks(
						pythonPackageManagersInstallBuildpack,
						pythonPackageManagersRunBuildpack,
						buildpack,
					).
					WithEnv(map[string]string{
						"BP_ENABLE_PACKAGE_MANAGERS": "true",
					}).
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
		})

		context("error case", func() {
			it("fails to build an oci image with conda-environment using old buildpacks", func() {
				var err error
				source, err = sourceWithCode(filepath.Join("testdata", "conda_app"))
				Expect(err).NotTo(HaveOccurred())

				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithPullPolicy("never").
					WithBuildpacks(
						minicondaBuildpack,
						condaEnvUpdateBuildpack,
						buildpack,
					).
					WithEnv(map[string]string{
						"BP_ENABLE_PACKAGE_MANAGERS": "true",
					}).
					Execute(name, source)
				Expect(err).To(HaveOccurred(), logs.String())
			})
		})
	})
}
