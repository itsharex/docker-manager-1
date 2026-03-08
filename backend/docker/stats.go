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

type DiskUsageSummary struct {
	TotalBytes int64 `json:"totalBytes"`
	UsedBytes  int64 `json:"usedBytes"`
}

func GetDiskUsageSummary() (DiskUsageSummary, error) {
	du, err := Cli.DiskUsage(Ctx(), types.DiskUsageOptions{})
	if err != nil {
		return DiskUsageSummary{}, err
	}

	var volumeBytes int64
	for _, v := range du.Volumes {
		if v.UsageData != nil {
			volumeBytes += v.UsageData.Size
		}
	}

	var imageBytes int64
	for _, img := range du.Images {
		imageBytes += img.Size
	}

	return DiskUsageSummary{
		TotalBytes: 0,
		UsedBytes:  volumeBytes + imageBytes,
	}, nil
}
