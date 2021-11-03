package filter

import (
	"path/filepath"

	"github.com/harvester/node-disk-manager/pkg/block"
	"github.com/harvester/node-disk-manager/pkg/util"
	"github.com/sirupsen/logrus"
)

const (
	devicePathFilterName = "device path filter"
)

// partDevicePathFilter filters devices based on given device path patterns and
// their parent device filters.
type partDevicePathFilter struct {
	filter *diskDevicePathFilter
}

// diskDevicePathFilter filters devices based on given device path patterns
type diskDevicePathFilter struct {
	devicePaths []string
}

func RegisterDevicePathFilter(filters ...string) *Filter {
	f := &diskDevicePathFilter{}
	for _, filter := range filters {
		if filter != "" {
			f.devicePaths = append(f.devicePaths, filter)
		}
	}
	return &Filter{
		Name:       devicePathFilterName,
		PartFilter: &partDevicePathFilter{filter: f},
		DiskFilter: f,
	}
}

// Match returns true if given device path matches the pattern.
func (f *partDevicePathFilter) Match(part *block.Partition) bool {
	devPath := util.GetFullDevPath(part.Name)
	if devPath == "" {
		return false
	}
	return match(devPath, f.filter.devicePaths)
}

// Match returns true if given device path matches the pattern.
func (f *diskDevicePathFilter) Match(disk *block.Disk) bool {
	devPath := util.GetFullDevPath(disk.Name)
	if devPath == "" {
		return false
	}
	return match(devPath, f.devicePaths)
}

func match(devPath string, patterns []string) bool {
	for _, pattern := range patterns {
		if pattern == "" || devPath == "" {
			return false
		}
		ok, err := filepath.Match(pattern, devPath)
		if err != nil {
			logrus.Errorf("failed to perform device path matching on disk %s for pattern %s: %s", devPath, pattern, err.Error())
			return false
		}
		if ok {
			return true
		}
	}
	return false
}
