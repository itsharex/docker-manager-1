package docker

import (
	"context"
	"fmt"

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

	if _, err := Cli.Ping(Ctx()); err != nil {
		return fmt.Errorf("failed to connect Docker daemon (DOCKER_HOST=%q): %w", Cli.DaemonHost(), err)
	}

	return nil
}

func Ctx() context.Context {
	return context.Background()
}
