package proxymate

import (
	"github.com/davecgh/go-spew/spew"
	"io"
	"net"
)

type SampleHandler struct{}

func (h *SampleHandler) Serve(conn net.Conn, die chan int) {
	const (
		REPLY = "pong"
	)

	recv, fault := Pump(conn, 16)

	for {
		select {
		case buf := <-recv:
			spew.Dump(buf)

			_, err := conn.Write([]byte(REPLY))
			if err != nil {
				panic(err)
			}

		case err := <-fault:
			switch err {
			// Client dropped the connection
			case io.EOF:
				goto continue_at_exit
			default:
				panic(err) // TODO: Maybe we can be more graceful?
			}

		case _, ok := <-die:
			if ok {
				panic("huh? we shouldn't be receiving anything on this channel!")
			}
			// Gracefully shutdown.
			goto continue_at_exit
		}
	}

continue_at_exit:
	return
}
