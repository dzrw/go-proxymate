package proxymate

import (
	"github.com/davecgh/go-spew/spew"
	"log"
	"net"
	"time"
)

type Client struct {
	Addr string
	Rate int

	conn net.Conn
	quit chan chan bool
}

func DialAndPing(addr string, rate int) (c *Client, err error) {
	c, err = Dial(addr)
	if err != nil {
		return
	}

	c.Ping(rate)
	return
}

func Dial(addr string) (c *Client, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}

	c = &Client{
		Addr: addr,
		conn: conn,
		quit: make(chan chan bool),
	}

	return
}

func (c *Client) Ping(rate int) {
	go func() {
		const (
			REQUEST = "ping"
		)

		defer c.conn.Close()

		recv, fault := Pump(c.conn, 32)

		for {
			select {
			case buf := <-recv:
				spew.Dump(buf)

			case err := <-fault:
				log.Println("client read error: ", err)
				break

			case <-time.After(time.Duration(rate) * time.Second):
				_, err := c.conn.Write([]byte(REQUEST))
				if err != nil {
					log.Println("client write error: ", err)
					break
				}

			case done := <-c.quit:
				done <- true
				break
			}
		}
	}()
}

func (c *Client) Stop() {
	log.Println("stopping client")
	done := make(chan bool)
	c.quit <- done
	<-done
}
