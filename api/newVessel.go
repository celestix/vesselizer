package api

import (
	"encoding/json"
	"io"
	"log"
	"vessel/common"
	"vessel/server"
)

func (s *Api) newVesselHandler(conn *server.SyncConn, pool *server.Pool, body json.RawMessage) (common.UpdateType, any, error) {
	var m common.NewVesselRequest
	if err := json.Unmarshal(body, &m); err != nil {
		return common.UPDATE_NEW_VESSEL, nil, err
	}
	vessel, err := s.vm.NewVessel(m.Name, m.AppDirectory, m.BaseImage, m.BuildFile, m.Entrypoint...)
	if err != nil {
		return common.UPDATE_NEW_VESSEL, nil, err
	}
	pool.SetConnection(vessel.Id, conn)
	var stdInPipe io.WriteCloser
	var stdOutPipe, stdErrPipe io.ReadCloser
	pool.SetVesselStdin(vessel.Id, &stdInPipe)
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
					log.Println("vessel["+vessel.Id+"] stdout:", "pipe has been closed")
					return
				}
				log.Println("vessel[" + vessel.Id + "] stdout:")
				return
			}
			pool.Broadcast(vessel.Id, server.MakeResult(common.UPDATE_RELAY_STDOUT, &common.RelayedData{Id: vessel.Id, Data: buf[:n]}))
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
					log.Println("vessel["+vessel.Id+"] stderr:", "pipe has been closed")
					return
				}
				log.Println("vessel[" + vessel.Id + "] stderr:")
				return
			}
			pool.Broadcast(vessel.Id, server.MakeResult(common.UPDATE_RELAY_STDERR, &common.RelayedData{Id: vessel.Id, Data: buf[:n]}))
		}
	}()
	// start vessel in a separate goroutine so that it doesn't block main goroutine
	go func() {
		err := vessel.Start(true, &stdInPipe, &stdOutPipe, &stdErrPipe)
		if err != nil {
			log.Println("vessel[" + vessel.Id + "] start:")
			return
		}
	}()
	return common.UPDATE_NEW_VESSEL, &common.NewVesselResponse{Id: vessel.Id}, nil
}
