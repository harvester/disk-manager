package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kevinburke/ssh_config"
	"github.com/melbahja/goph"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	diskv1 "github.com/harvester/node-disk-manager/pkg/apis/harvesterhci.io/v1beta1"
	clientset "github.com/harvester/node-disk-manager/pkg/generated/clientset/versioned"
)

/*
 * We have some assumption for the hotplug test:
 * 1. we will reuse the disk that is added on the initinal operation of ci test
 * 2. we use virsh command to remove disk/add back disk directly
 *
 * NOTE: The default qcow2 and xml location (created by initial operation) is `/tmp/hotplug_disks/`.
 *       File names are `node1-sda.qcow2` and `node1-sda.xml`.
 *       The target node name is `ndm-vagrant-rancherd_node1`.
 */

const (
	hotplugTargetNodeName  = "ndm-vagrant-rancherd_node1"
	hotplugDiskXMLFileName = "/tmp/hotplug_disks/node1-sda.xml"
	hotplugTargetDiskName  = "sda"
)

type HotPlugTestSuite struct {
	suite.Suite
	SSHClient      *goph.Client
	clientSet      *clientset.Clientset
	targetNodeName string
	targetDiskName string
}

func (s *HotPlugTestSuite) SetupSuite() {
	nodeName := ""
	f, _ := os.Open(filepath.Join(os.Getenv("HOME"), "ssh-config"))
	cfg, _ := ssh_config.Decode(f)
	// consider wildcard, so length shoule be 2
	require.Equal(s.T(), len(cfg.Hosts), 2, "number of Hosts on SSH-config should be 1")
	for _, host := range cfg.Hosts {
		if host.String() == "" {
			// wildcard, continue
			continue
		}
		nodeName = host.Patterns[0].String()
		break
	}
	require.NotEqual(s.T(), nodeName, "", "nodeName should not be empty.")
	s.targetNodeName = nodeName
	targetHost, _ := cfg.Get(nodeName, "HostName")
	targetUser, _ := cfg.Get(nodeName, "User")
	targetPrivateKey, _ := cfg.Get(nodeName, "IdentityFile")
	splitedResult := strings.Split(targetPrivateKey, "node-disk-manager/")
	privateKey := filepath.Join(os.Getenv("HOME"), splitedResult[len(splitedResult)-1])
	// Start new ssh connection with private key.
	auth, err := goph.Key(privateKey, "")
	require.Equal(s.T(), err, nil, "generate ssh auth key should not get error")

	s.SSHClient, err = goph.NewUnknown(targetUser, targetHost, auth)
	require.Equal(s.T(), err, nil, "New ssh connection should not get error")

	kubeconfig := filepath.Join(os.Getenv("HOME"), "kubeconfig")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	require.Equal(s.T(), err, nil, "Generate kubeconfig should not get error")

	s.clientSet, err = clientset.NewForConfig(config)
	require.Equal(s.T(), err, nil, "New clientset should not get error")
}

func (s *HotPlugTestSuite) AfterTest(_, _ string) {
	if s.SSHClient != nil {
		s.SSHClient.Close()
	}
}

func (s *HotPlugTestSuite) BeforeTest(_, _ string) {
	bdi := s.clientSet.HarvesterhciV1beta1().BlockDevices("longhorn-system")
	bdList, err := bdi.List(context.TODO(), v1.ListOptions{})
	require.Equal(s.T(), err, nil, "Get BlockdevicesList should not get error")
	diskCount := 0
	for _, blockdevice := range bdList.Items {
		if blockdevice.Spec.NodeName != s.targetNodeName {
			// focus the target node
			continue
		}
		bdStatus := blockdevice.Status
		if bdStatus.State == "Active" && bdStatus.ProvisionPhase == "Provisioned" {
			diskCount++
			s.targetDiskName = blockdevice.Name
		}
	}
	require.Equal(s.T(), diskCount, 1, "We should only have one disk.")
}

func TestHotPlugDisk(t *testing.T) {
	suite.Run(t, new(HotPlugTestSuite))
}

func (s *HotPlugTestSuite) Test_0_HotPlugRemoveDisk() {
	// remove disk dynamically
	cmd := fmt.Sprintf("virsh detach-disk %s %s --live", hotplugTargetNodeName, hotplugTargetDiskName)
	_, err := s.SSHClient.Run(cmd)
	require.Equal(s.T(), err, nil, "Running command `blkid` should not get error")

	// wait for controller handling
	time.Sleep(1 * time.Second)

	// check disk status
	require.NotEqual(s.T(), s.targetDiskName, "", "target disk name should not be empty before we start hotplug (remove) test")
	bdi := s.clientSet.HarvesterhciV1beta1().BlockDevices("longhorn-system")
	curBlockdevice, err := bdi.Get(context.TODO(), s.targetDiskName, v1.GetOptions{})
	require.Equal(s.T(), err, nil, "Get Blockdevices should not get error")

	require.Equal(s.T(), curBlockdevice.Status.State, diskv1.BlockDeviceInactive, "Disk status should be inactive after we remove disk")

}

func (s *HotPlugTestSuite) Test_1_HotPlugAddDisk() {
	// remove disk dynamically
	cmd := fmt.Sprintf("virsh attach-device --domain %s --file %s --live", hotplugTargetNodeName, hotplugDiskXMLFileName)
	_, err := s.SSHClient.Run(cmd)
	require.Equal(s.T(), err, nil, "Running command `blkid` should not get error")

	// wait for controller handling
	time.Sleep(1 * time.Second)

	// check disk status
	require.NotEqual(s.T(), s.targetDiskName, "", "target disk name should not be empty before we start hotplug (add) test")
	bdi := s.clientSet.HarvesterhciV1beta1().BlockDevices("longhorn-system")
	curBlockdevice, err := bdi.Get(context.TODO(), s.targetDiskName, v1.GetOptions{})
	require.Equal(s.T(), err, nil, "Get Blockdevices should not get error")

	require.Equal(s.T(), curBlockdevice.Status.State, diskv1.BlockDeviceActive, "Disk status should be inactive after we add disk")

}
