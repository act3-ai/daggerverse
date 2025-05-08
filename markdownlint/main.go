// A generated module for Markdownlint functions
//
// Package inspired by https://github.com/sagikazarmark/daggerverse/blob/main/golangci-lint/main.go.

package main

import (
	"dagger/markdownlint/internal/dagger"
	"fmt"
)

// defaultImageRepository is used when no image is specified.
const defaultImageRepository = "docker.io/davidanson/markdownlint-cli2"

type Markdownlint struct {
	Container *dagger.Container
}

func New(
	// Custom container to use as a base container. Must have 'markdownlint-cli2' available on PATH.
	// +optional
	Container *dagger.Container,

	// Version (image tag) to use as a markdownlint-cli2 binary source.
	// +optional
	// +default="latest"
	Version string,
) *Markdownlint {
	if Container == nil {
		Container = defaultContainer(Version)
	}

	return &Markdownlint{
		Container: Container,
	}
}

// Run markdownlint-cli2.
func (m *Markdownlint) Run(
	// Source directory containing markdown files to be linted.
	source *dagger.Directory,

	// Custom configuration file.
	// +optional
	config *dagger.File,

	// Additional arguments to pass to markdownlint-cli2.
	// +optional
	extraArgs []string,
) *dagger.Container {
	args := []string{"markdownlint-cli2"}

	return m.Container.
		WithWorkdir("/work/src").
		WithMountedDirectory(".", source).
		With(func(c *dagger.Container) *dagger.Container {
			if config != nil {
				c = c.WithMountedFile("/work/config", config)
				args = append(args, "--config", "/work/config")
			}
			return c
		}).
		WithExec(args)
}

func defaultContainer(version string) *dagger.Container {
	binary := dag.Container().
		From(fmt.Sprintf("%s:%s", defaultImageRepository, version)).
		File("/usr/local/bin/markdownlint-cli2")

	return dag.Container().
		WithFile("/usr/local/bin/markdownlint-cli2", binary, dagger.ContainerWithFileOpts{Permissions: 0755})
}
