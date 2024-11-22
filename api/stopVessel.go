package api

import (
	"encoding/json"
	"vessel/common"
	"vessel/server"
)

func (s *Api) stopVesselHandler(_ *server.SyncConn, pool *server.Pool, body json.RawMessage) (common.UpdateType, any, error) {
	var m common.VesselControlRequest
	if err := json.Unmarshal(body, &m); err != nil {
		return common.UPDATE_STOP_VESSEL, nil, err
	}
	pool.RemoveConnection(m.Id)
	err := s.vm.StopVessel(m.Id)
	return common.UPDATE_STOP_VESSEL, &common.VesselControlResponse{Success: err == nil}, err
}
