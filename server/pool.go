package server

import (
	"io"
	"log"
	"sync"
	"vessel/common"
)

type Pool struct {
	l  *log.Logger
	mu *sync.RWMutex
	m  map[string]*SyncConn
	e  map[string]*Error
	// Vessel Stdin
	rmu  *sync.RWMutex
	rMap map[string]*io.WriteCloser
}

func NewPool(l *log.Logger) *Pool {
	return &Pool{
		l:    l,
		mu:   &sync.RWMutex{},
		m:    make(map[string]*SyncConn),
		e:    make(map[string]*Error),
		rmu:  &sync.RWMutex{},
		rMap: make(map[string]*io.WriteCloser),
	}
}

func (p *Pool) HasConnection(uid string) bool {
	p.mu.RLock()
	_, ok := p.m[uid]
	p.mu.RUnlock()
	return ok
}

func (p *Pool) RemoveConnection(uid string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	conn, ok := p.m[uid]
	if !ok {
		return
	}
	defer conn.Conn.Close()
	delete(p.m, uid)
}

func (p *Pool) SetConnection(uid string, sconn *SyncConn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.m[uid] = sconn
}

func (p *Pool) Broadcast(uid string, data []byte) {
	head := common.IntToBytes(uint32(len(data)))
	p.mu.RLock()
	sconn := p.m[uid]
	p.mu.RUnlock()
	// connection not found
	if sconn == nil {
		return
	}
	sconn.wmu.Lock()
	defer sconn.wmu.Unlock()
	_, err := sconn.Conn.Write(head)
	if err != nil {
		p.RemoveConnection(uid)
		return
	}
	_, err = sconn.Conn.Write(data)
	if err != nil {
		p.RemoveConnection(uid)
		return
	}
}

func (p *Pool) SetVesselStdin(uid string, inPipe *io.WriteCloser) {
	p.rmu.Lock()
	defer p.rmu.Unlock()
	p.rMap[uid] = inPipe
}

func (p *Pool) GetVesselStdin(uid string) io.WriteCloser {
	p.rmu.RLock()
	defer p.rmu.RUnlock()
	wr := p.rMap[uid]
	if wr == nil {
		return nil
	}
	return *wr
}

func (p *Pool) WriteError(uid string, errType ErrorType, errMessage string) {
	p.mu.RLock()
	err, ok := p.e[uid]
	if ok && err.Type == ErrorTypeCritical && errType != ErrorTypeCritical {
		p.mu.RUnlock()
		return
	}
	p.mu.RUnlock()
	p.mu.Lock()
	defer p.mu.Unlock()
	p.e[uid] = &Error{errType, errMessage}
}

func (p *Pool) ForceWriteError(uid string, errType ErrorType, errMessage string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.e[uid] = &Error{errType, errMessage}
}

func (p *Pool) GetError(uid string) *Error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.e[uid]
}
