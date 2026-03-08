package docker

import (
	"github.com/docker/docker/api/types"
)

func ListNetworks() ([]types.NetworkResource, error) {
	return Cli.NetworkList(
		Ctx(),
		types.NetworkListOptions{},
	)
}

func RemoveNetwork(id string) error {
	return Cli.NetworkRemove(
		Ctx(),
		id,
	)
}
