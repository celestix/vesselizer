package api

import (
	"log"

	"vessel/common"
	"vessel/manager"
	"vessel/server"
)

type Api struct {
	log *log.Logger
	vm  *manager.VesselsManager
}

func NewApi(l *log.Logger, vm *manager.VesselsManager) *Api {
	return &Api{
		log: l,
		vm:  vm,
	}
}

func (s *Api) RegisterHandlers(server *server.Server) {
	server.RegisterHandler(common.UPDATE_ECHO, s.echoHandler)
	server.RegisterHandler(common.UPDATE_NEW_VESSEL, s.newVesselHandler)
	server.RegisterHandler(common.UPDATE_STOP_VESSEL, s.stopVesselHandler)
	server.RegisterHandler(common.UPDATE_START_VESSEL, s.startVesselHandler)
	server.RegisterHandler(common.UPDATE_RELAY_STDIN, s.relayInput)
}
