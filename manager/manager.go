package manager

import (
	"fmt"
	"io"
	"sync"
	"vessel/errors"
	"vessel/runtime"
)

type VesselsManager struct {
	rmu     *sync.RWMutex
	vessels map[string]*runtime.Vessel
}

func NewVesselsManager() *VesselsManager {
	return &VesselsManager{
		rmu:     new(sync.RWMutex),
		vessels: make(map[string]*runtime.Vessel),
	}
}

func (m *VesselsManager) NewVessel(name, appDir, baseImage, buildFile string, entryPoint ...string) (*runtime.Vessel, error) {
	vessel, err := runtime.Create(name, appDir, baseImage, buildFile, entryPoint...)
	if err != nil {
		return nil, err
	}
	m.rmu.Lock()
	defer m.rmu.Unlock()
	m.vessels[vessel.Id] = vessel
	fmt.Println(m.vessels)
	return vessel, nil
}

func (m *VesselsManager) StartVessel(id string, inPipe *io.WriteCloser, outPipe, errPipe *io.ReadCloser) error {
	m.rmu.RLock()
	vessel, ok := m.vessels[id]
	m.rmu.RUnlock()
	if !ok {
		return errors.ErrVesselNotFound
	}
	return vessel.Start(false, inPipe, outPipe, errPipe)
}

func (m *VesselsManager) StopVessel(id string) error {
	m.rmu.RLock()
	vessel, ok := m.vessels[id]
	m.rmu.RUnlock()
	if !ok {
		return errors.ErrVesselNotFound
	}
	return vessel.Stop()
}
