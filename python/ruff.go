package main

import (
	"context"
	"dagger/python/internal/dagger"
	"fmt"
)

type RuffCheckFormat string

const (
	text       RuffCheckFormat = "text"
	stdout     RuffCheckFormat = "stdout"
	jsonFormat RuffCheckFormat = "json"
	gitlab     RuffCheckFormat = "gitlab"
	concise    RuffCheckFormat = "concise"
	full       RuffCheckFormat = "full"
	junit      RuffCheckFormat = "junit"
	pylint     RuffCheckFormat = "pylint"
)

// Return the result of running ruff check
func (python *Python) RuffCheck(ctx context.Context,
	//output format of ruff lint check, valid values: concise, full, json, json-lines, junit, grouped, github, gitlab
	// +optional
	// +default="full"
	ruffCheckFormat RuffCheckFormat,
) *dagger.File {

	outputFile := fmt.Sprintf("ruff-check.%s", ruffCheckFormat)

	c := python.Container()
	// Run the Ruff linter with the provided output format
	// The output format is passed to the --ruff-check-format flag
	ruffResults := c.WithExec(
		[]string{
			"uv",
			"run",
			"--with=ruff",
			"ruff",
			"check", ".",
			"--output-file", outputFile,
			"--output-format", string(ruffCheckFormat)})

	// Return the formatted output of the Ruff check as a string
	// The output could include details about any code formatting issues.
	return ruffResults.File(outputFile)
}

// Return the result of running ruff format
func (python *Python) RuffFormat(ctx context.Context) (*dagger.File, error) {

	c := python.Container().
		WithExec(
			[]string{
				"uv",
				"run",
				"--with=ruff",
				"ruff",
				"format",
				"--check",
				"--diff", "."})

	// Return the formatted output of Ruff format as a string
	// The output could include details about any code formatting issues.
	results, err := c.Stdout(ctx)
	if err != nil {
		return nil, err
	}
	return dag.Directory().WithNewFile("ruff-format.txt", results).File("ruff-format.txt"), nil
}
