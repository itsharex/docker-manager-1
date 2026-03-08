package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/system"
)

func GetSystemInfo() (system.Info, error) {
	return Cli.Info(Ctx())
}

func GetVersion() (types.Version, error) {
	return Cli.ServerVersion(Ctx())
}
