package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"vessel/common"
	"vessel/server"
)

func (s *Api) relayInput(_ *server.SyncConn, pool *server.Pool, body json.RawMessage) (common.UpdateType, any, error) {
	var m common.RelayedData
	if err := json.Unmarshal(body, &m); err != nil {
		fmt.Println("errorsdiefef", err)
		return common.UPDATE_RELAY_STDIN, nil, err
	}
	fmt.Println("Received:", m.Id, string(m.Data))
	stdInPipe := pool.GetVesselStdin(m.Id)
	if stdInPipe == nil {
		fmt.Println("errorsdiefef", "efieifh")
		return common.UPDATE_RELAY_STDIN, nil, errors.New("stdin pipe not found")
	}
	defer stdInPipe.Close()
	_, err := stdInPipe.Write(m.Data)
	fmt.Println("fhuehfeuf")
	return common.UPDATE_RELAY_STDIN, true, err
}
