package main

import (
	"context"
	"dagger/python/internal/dagger"
)

// Return the result of running Pyright
func (python *Python) Pyright(ctx context.Context) (*dagger.File, error) {

	c := python.Container().
		WithExec(
			[]string{
				"uv",
				"run",
				"--with=pyright",
				"pyright",
				".",
			})

	results, err := c.Stdout(ctx)
	if err != nil {
		return nil, err
	}
	return dag.Directory().WithNewFile("pyright.txt", results).File("pyright.txt"), nil
}
