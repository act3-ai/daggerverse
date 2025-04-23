package main

import (
	"context"
	"dagger/python/internal/dagger"
)

// // build mypy container
// func (python *Python) mypyContainer(ctx context.Context, source *dagger.Directory) *dagger.Container {

// 	// create container and install mypy
// 	return python.MypyContainer.
// 		WithDirectory("/app", source).
// 		WithWorkdir("/app").
// 		WithExec([]string{"uv", "tool", "install", "mypy"})

// }

// Return the result of running mypy
func (python *Python) Mypy(ctx context.Context) (*dagger.File, error) {

	c := python.Container().
		WithExec(
			[]string{
				"uv",
				"run",
				"--with=mypy",
				"mypy",
				"--junit-xml",
				"mypy-junit.xml",
				".",
			},
		)

	// Return the formatted output of the mypy check in a file
	return c.File("mypy-junit.xml"), nil
}
