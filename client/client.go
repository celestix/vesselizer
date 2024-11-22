package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"vessel/common"
)

type Client struct {
	mu     *sync.RWMutex
	d      *Dispatcher
	conn   net.Conn
	listen bool
}

func NewClient() (*Client, error) {
	socketPath := filepath.Join(os.TempDir(), "vesselizer.sock")
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		err = fmt.Errorf("error connecting to server: %s", err.Error())
		return nil, err
	}
	return &Client{
		conn: conn,
		mu:   &sync.RWMutex{},
		d: &Dispatcher{
			Handlers: make(map[common.UpdateType][]Handler),
		},
	}, nil
}

func (c *Client) Listen() (err error) {
	defer c.conn.Close()
	c.listen = true
	for c.listen {
		c.mu.RLock()
		var buf []byte
		buf, err = read(c.conn)
		if err != nil {
			c.mu.RUnlock()
			if err == io.EOF {
				err = nil
				return
			}
			err = fmt.Errorf("error reading: %s", err.Error())
			return
		}
		err = c.d.process(buf)
		if err != nil {
			c.mu.RUnlock()
			if err == ErrDisconnect {
				break
			}
			err = fmt.Errorf("error processing: %s", err.Error())
			return
		}
		c.mu.RUnlock()
	}
	return
}

func (c *Client) AddHandler(t common.UpdateType, h Handler) {
	c.d.AddHandler(t, h)
}

func (c *Client) RemoveHandler(t common.UpdateType) {
	c.d.RemoveHandler(t)
}

func (c *Client) Disconnect() {
	c.listen = false
}

func (c *Client) invoke(method common.UpdateType, message any) (json.RawMessage, error) {
	// block updates listener while invoking a method
	// to retrieve the message update here instead
	c.mu.Lock()
	defer c.mu.Unlock()
	buf, err := json.Marshal(&Request{
		Method:  method,
		Message: message,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to invoke %s: %s", common.UpdateTypeToString[method], err.Error())
	}
	err = write(c.conn, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke %s: %s", common.UpdateTypeToString[method], err.Error())
	}
	buf, err = read(c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke %s: %s", common.UpdateTypeToString[method], err.Error())
	}
	var res Response
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %s", common.UpdateTypeToString[method], err.Error())
	}
	if !res.Ok {
		return nil, errors.New(res.Error)
	}
	return res.Update.Message, nil
}
