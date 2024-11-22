package client

import (
	"encoding/json"
	"vessel/common"
)

func invoke[T any](c *Client, method common.UpdateType, message any) (*T, error) {
	resp, err := c.invoke(method, message)
	if err != nil {
		return nil, err
	}
	var d T
	return &d, json.Unmarshal(resp, &d)
}

func (c *Client) Echo() (string, error) {
	m, err := invoke[common.EchoUpdate](c, common.UPDATE_ECHO, "")
	if err != nil {
		return "", err
	}
	return m.Message, err
}

func (c *Client) NewVessel(name, appDir, baseImage, buildFile string, entrypoint ...string) (string, error) {
	m, err := invoke[common.NewVesselResponse](c, common.UPDATE_NEW_VESSEL, &common.NewVesselRequest{
		Name:         name,
		BaseImage:    baseImage,
		Entrypoint:   entrypoint,
		BuildFile:    buildFile,
		AppDirectory: appDir,
	})
	if err != nil {
		return "", err
	}
	return m.Id, err
}

func (c *Client) StopVessel(id string) (bool, error) {
	m, err := invoke[common.VesselControlResponse](c, common.UPDATE_STOP_VESSEL, &common.VesselControlRequest{Id: id})
	return m != nil && m.Success, err
}

func (c *Client) StartVessel(id string) (bool, error) {
	m, err := invoke[common.VesselControlResponse](c, common.UPDATE_START_VESSEL, &common.VesselControlRequest{Id: id})
	return m != nil && m.Success, err
}

func (c *Client) RelayStdin(id string, data []byte) error {
	_, err := c.invoke(common.UPDATE_RELAY_STDIN, &common.RelayedData{Id: id, Data: data})
	return err
}
