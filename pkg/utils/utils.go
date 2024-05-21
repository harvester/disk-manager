package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/longhorn/longhorn-manager/util"
)

const (
	// ProcPath is a vfs storing process info for Linux.
	ProcPath = "/proc"
	// HostProcPath is the convention path where host `/proc` is mounted.
	HostProcPath = "/host/proc"
	// DiskRemoveTag indicates a Longhorn is pending to remove.
	DiskRemoveTag = "harvester-ndm-disk-remove"
)

var CmdTimeoutError error

var ext4MountOptions = strings.Join([]string{
	"journal_checksum",
	"journal_ioprio=0",
	"barrier=1",
	"errors=remount-ro",
}, ",")

// IsHostProcMounted checks if host's proc info `/proc` is mounted on `/host/proc`
func IsHostProcMounted() (bool, error) {
	_, err := os.Stat(HostProcPath)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	return err == nil, nil
}

func GetFullDevPath(shortPath string) string {
	if shortPath == "" {
		return ""
	}
	return fmt.Sprintf("/dev/%s", shortPath)
}

func MatchesIgnoredCase(s []string, k string) bool {
	for _, e := range s {
		if strings.EqualFold(e, k) {
			return true
		}
	}
	return false
}

func ContainsIgnoredCase(s []string, k string) bool {
	k = strings.ToLower(k)
	for _, e := range s {
		if strings.Contains(k, strings.ToLower(e)) {
			return true
		}
	}
	return false
}

func MakeExt4DiskFormatting(devPath, uuid string) error {
	args := []string{"-F", devPath}
	if uuid != "" {
		args = append(args, "-U", uuid)
	}
	cmd := exec.Command("mkfs.ext4", args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to format %s. err: %v", devPath, err)
	}
	return nil
}

// MountDisk mounts the specified ext4 volume device to the specified path
func MountDisk(devPath, mountPoint string) error {
	var needMkdir bool
	if _, err := os.Stat(mountPoint); err != nil && !os.IsNotExist(err) {
		return err
	} else if os.IsNotExist(err) {
		needMkdir = true
	}

	isHostProcMounted, err := IsHostProcMounted()
	if err != nil {
		return err
	}

	if needMkdir {
		if isHostProcMounted {
			if _, err := executeOnHostNamespace("mkdir", []string{"-p", mountPoint}); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(mountPoint, os.ModeDir); err != nil {
				return err
			}
		}
	}

	if isHostProcMounted {
		return mountExt4OnHostNamespace(devPath, mountPoint, false)
	}

	return mountExt4(devPath, mountPoint, false)
}

// UmountDisk unmounts the specified volume device to the specified path
func UmountDisk(path string) error {
	isHostProcMounted, err := IsHostProcMounted()
	if err != nil {
		return err
	}
	if isHostProcMounted {
		_, err := executeOnHostNamespace("umount", []string{path})
		return err
	}
	err = syscall.Unmount(path, 0)
	return os.NewSyscallError("umount", err)
}

func ForceUmountWithTimeout(path string, timeout time.Duration) error {
	isHostProcMounted, err := IsHostProcMounted()
	if err != nil {
		return err
	}
	if isHostProcMounted {
		_, err := executeOnHostNamespaceWithTimeout("umount", []string{"-f", path}, timeout)
		return err
	}
	// flags, MNT_FORCE -> 1
	err = syscall.Unmount(path, 1)
	return os.NewSyscallError("umount", err)
}

func mountExt4(device, path string, readonly bool) error {
	var flags uintptr
	flags = syscall.MS_RELATIME
	if readonly {
		flags |= syscall.MS_RDONLY
	}
	err := syscall.Mount(device, path, "ext4", flags, ext4MountOptions)
	return os.NewSyscallError("mount", err)
}

// mountExt4OnHostNamespace provides the same functionality as mountExt4 but on host namespace.
func mountExt4OnHostNamespace(device, path string, readonly bool) error {
	ns := GetHostNamespacePath(util.HostProcPath)
	executor, err := NewExecutorWithNS(ns)
	if err != nil {
		return err
	}

	opts := ext4MountOptions + ",relatime"
	if readonly {
		opts = opts + ",ro"
	}

	_, err = executor.Execute("mount", []string{"-t", "ext4", "-o", opts, device, path})
	return err
}

func executeOnHostNamespace(cmd string, args []string) (string, error) {
	ns := GetHostNamespacePath(util.HostProcPath)
	executor, err := NewExecutorWithNS(ns)
	if err != nil {
		return "", err
	}
	return executor.Execute(cmd, args)
}

func executeOnHostNamespaceWithTimeout(cmd string, args []string, cmdTimeout time.Duration) (string, error) {
	ns := GetHostNamespacePath(util.HostProcPath)
	executor, err := NewExecutorWithNS(ns)
	executor.SetTimeout(cmdTimeout)
	if err != nil {
		return "", err
	}
	return executor.Execute(cmd, args)
}

func IsFSCorrupted(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "wrong fs type")
}

func IsSupportedFileSystem(fsType string) bool {
	if fsType == "ext4" || fsType == "xfs" {
		return true
	}
	return false
}

// CallerWithLock is a helper function to call a function with a condition lock
func CallerWithCondLock[T any](cond *sync.Cond, f func() T) T {
	cond.L.Lock()
	defer cond.L.Unlock()
	return f()
}
