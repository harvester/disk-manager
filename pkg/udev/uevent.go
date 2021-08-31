package udev

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/pilebones/go-udev/netlink"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/longhorn/node-disk-manager/pkg/apis/longhorn.io/v1beta1"
	"github.com/longhorn/node-disk-manager/pkg/block"
	"github.com/longhorn/node-disk-manager/pkg/controller/blockdevice"
	"github.com/longhorn/node-disk-manager/pkg/filter"
	ctldiskv1 "github.com/longhorn/node-disk-manager/pkg/generated/controllers/longhorn.io/v1beta1"
	"github.com/longhorn/node-disk-manager/pkg/option"
)

const (
	defaultDuration time.Duration = 1
)

type Udev struct {
	namespace  string
	nodeName   string
	startOnce  sync.Once
	controller *blockdevice.Controller
}

func NewUdev(block *block.Info, blockdevices ctldiskv1.BlockDeviceController, opt *option.Option, filters []*filter.Filter) *Udev {
	controller := &blockdevice.Controller{
		BlockInfo:        block,
		Blockdevices:     blockdevices,
		BlockdeviceCache: blockdevices.Cache(),
		Filters:          filters,
	}
	return &Udev{
		startOnce:  sync.Once{},
		namespace:  opt.Namespace,
		nodeName:   opt.NodeName,
		controller: controller,
	}
}

func (u *Udev) Monitor(ctx context.Context) {
	u.startOnce.Do(func() {
		u.monitor(ctx)
	})
}

func (u *Udev) monitor(ctx context.Context) {
	logrus.Infoln("Start monitoring udev processed events")

	matcher, err := getOptionalMatcher(nil)
	if err != nil {
		logrus.Fatalf("Failed to get udev config, error: %s", err.Error())
	}

	conn := new(netlink.UEventConn)
	if err := conn.Connect(netlink.UdevEvent); err != nil {
		logrus.Fatalf("Unable to connect to Netlink Kobject UEvent socket, error: %s", err.Error())
	}
	defer conn.Close()

	uqueue := make(chan netlink.UEvent)
	errors := make(chan error)
	quit := conn.Monitor(uqueue, errors, matcher)

	// Handling message from udev queue
	for {
		select {
		case uevent := <-uqueue:
			u.ActionHandler(uevent)
		case err := <-errors:
			logrus.Errorf("failed to parse udev event, error: %s", err.Error())
		case <-ctx.Done():
			close(quit)
			return
		}
	}
}

func (u *Udev) ActionHandler(uevent netlink.UEvent) {
	udevDevice := InitUdevDevice(uevent.Env)
	if !udevDevice.IsDisk() && !udevDevice.IsPartition() {
		return
	}

	disk := &block.Disk{}
	part := &block.Partition{}
	blockdevice := &v1beta1.BlockDevice{}
	if udevDevice.IsDisk() {
		disk = u.controller.BlockInfo.GetDiskByDevPath(udevDevice.GetShortName())
		blockdevice = udevDevice.DeviceInfoFromUdevDisk(disk, u.nodeName, u.namespace)
	} else {
		parentPath, err := block.GetParentDevName(udevDevice.GetDevName())
		logrus.Infof("debug: parent path %s", parentPath)
		if err != nil {
			logrus.Errorf("failed to get parent dev name, %s", err.Error())
		}
		part = u.controller.BlockInfo.GetPartitionByDevPath(parentPath, udevDevice.GetDevName())
		disk = part.Disk
		blockdevice = udevDevice.DeviceInfoFromUdevPartition(part, u.nodeName, u.namespace)
	}

	if u.controller.ApplyFilter(disk) {
		return
	}

	switch uevent.Action {
	//case netlink.ADD:
	//	u.AddBlockDevice(blockdevice, defaultDuration)
	case netlink.REMOVE:
		if udevDevice.IsDisk() {
			u.RemoveBlockDevice(blockdevice, defaultDuration)
		}
	}
}

// AddBlockDevice add new block device and partitions by watching the udev add action
func (u *Udev) AddBlockDevice(device *v1beta1.BlockDevice, duration time.Duration) {
	if duration > defaultDuration {
		time.Sleep(duration)
	}
	logrus.Debugf("uevent add block deivce %s", device.Spec.DevPath)

	bdList, err := u.controller.BlockdeviceCache.List(u.namespace, labels.Everything())
	if err != nil {
		logrus.Errorf("Failed to add block device via udev event, error: %s, retry in %s", err.Error(), 2*duration)
		u.AddBlockDevice(device, 2*duration)
	}

	if err := u.controller.SaveBlockDevice(device, bdList); err != nil {
		logrus.Errorf("failed to save block device %s, error: %s", device.Name, err.Error())
		//u.AddBlockDevice(device, 2*defaultDuration)
	}
}

// RemoveBlockDevice will set the existing block device to detached state
func (u *Udev) RemoveBlockDevice(device *v1beta1.BlockDevice, duration time.Duration) {
	if duration > defaultDuration {
		time.Sleep(duration)
	}

	logrus.Debugf("uevent remove block deivce %s", device.Spec.DevPath)

	err := u.controller.Blockdevices.Delete(u.namespace, device.Name, &metav1.DeleteOptions{})
	if err != nil && errors.IsNotFound(err) {
		logrus.Errorf("failed to delete block device, %s is not found", device.Name)
	} else if err != nil {
		logrus.Errorf("faield to delete the block device %s, error: %s", device.Name, err.Error())
		u.RemoveBlockDevice(device, 2*duration)
	}
}

// getOptionalMatcher Parse and load config file which contains rules for matching
func getOptionalMatcher(filePath *string) (matcher netlink.Matcher, err error) {
	if filePath == nil || *filePath == "" {
		return nil, nil
	}

	stream, err := ioutil.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}

	if stream == nil {
		return nil, fmt.Errorf("empty, no rules provided in \"%s\", err: %w", *filePath, err)
	}

	var rules netlink.RuleDefinitions
	if err := json.Unmarshal(stream, &rules); err != nil {
		return nil, fmt.Errorf("wrong rule syntax, err: %w", err)
	}

	return &rules, nil
}
