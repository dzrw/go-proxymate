package proxymate

import (
	"log"
	"net"
	"sync"
)

type Handler interface {
	Serve(conn net.Conn, die chan int)
}

type Server struct {
	Addr    string
	Handler Handler

	wl WaitListener

	mu       *sync.Mutex
	stopping bool
}

func ListenAndServe(addr string, handler Handler) (srv *Server, err error) {
	l, err := ListenTCP(addr)
	if err != nil {
		return
	}

	log.Printf("Listening on %s\n", addr)

	srv = &Server{Addr: addr, Handler: handler, wl: l, mu: &sync.Mutex{}}
	srv.Serve()

	return
}

func (srv *Server) Stop() (err error) {
	log.Println("stopping server")

	srv.mu.Lock()
	srv.stopping = true
	srv.mu.Unlock()

	err = srv.wl.Close()
	if err != nil {
		return
	}

	log.Println("refusing new connections to " + srv.Addr)
	log.Println("waiting for existing requests to drain...")

	// TODO: But what if they don't drain? How can we force them to die when
	// we're stuck in a wait?
	srv.wl.Wait()

	log.Printf("remaining connections have drained")

	return
}

func (srv *Server) IsStopping() bool {
	srv.mu.Lock()
	b := srv.stopping
	srv.mu.Unlock()
	return b
}

func (srv *Server) Serve() {
	die := make(chan int)

	go func() {
		for {
			// Wait for a connection.
			conn, err := srv.wl.Accept()
			if err != nil {
				if srv.IsStopping() {
					// This only works because we expect Handlers to set up some
					// kind of for-select loop and test whether this channel is closed.
					// If we couldn't assume that, then I think we'd pass a <Mutex, bool>
					// tuple to each connection which the handler would check in its
					// processing loop.
					close(die)
					break
				}

				log.Fatal(err)
			}

			// Handle the connection in a new goroutine.
			// The loop then returns to accepting, so that
			// multiple connections may be served concurrently.
			go handleConn(conn, srv.Handler, die)
		}
	}()
}

func handleConn(conn net.Conn, handler Handler, die chan int) {
	defer conn.Close()

	defer func() {
		if r := recover(); r != nil {
			log.Println("panic recovered in serve: ", r)
		}
	}()

	handler.Serve(conn, die)
}

// Convenience function for reading byte buffers into a Go channel.
func Pump(conn net.Conn, bufferSize int) (recv chan []byte, fault chan error) {
	recv = make(chan []byte)
	fault = make(chan error)

	go func() {
		for {
			buf := make([]byte, bufferSize)

			_, err := conn.Read(buf)
			if err != nil {
				fault <- err
				return
			}

			recv <- buf
		}
	}()

	return
}
