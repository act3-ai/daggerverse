// A generated module for Tests functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/tests/internal/dagger"
	"errors"
)

type Tests struct{}

// Run all tests.
func (m *Tests) All(ctx context.Context) error {
	var errs []error

	// TODO: conc pkg will be useful once this grows, be sure to limit goroutines.
	errs = append(errs, m.TestBuildAll(ctx))
	errs = append(errs, m.TestBuildPlatform(ctx))

	return errors.Join(errs...)
}

// Test build for all platforms defined in goreleaser config.
func (m *Tests) TestBuildAll(ctx context.Context) error {

	_, err := dag.Goreleaser(testDir()).
		Build().
		All().
		Stdout(ctx)

	return err
}

// Test build for a specific platform.
func (m *Tests) TestBuildPlatform(ctx context.Context) error {

	_, err := dag.Goreleaser(testDir()).
		Build().
		All().
		Stdout(ctx)

	return err
}

// testDir provides a git repository used for testing.
func testDir() *dagger.Directory {
	const (
		// TODO: Relying on an outside repository may not a great practice.
		// Note only is this repo external (relative to this repo), but
		// a larger project takes longer to build.
		testGitRepo    = "https://github.com/act3-ai/data-tool.git"
		testGitRepoTag = "v1.15.33"
	)

	return dagger.Connect().
		Git(testGitRepo).
		Tag(testGitRepoTag).
		Tree()

}
