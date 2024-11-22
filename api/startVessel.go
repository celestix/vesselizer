package api

import (
	"encoding/json"
	"io"
	"log"
	"vessel/common"
	"vessel/server"
)

func (s *Api) startVesselHandler(sconn *server.SyncConn, pool *server.Pool, body json.RawMessage) (common.UpdateType, any, error) {
	var m common.VesselControlRequest
	if err := json.Unmarshal(body, &m); err != nil {
		return common.UPDATE_START_VESSEL, nil, err
	}
	pool.SetConnection(m.Id, sconn)
	var stdInPipe io.WriteCloser
	var stdOutPipe, stdErrPipe io.ReadCloser
	pool.SetVesselStdin(m.Id, &stdInPipe)
	// relay stdout
	go func() {
		for stdOutPipe == nil {
			continue
		}
		log.Println("stdOutPipe is available")
		for {
			buf := make([]byte, 1024)
			n, err := stdOutPipe.Read(buf)
			if err != nil {
				if err == io.EOF {
					log.Println("vessel["+m.Id+"] stdout:", "pipe has been closed")
					return
				}
				log.Println("vessel[" + m.Id + "] stdout:")
				return
			}
			pool.Broadcast(m.Id, server.MakeResult(common.UPDATE_RELAY_STDOUT, &common.RelayedData{Id: m.Id, Data: buf[:n]}))
		}
	}()
	// relay stderr
	go func() {
		for stdErrPipe == nil {
			continue
		}
		log.Println("stdErrPipe is available")
		for {
			buf := make([]byte, 1024)
			n, err := stdErrPipe.Read(buf)
			if err != nil {
				if err == io.EOF {
					log.Println("vessel["+m.Id+"] stderr:", "pipe has been closed")
					return
				}
				log.Println("vessel[" + m.Id + "] stderr:")
				return
			}
			pool.Broadcast(m.Id, server.MakeResult(common.UPDATE_RELAY_STDERR, &common.RelayedData{Id: m.Id, Data: buf[:n]}))
		}
	}()
	// start vessel in a separate goroutine so that it doesn't block main goroutine
	go func() {
		log.Println("vessel["+m.Id+"] start:", s.vm.StartVessel(m.Id, &stdInPipe, &stdOutPipe, &stdErrPipe))
	}()
	return common.UPDATE_START_VESSEL, &common.VesselControlResponse{Success: true}, nil
}
