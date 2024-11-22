package server

import (
	"io"
	"sync"
	"vessel/common"
)

func read(mu *sync.Mutex, conn io.Reader) ([]byte, error) {
	mu.Lock()
	defer mu.Unlock()
	head := make([]byte, 4)
	_, err := conn.Read(head)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, common.BytesToInt(head))
	_, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func write(mu *sync.Mutex, conn io.Writer, b []byte) error {
	mu.Lock()
	defer mu.Unlock()
	_, err := conn.Write(common.IntToBytes(uint32(len(b))))
	if err != nil {
		return err
	}
	_, err = conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}
