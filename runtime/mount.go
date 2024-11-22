package runtime

import (
	"io/fs"
	"log"
	"os"
	"syscall"

	"github.com/pkg/errors"
)

// Mount bind (mount --bind x x)
func MountBind(target string) error {
	return syscall.Mount(target, target, "", syscall.MS_BIND, "")
}

// Lazy unmount (umount -l x)
func Unmount(target string) error {
	return syscall.Unmount(target, syscall.MNT_DETACH)
}

func createDevice(path string, perm fs.FileMode, mode uint32, major, minor int) error {
	// Calculate the device number as (major << 8) | minor
	dev := (major << 8) | minor
	err := syscall.Mknod(path, mode, dev)
	if err != nil {
		return errors.Wrap(err, "failed to create device")
	}
	return os.Chmod(path, perm)
}

func PerformSysMounts() error {
	log.Println("Mounting proc")
	err := syscall.Mount("proc", "proc", "proc", 0, "")
	if err != nil {
		return errors.Wrap(err, "failed to mount proc")
	}
	// mount -t tmpfs -o rw,nosuid,size=65536k,mode=755,inode64 tmpfs /dev
	log.Println("Mounting tmpfs")
	err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID, "size=65536k,mode=755,inode64")
	if err != nil {
		return errors.Wrap(err, "failed to mount tmpfs")
	}
	// mount -t devpts -o rw,nosuid,noexec,relatime,gid=5,mode=620,ptmxmode=666 devpts /dev/pts
	log.Println("Mounting devpts")
	err = os.Mkdir("/dev/pts", 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/pts")
	}
	err = syscall.Mount("devpts", "/dev/pts", "devpts", syscall.MS_RELATIME|syscall.MS_NOSUID|syscall.MS_NOEXEC, "gid=5,mode=620,ptmxmode=666")
	if err != nil {
		return errors.Wrap(err, "failed to mount devpts")
	}
	// 	mount -t sysfs -o ro,nosuid,nodev,noexec,relatime sysfs /sys
	log.Println("Mounting sysfs")
	err = syscall.Mount("sysfs", "/sys", "sysfs", syscall.MS_RDONLY|syscall.MS_NOSUID|syscall.MS_NODEV|syscall.MS_NOEXEC, "")
	if err != nil {
		return errors.Wrap(err, "failed to mount sysfs")
	}
	// mount -t tmpfs -o rw,nosuid,nodev,noexec,relatime,size=65536k,inode64 tmpfs /dev/shm
	log.Println("Mounting tmpfs at /dev/shm")
	err = os.Mkdir("/dev/shm", 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/shm")
	}
	err = syscall.Mount("tmpfs", "/dev/shm", "tmpfs", syscall.MS_NOSUID|syscall.MS_NODEV|syscall.MS_NOEXEC, "size=65536k,inode64")
	if err != nil {
		return errors.Wrap(err, "failed to mount tmpfs at /dev/shm")
	}
	// mount -t devpts -o rw,nosuid,noexec,relatime,gid=5,mode=620,ptmxmode=666 devpts /dev/console
	log.Println("Mounting devpts at /dev/console")
	err = os.Mkdir("/dev/console", 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/console")
	}
	err = syscall.Mount("devpts", "/dev/console", "devpts", syscall.MS_NOSUID|syscall.MS_NOEXEC, "gid=5,mode=620,ptmxmode=666")
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/console")
	}
	// mount -t mqueue -o rw,nosuid,nodev,noexec,relatime mqueue /dev/mqueue
	log.Println("Mounting mqueue at /dev/mqueue")
	err = os.Mkdir("/dev/mqueue", 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/mqueue")
	}
	err = syscall.Mount("mqueue", "/dev/mqueue", "mqueue", syscall.MS_NOSUID|syscall.MS_NODEV|syscall.MS_NOEXEC|syscall.MS_RELATIME, "")
	if err != nil {
		return errors.Wrap(err, "failed to mount mqueue at /dev/mqueue")
	}
	// mknod -m 666 /dev/null c 1 3
	log.Println("Creating /dev/null")
	err = createDevice("/dev/null", 0666, syscall.S_IFCHR, 1, 3)
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/null")
	}
	// ln -s /proc/self/fd /dev/fd
	log.Println("Creating /dev/fd")
	err = os.Symlink("/proc/self/fd", "/dev/fd")
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/fd")
	}
	// mknod -m 666 /dev/random c 1 8
	log.Println("Creating /dev/random")
	err = createDevice("/dev/random", 0666, syscall.S_IFCHR, 1, 8)
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/random")
	}
	// mknod -m 666 /dev/urandom c 1 9
	log.Println("Creating /dev/urandom")
	err = createDevice("/dev/urandom", 0666, syscall.S_IFCHR, 1, 9)
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/urandom")
	}
	// mknod -m 666 /dev/zero c 1 5
	log.Println("Creating /dev/zero")
	err = createDevice("/dev/zero", 0666, syscall.S_IFCHR, 1, 5)
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/zero")
	}
	// mknod -m 666 /dev/full c 1 7
	log.Println("Creating /dev/full")
	err = createDevice("/dev/full", 0666, syscall.S_IFCHR, 1, 7)
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/full")
	}
	// mknod -m 666 /dev/tty c 5 0
	log.Println("Creating /dev/tty")
	err = createDevice("/dev/tty", 0666, syscall.S_IFCHR, 5, 0)
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/tty")
	}
	// ln -s /dev/pts/ptmx /dev/ptmx
	log.Println("Creating /dev/ptmx")
	err = os.Symlink("/dev/pts/ptmx", "/dev/ptmx")
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/ptmx")
	}
	// ln -f /proc/self/fd/0 /dev/stdin
	log.Println("Creating /dev/stdin")
	err = os.Symlink("/proc/self/fd/0", "/dev/stdin")
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/stdin")
	}
	// ln -f /proc/self/fd/1 /dev/stdout
	log.Println("Creating /dev/stdout")
	err = os.Symlink("/proc/self/fd/1", "/dev/stdout")
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/stdout")
	}
	// ln -f /proc/self/fd/2 /dev/stderr
	log.Println("Creating /dev/stderr")
	err = os.Symlink("/proc/self/fd/2", "/dev/stderr")
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/stderr")
	}
	// ln -f /proc/kcore /dev/core
	log.Println("Creating /dev/core")
	err = os.Symlink("/proc/kcore", "/dev/core")
	if err != nil {
		return errors.Wrap(err, "failed to create /dev/core")
	}
	return nil
}
