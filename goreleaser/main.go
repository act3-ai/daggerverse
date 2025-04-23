// A module for running the Goreleaser CLI.
//
// This module aids in building executables and releasing. The bulk of configuration
// should be done in a .goreleaser.yaml file.

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

// Goreleaser represents the `goreleaser` command.
type Goreleaser struct {
	// +private
	Container *dagger.Container
	// +private
	RegistryConfig *dagger.RegistryConfig
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

	gr.Container = dag.Container().
		From(fmt.Sprintf("%s:%s", imageGoReleaser, version)).
		WithWorkdir("/work/src").
		WithMountedDirectory("/work/src", Source)

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
	gr.Container = gr.Container.WithEnvVariable(name, value, opts...)
	return gr
}

// WithSecretVariable adds an env variable containing a secret to the goreleaser container.
//
// This is useful for reusability and readability by not breaking the goreleaser calling chain.
func (gr *Goreleaser) WithSecretVariable(name string, secret *dagger.Secret) *Goreleaser {
	gr.Container = gr.Container.WithSecretVariable(name, secret)
	return gr
}

// Add netrc credentials.
func (gr *Goreleaser) WithNetrc(
	// NETRC credentials
	netrc *dagger.Secret,
) *Goreleaser {
	gr.Container = gr.Container.WithMountedSecret("/root/.netrc", netrc)
	return gr
}

// Add registry credentials.
func (gr *Goreleaser) WithRegistryAuth(
	// registry's hostname
	address string,
	// username in registry
	username string,
	// password or token for registry
	secret *dagger.Secret,
) *Goreleaser {
	gr.RegistryConfig = gr.RegistryConfig.WithRegistryAuth(address, username, secret)
	return gr
}

// Run goreleaser.
//
// Run is a "catch-all" in case functions are not implemented.
func (gr *Goreleaser) Run(
	// arguments and flags, without `goreleaser`.
	args []string,
) *dagger.Container {
	return gr.Container.WithExec(append([]string{"goreleaser"}, args...))
}
