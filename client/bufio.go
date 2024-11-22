package client

import (
	"net"
	"vessel/common"
)

func read(conn net.Conn) ([]byte, error) {
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

func write(conn net.Conn, b []byte) error {
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
