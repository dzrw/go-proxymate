package proxymate

import (
	"log"
	"net"
	"sync"
)

func ListenTCP(addr string) (wl WaitListener, err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	wl = &waitListener{
		Listener: l,
		wg:       &sync.WaitGroup{},
	}
	return
}

// ----------------------------------------------

type WaitListener interface {
	net.Listener
	Wait()
}

// ----------------------------------------------

type waitListener struct {
	net.Listener
	wg *sync.WaitGroup
}

func (l *waitListener) Wait() {
	l.wg.Wait()
}

// Add to the wait group on success and returns a waitConn.
func (l *waitListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return c, err
	}

	log.Printf("connection accepted: %s", c.RemoteAddr().String())

	l.wg.Add(1)
	return &waitConn{c, l.wg}, nil
}

// ----------------------------------------------

type waitConn struct {
	net.Conn
	wg *sync.WaitGroup
}

// Calls Done() on the wait group.
func (c *waitConn) Close() error {
	err := c.Conn.Close()
	c.wg.Done()
	log.Printf("connection closed: %s", c.RemoteAddr().String())
	return err
}
