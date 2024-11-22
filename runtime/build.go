package runtime

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"vessel/config"
)

func debugPrintln(debug bool, a ...any) {
	if !debug {
		return
	}
	log.Println(a...)
}

func StartVessel(id string, build bool, debug bool) error {
	vesselDir := filepath.Join(config.DATA_DIR, config.VESSELS_DIR, id)
	vesselConfigFile, err := os.Open(filepath.Join(vesselDir, "config.json"))
	if err != nil {
		return err
	}
	defer vesselConfigFile.Close()
	var vesselConfig config.VesselConfig
	err = json.NewDecoder(vesselConfigFile).Decode(&vesselConfig)
	if err != nil {
		return err
	}
	debugPrintln(debug, "Starting Vessel:", vesselConfig.Name)
	rootFSDir := filepath.Join(vesselDir, "fs")
	debugPrintln(debug, "Binding Root Dir:", rootFSDir)
	if err = MountBind(rootFSDir); err != nil {
		return err
	}
	debugPrintln(debug, "Preparing for pivot_root")
	debugPrintln(debug, "Creating a directory for old root")
	oldRootDir := filepath.Join(rootFSDir, "oldroot")
	_ = os.Mkdir(oldRootDir, 0700)
	debugPrintln(debug, "Performing pivot_root between", rootFSDir, "and", oldRootDir)
	err = syscall.PivotRoot(rootFSDir, oldRootDir)
	if err != nil {
		return err
	}
	debugPrintln(debug, "Change directory to new root")
	err = syscall.Chdir("/")
	if err != nil {
		return err
	}
	debugPrintln(debug, "Unmounting old root")
	err = syscall.Unmount("/oldroot", syscall.MNT_DETACH)
	if err != nil {
		log.Println("Warning: failed to unmount old root directory:", err)
	}
	err = os.Remove("/oldroot")
	if err != nil {
		log.Println("Warning: failed to remove old root directory:", err)
	}
	debugPrintln(debug, "Writing resolv.conf")
	err = WriteResolvConf()
	if err != nil {
		log.Println("Warning: failed to write resolv.conf:", err)
	}
	debugPrintln(debug, "Mounting sys files")
	err = PerformSysMounts()
	if err != nil {
		log.Println("Warning: failed to mount sys files:", err)
	}
	var env []string = []string{"PATH=/bin:/usr/bin:/sbin:/usr/sbin:/usr/local/sbin:/usr/local/bin", "HOME=/root"}
	if build {
		debugPrintln(debug, "Running buildfile of vessel")
		cmd := exec.Command("/bin/sh", fmt.Sprintf("/app/%s", vesselConfig.BuildFile))
		cmd.Env = env
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
		debugPrintln(debug, "Build completed successfully")
	}
	debugPrintln(debug, "Running entrypoint of vessel")
	cmd := exec.Command(vesselConfig.Entrypoint[0], vesselConfig.Entrypoint[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func WriteResolvConf() error {
	resolvConf := `nameserver 1.1.1.1`
	return os.WriteFile("/etc/resolv.conf", []byte(resolvConf), 0644)
}
