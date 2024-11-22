package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type VesselConfig struct {
	Name       string
	Id         string
	BaseImage  string
	BuildFile  string
	Entrypoint []string
}

func (v *VesselConfig) JSON() ([]byte, error) {
	return json.Marshal(v)
}

func (v *VesselConfig) Save(vesselDir string) error {
	file, err := os.Create(filepath.Join(vesselDir, "config.json"))
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := v.JSON()
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}
