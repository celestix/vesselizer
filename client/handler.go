package client

import (
	"encoding/json"
	"vessel/common"
)

type Handler interface {
	Handle(json.RawMessage) error
}

type CallbackFunction func(string, []byte) error

type RelayedOutputHandler struct {
	callback CallbackFunction
}

func NewRelayedOutputHandler(cb CallbackFunction) *RelayedOutputHandler {
	return &RelayedOutputHandler{
		callback: cb,
	}
}

func (r *RelayedOutputHandler) Handle(msg json.RawMessage) error {
	var data common.RelayedData
	err := json.Unmarshal(msg, &data)
	if err != nil {
		return err
	}
	return r.callback(data.Id, data.Data)
}
