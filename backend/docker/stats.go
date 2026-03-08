package docker

import (
	"github.com/docker/docker/api/types"
)

func GetSystemInfo() (types.Info, error) {
	return Cli.Info(Ctx())
}

func GetVersion() (types.Version, error) {
	return Cli.ServerVersion(Ctx())
}
