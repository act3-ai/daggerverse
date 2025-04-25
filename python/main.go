// A generated module for Python testing functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.

package main

import (
	"context"
	"dagger/python/internal/dagger"
	"fmt"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// TODO add renovate to this
const uvImageDefault = "ghcr.io/astral-sh/uv:0.6.16-debian"

type Python struct {
	// +private
	Base *dagger.Container

	// +private
	Source *dagger.Directory

	// +private
	Netrc *dagger.Secret

	// +private
	SyncArgs []string
}

func New(
	// base development container
	// +optional
	base *dagger.Container,

	// top-level source code directory
	// +defaultPath="/"
	src *dagger.Directory,

	// .netrc file for private modules can be passed as env var or file --netrc env:var_name, file:/filepath/.netrc
	// +optional
	netrc *dagger.Secret,

	// extra arguments for uv sync command
	// +optional
	syncArgs []string,
) *Python {
	if base == nil {
		base = dag.Container().From(uvImageDefault)
	}

	if syncArgs == nil {
		syncArgs = []string{
			"--frozen",
			"--all-extras",
		}
	}

	return &Python{
		Base:     base,
		Source:   src,
		Netrc:    netrc,
		SyncArgs: syncArgs,
	}
}

// base UV container (with caching, source, and credentials injected)
func (python *Python) UV() *dagger.Container {
	return python.Base.
		WithDirectory("/app", python.Source).
		WithWorkdir("/app").
		WithMountedCache("/root/.cache/uv", dag.CacheVolume("uv-cache")).
		WithEnvVariable("UV_NATIVE_TLS", "true").
		WithEnvVariable("UV_CACHE_DIR", "/root/.cache/uv"). // This is the default location for the UV_CACHE_DIR but we set it just to be safe.
		With(func(c *dagger.Container) *dagger.Container {
			if python.Netrc != nil {
				return c.WithMountedSecret("/root/.netrc", python.Netrc)
			}
			return c
		})
}

// build dev dependencies first before running test
func (python *Python) Container() *dagger.Container {
	return python.UV().
		WithExec(
			append(
				[]string{"uv", "sync"},
				python.SyncArgs...,
			),
		)
}

// add an environment variable to the base container
func (python *Python) WithEnvVariable(name, value string) *Python {
	python.Base = python.Base.WithEnvVariable(name, value)
	return python
}

// check that the lockfile is in sync with pyproject.toml
func (python *Python) CheckLock(ctx context.Context) (string, error) {
	return python.UV().
		WithExec([]string{"uv", "lock", "--check"}).
		Stdout(ctx)
}

// Return the result of all lint checks
func (python *Python) Lint(ctx context.Context,
	//output format of ruff check, valid values: full, concise, gitlab, json, stdout, text
	// +optional
	// +default="full"
	ruffCheckFormat RuffCheckFormat,

	// +optional
	skip []string,
) (*dagger.Directory, error) {
	checks := map[string]func(context.Context) (*dagger.File, error){
		"ruff-check": func(ctx context.Context) (*dagger.File, error) {
			return python.RuffCheck(ctx, ruffCheckFormat), nil
		},
		"ruff-format": func(ctx context.Context) (*dagger.File, error) {
			return python.RuffCheck(ctx, ruffCheckFormat), nil
		},
		"mypy":    python.Mypy,
		"pylint":  python.Pylint,
		"pyright": python.Pyright,
	}

	for _, check := range skip {
		delete(checks, check)
	}

	p := pool.NewWithResults[*dagger.File]().WithContext(ctx).WithMaxGoroutines(3) //.WithCollectErrored()
	for name, check := range checks {
		p.Go(func(ctx context.Context) (*dagger.File, error) {
			ctx, span := Tracer().Start(ctx, name)
			defer span.End()
			return check(ctx)
		})
	}

	// Wait for all goroutines to finish
	files, err := p.Wait()
	if err != nil {
		return nil, err
	}

	//create new directory with result files
	return dag.Directory().WithFiles("/", files), nil
}

// Return the result of running all tests(lint and unit test)
func (python *Python) Test(ctx context.Context,
	// Output format of ruff check, valid values: gitlab, json, table, stdout, text
	// +optional
	// +default="full"
	ruffCheckFormat RuffCheckFormat,
	// unit test directoy
	// +optional
	// +default="test"
	unitTestDir string,

	// +optional
	skip []string,
) (*dagger.Directory, error) {

	var combinedErr []string // To aggregate errors

	// Run Lint
	lintResultsDirectory, lintErr := python.Lint(ctx, ruffCheckFormat, skip)

	if lintErr != nil {
		combinedErr = append(combinedErr, "Lint Error: "+lintErr.Error())
	}

	// run unit test
	unitTestResults, err := python.UnitTest(ctx, unitTestDir)
	if err != nil {
		return nil, err
	}

	if unitTestResults == nil {
		combinedErr = append(combinedErr, "Unit Test Error")
	}

	// If there are any errors, combine them into a single error
	if len(combinedErr) > 0 {
		return nil, fmt.Errorf(strings.Join(combinedErr, "\n"))
	}

	testResultsDir := dag.Directory().WithDirectory("lint-results", lintResultsDirectory).WithDirectory("unit-test-results", unitTestResults)

	return testResultsDir, nil
}
