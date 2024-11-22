package runtime

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"vessel/config"
	"vessel/errors"
)

type VesselStatus int

const (
	VesselRunning VesselStatus = iota
	VesselStopped
	VesselBuilt
)

type Vessel struct {
	// a unique id of vessel
	Id string
	// name of vessel
	Name string
	// Base image used by the vessel
	BaseImage string
	// Path of root in host machine
	RootPath string
	// Build file that is run just once
	BuildFile string
	// Entrypoint of the vessel
	Entrypoint []string
	// Vessel status
	Status VesselStatus

	cmd *exec.Cmd
}

func Create(name string, appDir string, baseImage string, buildFile string, entryPoint ...string) (*Vessel, error) {
	if _, ok := config.Images[baseImage]; !ok {
		return nil, errors.ErrSpecifiedBaseImageNotFound
	}
	id := generateId()
	v := Vessel{
		Id:         id,
		Name:       name,
		BaseImage:  baseImage,
		RootPath:   filepath.Join(config.DATA_DIR, config.VESSELS_DIR, id),
		BuildFile:  buildFile,
		Entrypoint: entryPoint,
	}
	vesselDir := filepath.Join(config.DATA_DIR, config.VESSELS_DIR, v.Id)
	err := os.MkdirAll(vesselDir, 0755)
	if err != nil {
		return nil, err
	}
	// /usr/lib/vesselizer/vessels/<id>/config.json
	err = v.GenerateConfig(vesselDir)
	if err != nil {
		return nil, err
	}
	// /usr/lib/vesselizer/vessels/<id>/fs
	err = v.GenerateRoot(vesselDir)
	if err != nil {
		return nil, err
	}
	err = v.CopyHostDirectory(vesselDir, appDir)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (v *Vessel) GenerateRoot(vesselDir string) error {
	// Create the root directory
	rootDir := filepath.Join(vesselDir, "fs")
	err := os.Mkdir(rootDir, 0755)
	if err != nil {
		return err
	}
	// Generate root
	image := config.Images[v.BaseImage]
	imageDir := filepath.Join(config.DATA_DIR, config.IMAGES_DIR, image)
	cmd := exec.Command("tar", "-xf", imageDir, "-C", rootDir)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (v *Vessel) GenerateConfig(vesselDir string) error {
	cv := config.VesselConfig{
		Id:         v.Id,
		Name:       v.Name,
		BaseImage:  v.BaseImage,
		BuildFile:  v.BuildFile,
		Entrypoint: v.Entrypoint,
	}
	return cv.Save(vesselDir)
}

func (v *Vessel) CopyHostDirectory(vesselDir string, src string) error {
	tgtAppDir := filepath.Join(vesselDir, "fs", "app")
	return copyDir(src, tgtAppDir)
}

func (v *Vessel) Start(build bool, inPipe *io.WriteCloser, outPipe, errPipe *io.ReadCloser) error {
	cmd := exec.Command("/proc/self/exe", "start_pre_iso_vessel", v.Id, fmt.Sprint(build))
	stdInPipe, _ := cmd.StdinPipe()
	stdOutPipe, _ := cmd.StdoutPipe()
	stdErrPipe, _ := cmd.StderrPipe()
	*inPipe = stdInPipe
	*outPipe = stdOutPipe
	*errPipe = stdErrPipe
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	log.Println("Running Command")
	v.cmd = cmd
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Println("Ran Command")
	log.Println(cmd)
	return nil
}

func (v *Vessel) Stop() error {
	if v.cmd == nil {
		return errors.ErrVesselNotRunning
	}
	err := v.cmd.Process.Signal(syscall.SIGKILL)
	if err != nil {
		return err
	}
	v.cmd = nil
	return nil
}

func generateId() string {
	var b = make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}
