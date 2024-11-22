package api

import (
	"encoding/json"
	"vessel/common"
	"vessel/server"
)

func (s *Api) echoHandler(conn *server.SyncConn, pool *server.Pool, body json.RawMessage) (common.UpdateType, any, error) {
	return common.UPDATE_ECHO, &common.EchoUpdate{Message: "server is up and running"}, nil
}
