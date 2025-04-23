// A module for running the Goreleaser CLI.
//
// This module aids in building executables and releasing. The bulk of configuration
// should be done in .goreleaser.yaml.

package main

import (
	"dagger/goreleaser/internal/dagger"
	"fmt"
	"os"
)

// environment variable names
const (
	envGOMAXPROCS = "GOMAXPROCS"
	envGOMEMLIMIT = "GOMEMLIMIT"
	envGOOS       = "GOOS"
	envGOARCH     = "GOARCH"
	envGOARM      = "GOARM"
)

const (
	imageGoReleaser = "ghcr.io/goreleaser/goreleaser" // defaults to "latest"
)

type Goreleaser struct {
	// +private
	container *dagger.Container
}

func New(
	// Git repository source.
	Source *dagger.Directory,

	// Version (image tag) to use as a goreleaser binary source.
	// +optional
	// +default="latest"
	version string,
) *Goreleaser {
	gr := &Goreleaser{}

	gr.container = dag.Container().
		From(fmt.Sprintf("%s:%s", imageGoReleaser, version))

	// inherit from host, overriden by WithEnvVariable
	val, ok := os.LookupEnv(envGOMAXPROCS)
	if ok {
		gr = gr.WithEnvVariable(envGOMAXPROCS, val)
	}
	val, ok = os.LookupEnv(envGOMEMLIMIT)
	if ok {
		gr = gr.WithEnvVariable(envGOMEMLIMIT, val)
	}

	return gr
}

// WithEnvVariable adds an environment variable to the goreleaser container.
//
// This is useful for reusability and readability by not breaking the goreleaser calling chain.
func (gr *Goreleaser) WithEnvVariable(name, value string, opts ...dagger.ContainerWithEnvVariableOpts) *Goreleaser {
	gr.container = gr.container.WithEnvVariable(name, value, opts...)
	return gr
}

// WithSecretVariable adds an env variable containing a secret to the goreleaser container.
//
// This is useful for reusability and readability by not breaking the goreleaser calling chain.
func (gr *Goreleaser) WithSecretVariable(name string, secret *dagger.Secret) *Goreleaser {
	gr.container = gr.container.WithSecretVariable(name, secret)
	return gr
}

// Run goreleaser.
func (gr *Goreleaser) Run(
	// arguments and flags, without `goreleaser`.
	args []string,
) *dagger.Container {
	return gr.container.WithExec(append([]string{"goreleaser"}, args...))
}
