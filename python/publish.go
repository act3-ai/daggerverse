package main

import (
	"context"
	"dagger/python/internal/dagger"
)

// publish python package/wheel
func (python *Python) Publish(ctx context.Context,
	publishUrl string,
	username string,
	password *dagger.Secret,
) (*dagger.Container, error) {

	c := python.Container().
		WithEnvVariable("UV_PUBLISH_CHECK_URL", publishUrl+"/simple").
		WithEnvVariable("UV_PUBLISH_URL", publishUrl).
		WithEnvVariable("UV_PUBLISH_USERNAME", username).
		WithSecretVariable("UV_PUBLISH_PASSWORD", password).
		WithExec(
			[]string{
				"uv",
				"build",
			}).
		WithExec(
			[]string{
				"uv",
				"publish",
			})

	return c, nil
}
