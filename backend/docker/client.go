package docker

import (
	"context"
	"github.com/docker/docker/client"
)

var Cli *client.Client

func Init() error {

	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return err
	}

	Cli = cli

	return nil
}

func Ctx() context.Context {
	return context.Background()
}