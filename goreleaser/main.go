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
	Ctr *dagger.Container
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

	gr.Ctr = dag.Container().
		From(fmt.Sprintf("%s:%s", imageGoReleaser, version)).
		WithWorkdir("/work/src").
		WithMountedDirectory("/work/src", Source)

	// inherit from host, overriden by WithEnvVariable
	val, ok := os.LookupEnv(envGOMAXPROCS)
	if ok {
		gr = gr.WithEnvVariable(envGOMAXPROCS, val, false)
	}
	val, ok = os.LookupEnv(envGOMEMLIMIT)
	if ok {
		gr = gr.WithEnvVariable(envGOMEMLIMIT, val, false)
	}

	return gr
}

// WithEnvVariable adds an environment variable to the goreleaser container.
//
// This is useful for reusability and readability by not breaking the goreleaser calling chain.
func (gr *Goreleaser) WithEnvVariable(
	// The name of the environment variable (e.g., "HOST").
	name string,
	// The value of the environment variable (e.g., "localhost").
	value string,
	// Replace `${VAR}` or $VAR in the value according to the current environment
	// variables defined in the container (e.g., "/opt/bin:$PATH").
	//
	// +optional
	expand bool,
) *Goreleaser {
	gr.Ctr = gr.Ctr.WithEnvVariable(
		name,
		value,
		dagger.ContainerWithEnvVariableOpts{
			Expand: expand,
		},
	)
	return gr
}

// WithSecretVariable adds an env variable containing a secret to the goreleaser container.
//
// This is useful for reusability and readability by not breaking the goreleaser calling chain.
func (gr *Goreleaser) WithSecretVariable(
	// The name of the environment variable containing a secret (e.g., "PASSWORD").
	name string,
	// The value of the environment variable containing a secret.
	secret *dagger.Secret,
) *Goreleaser {
	gr.Ctr = gr.Ctr.WithSecretVariable(name, secret)
	return gr
}

// Add netrc credentials.
func (gr *Goreleaser) WithNetrc(
	// NETRC credentials
	netrc *dagger.Secret,
) *Goreleaser {
	gr.Ctr = gr.Ctr.WithMountedSecret("/root/.netrc", netrc)
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
	return gr.Ctr.WithExec(append([]string{"goreleaser"}, args...))
}

// Fetch the goreleaser container in its current state. All modifications are preserved, e.g. environment variables.
func (gr *Goreleaser) Container() *dagger.Container {
	return gr.Ctr
}
