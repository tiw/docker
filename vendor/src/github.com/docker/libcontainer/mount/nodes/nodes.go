// +build linux

package nodes

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/docker/libcontainer/devices"
	"github.com/dotcloud/docker/pkg/system"
)

// Create the device nodes in the container.
func CreateDeviceNodes(rootfs string, nodesToCreate []*devices.Device) error {
	oldMask := system.Umask(0000)
	defer system.Umask(oldMask)

	for _, node := range nodesToCreate {
		if err := CreateDeviceNode(rootfs, node); err != nil {
			return err
		}
	}
	return nil
}

// Creates the device node in the rootfs of the container.
func CreateDeviceNode(rootfs string, node *devices.Device) error {
	var (
		dest   = filepath.Join(rootfs, node.Path)
		parent = filepath.Dir(dest)
	)

	if err := os.MkdirAll(parent, 0755); err != nil {
		return err
	}

	fileMode := node.FileMode
	switch node.Type {
	case 'c':
		fileMode |= syscall.S_IFCHR
	case 'b':
		fileMode |= syscall.S_IFBLK
	default:
		return fmt.Errorf("%c is not a valid device type for device %s", node.Type, node.Path)
	}

	if err := system.Mknod(dest, uint32(fileMode), devices.Mkdev(node.MajorNumber, node.MinorNumber)); err != nil && !os.IsExist(err) {
		return fmt.Errorf("mknod %s %s", node.Path, err)
	}
	return nil
}
