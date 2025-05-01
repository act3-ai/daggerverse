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
	Version string,

	// Disable mounting cache volumes.
	//
	// +optional
	disableCache bool,
) *Goreleaser {
	gr := &Goreleaser{
		Container:      defaultContainer(Source, Version),
		RegistryConfig: dag.RegistryConfig(),
	}

	if !disableCache {
		gr = gr.WithGoModuleCache(dag.CacheVolume("go-mod"), nil, "").
			WithBuildCache(dag.CacheVolume("go-build"), nil, "")
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
	gr.Container = gr.Container.WithEnvVariable(
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

// Mount a cache volume for Go module cache.
func (gr *Goreleaser) WithGoModuleCache(
	cache *dagger.CacheVolume,

	// Identifier of the directory to use as the cache volume's root.
	//
	// +optional
	source *dagger.Directory,

	// Sharing mode of the cache volume.
	//
	// +optional
	sharing dagger.CacheSharingMode,
) *Goreleaser {
	gr.Container = gr.Container.WithMountedCache(
		"/go/pkg/mod",
		cache,
		dagger.ContainerWithMountedCacheOpts{
			Source:  source,
			Sharing: sharing,
		},
	)

	return gr
}

// Mount a cache volume for Go build cache.
func (gr *Goreleaser) WithBuildCache(
	cache *dagger.CacheVolume,

	// Identifier of the directory to use as the cache volume's root.
	//
	// +optional
	source *dagger.Directory,

	// Sharing mode of the cache volume.
	//
	// +optional
	sharing dagger.CacheSharingMode,
) *Goreleaser {
	gr.Container = gr.Container.WithMountedCache(
		"/root/.cache/go-build",
		cache,
		dagger.ContainerWithMountedCacheOpts{
			Source:  source,
			Sharing: sharing,
		},
	)

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

// defaultContainer constructs a minimal container containing a source git repository.
func defaultContainer(source *dagger.Directory, version string) *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("%s:%s", imageGoReleaser, version)).
		WithWorkdir("/work/src").
		WithMountedDirectory("/work/src", source).
		With(func(r *dagger.Container) *dagger.Container {
			// inherit from host, overriden by WithEnvVariable
			val, ok := os.LookupEnv(envGOMAXPROCS)
			if ok {
				r = r.WithEnvVariable(envGOMAXPROCS, val)
			}
			return r
		}).
		With(func(r *dagger.Container) *dagger.Container {
			// inherit from host, overriden by WithEnvVariable
			val, ok := os.LookupEnv(envGOMEMLIMIT)
			if ok {
				r = r.WithEnvVariable(envGOMEMLIMIT, val)
			}
			return r
		})
}
