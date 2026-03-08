package docker

import (
	"github.com/docker/docker/api/types"
)

func ListImages() ([]types.ImageSummary, error) {
	return Cli.ImageList(
		Ctx(),
		types.ImageListOptions{All: true},
	)
}

func RemoveImage(id string) error {
	_, err := Cli.ImageRemove(
		Ctx(),
		id,
		types.ImageRemoveOptions{Force: true},
	)
	return err
}
