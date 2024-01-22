package internal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	jdi "github.com/kyo-w/jdwp"
	"io"
	"log"
	"reflect"
	"sync"
	"time"
)

const cmdCompositeEvent = cmdID(100)

var (
	handshake = []byte("JDWP-Handshake")

	defaultIDSizes = jdi.IDSizes{
		FieldIDSize:         8,
		MethodIDSize:        8,
		ObjectIDSize:        8,
		ReferenceTypeIDSize: 8,
		FrameIDSize:         8,
	}
)

type Connection struct {
	in           io.Reader
	r            Reader
	w            Writer
	flush        func() error
	idSizes      jdi.IDSizes
	nextPacketID packetID
	Events       map[jdi.EventRequestID]chan<- jdi.EventResponse
	// 这与JDWP通信包相关，每一个包都有一个ID表示，发送包时自行指定，响应时自行从映射中获取
	replies map[packetID]chan<- replyPacket
	sync.Mutex
}

func Open(ctx context.Context, conn io.ReadWriteCloser) (*Connection, error) {
	if err := exchangeHandshakes(conn); err != nil {
		return nil, err
	}
	buf := bufio.NewWriterSize(conn, 1024)
	r := ByteOrderReader(conn, BigEndian)
	w := ByteOrderWriter(buf, BigEndian)
	c := &Connection{
		in:      conn,
		r:       r,
		w:       w,
		flush:   buf.Flush,
		idSizes: defaultIDSizes,
		Events:  map[jdi.EventRequestID]chan<- jdi.EventResponse{},
		replies: map[packetID]chan<- replyPacket{},
	}

	go c.recv(ctx)
	var err error
	c.idSizes, err = c.GetIDSizes()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func exchangeHandshakes(conn io.ReadWriter) error {
	if _, err := conn.Write(handshake); err != nil {
		return err
	}
	ok, err := expect(conn, handshake)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("bad handshake")
	}
	return nil
}
func expect(conn io.Reader, expected []byte) (bool, error) {
	got := make([]byte, len(expected))
	for len(expected) > 0 {
		n, err := conn.Read(got)
		if err != nil {
			return false, err
		}
		for i := 0; i < n; i++ {
			if got[i] != expected[i] {
				return false, nil
			}
		}
		got, expected = got[n:], expected[n:]
	}
	return true, nil
}
func (c *Connection) SendCommand(cmd Cmd, req interface{}, out interface{}) error {
	p, err := c.req(cmd, req)
	if err != nil {
		return err
	}
	return p.wait(out)
}

func (c *Connection) req(cmd Cmd, req interface{}) (*pending, error) {
	data := bytes.Buffer{}
	if req != nil {
		e := ByteOrderWriter(&data, BigEndian)
		if err := c.encode(e, reflect.ValueOf(req)); err != nil {
			return nil, err
		}
	}

	id, replyChan := c.newReplyHandler()

	p := cmdPacket{id: id, cmdSet: cmd.set, cmdID: cmd.id, data: data.Bytes()}

	c.Lock()
	defer c.Unlock()

	if err := p.write(c.w); err != nil {
		return nil, err
	}
	if err := c.flush(); err != nil {
		return nil, err
	}
	return &pending{c, replyChan, id}, nil
}

type pending struct {
	c  *Connection
	p  <-chan replyPacket
	id packetID
}

func (p *pending) wait(out interface{}) error {
	select {
	case reply := <-p.p:
		if reply.err != ErrNone {
			log.Printf("<%v> recv err: %+v", p.id, reply.err)
			fmt.Printf("<%v> recv err: %+v", p.id, reply.err)
			return reply.err
		}
		if out == nil {
			return nil
		}
		r := bytes.NewReader(reply.data)
		d := ByteOrderReader(r, BigEndian)
		if err := p.c.decode(d, reflect.ValueOf(out)); err != nil {
			return err
		}
		if offset, _ := r.Seek(0, 1); offset != int64(len(reply.data)) {
			panic(fmt.Errorf("Only %d/%d bytes read from reply packet", offset, len(reply.data)))
		}
		return nil
	case <-time.After(time.Second * 120):
		return fmt.Errorf("timeout")
	}
}
func (c *Connection) newReplyHandler() (packetID, <-chan replyPacket) {
	reply := make(chan replyPacket, 1)
	c.Lock()
	id := c.nextPacketID
	c.nextPacketID++
	c.replies[id] = reply
	c.Unlock()
	return id, reply
}
func (c *Connection) GetIDSizes() (jdi.IDSizes, error) {
	res := jdi.IDSizes{}
	err := c.SendCommand(CmdVirtualMachineIDSizes, struct{}{}, &res)
	return res, err
}
func (c *Connection) recv(ctx context.Context) {
	for !Stopped(ctx) {
		packet, err := c.readPacket()
		switch err {
		case nil:
		case io.EOF:
			return
		default:
			if !Stopped(ctx) {
				// TODO: turn it into a log
				fmt.Printf("Failed to read packet. Error: %v\n", err)
			}
			return
		}

		switch packet := packet.(type) {
		case replyPacket:
			c.Lock()
			out, ok := c.replies[packet.id]
			delete(c.replies, packet.id)
			c.Unlock()
			if !ok {
				// TODO: turn it into a log
				fmt.Printf("Unexpected reply for packet %d\n", packet.id)
				continue
			}
			out <- packet

		case cmdPacket:
			switch {
			case packet.cmdSet == cmdSetEvent && packet.cmdID == cmdCompositeEvent:
				d := ByteOrderReader(bytes.NewReader(packet.data), BigEndian)
				l := jdi.EventsResponse{}
				if err := c.decode(d, reflect.ValueOf(&l)); err != nil {
					// TODO: turn it into a log
					fmt.Printf("Couldn't decode composite event data. Error: %v\n", err)
					continue
				}

				for _, ev := range l.Events {
					//fmt.Printf("<%v> event: %T %+v", ev.request(), ev, ev)

					c.Lock()
					handler, ok := c.Events[ev.GetRequest()]
					c.Unlock()

					if ok {
						handler <- ev
					} else {
						fmt.Printf("No event handler registered for %+v", ev)
					}
				}

			default:
				fmt.Printf("received unknown packet %+v", packet)
			}
		}
	}
}
