package server

import (
	"encoding/json"

	"vessel/common"
)

type HandlerFunc func(
	conn *SyncConn,
	pool *Pool,
	body json.RawMessage,
) (
	common.UpdateType,
	any,
	error,
)
