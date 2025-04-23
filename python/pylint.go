package main

import (
	"context"
	"dagger/python/internal/dagger"
)

// Return the result of running pylint
func (python *Python) Pylint(ctx context.Context) (*dagger.File, error) {

	c := python.Container().
		WithExec(
			[]string{
				"uv",
				"run",
				"--with=pylint",
				"pylint",
				"--recursive=y",
				// "--reports=y",
				".",
			})

	results, err := c.Stdout(ctx)
	if err != nil {
		return nil, err
	}
	return dag.Directory().WithNewFile("pylint.txt", results).File("pylint.txt"), nil
}
