package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type AelConn struct {
	C net.Conn
	P *textproto.Conn
	Closed bool
}

func NewAelConn(c net.Conn, p *textproto.Conn) *AelConn {
	a := &AelConn{
		C: c,
		P: p,
		Closed: false,
	}
	return a
}

func (a *AelConn) Close() {
	if a.Closed == false {
		a.C.Close()
		a.Closed = true
		return
	} else {
		return
	}
}

type Service struct {
	ch        chan bool
	waitGroup *sync.WaitGroup
	openConns []AelConn
}

func NewService() *Service {
	s := &Service{
		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
	s.waitGroup.Add(1)
	return s
}

func (s *Service) Stop() {
	close(s.ch)
	log.Println("Stopping Service")
	//Try to wait out
	if waitTimeout(s.waitGroup,10 * time.Second) {
		//timed out, kill them all
		log.Print("Closing remaining connections")
		for _,c := range s.openConns {
			c.P.PrintfLine("END aelita " + PROTOV)
			c.Close()
		}
	} else {
		//Everyone's happy!
		log.Print("All connections closed normally")
	}
	log.Println("Service Stopped")
}
//https://stackoverflow.com/a/32843750
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
    c := make(chan struct{})
    go func() {
        defer close(c)
        wg.Wait()
    }()
    select {
    case <-c:
        return false // completed normally
    case <-time.After(timeout):
        return true // timed out
    }
}


func StartServer(addr string, ael *Controller) {
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if nil != err {
		log.Fatalln(err)
	}
	ln, err := net.ListenTCP("tcp",laddr)
	if err != nil {
		log.Fatalf("Listen failed: %v\n", err)
		os.Exit(1)
	}
	service := NewService()
	go service.Serve(ln, ael)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	// Stop the service gracefully.
	service.Stop()
}

func (s *Service) Serve(listener *net.TCPListener, ael *Controller) {
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			CleanUpListener(listener)
			return
		default:
		}
		listener.SetDeadline(time.Now().Add(1e9))
		conn, err := listener.AcceptTCP()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}
		log.Println(conn.RemoteAddr(), "connected")
		a := NewAelConn(conn,textproto.NewConn(conn))
		s.waitGroup.Add(1)
		go s.serve(a, ael)
	}
}

func (s *Service) serve(ac *AelConn, ael *Controller) {
	defer ac.Close()
	defer s.waitGroup.Done()
	//conn.SetDeadline(time.Now().Add(1e9))
	id := ac.P.Next()
	ac.P.StartRequest(id)
	header, err := ac.P.ReadLine()
	if err == io.EOF {
		return
	}
	if err != nil {
		log.Printf("reading request failed: %v\n", err)
		return
	}
	ac.P.EndRequest(id)
	hcheck, res := checkHeader(header)
	ac.P.StartResponse(id)
	ac.P.PrintfLine(res)
	ac.P.EndResponse(id)
	if hcheck == false {
		ac.Close()
		return
	}
	for {
		select {
		case <-s.ch:
			CleanUpConnection(ac)
			return
		default:
		}
		id := ac.P.Next()
		ac.P.StartRequest(id)
		cmd, err := ac.P.ReadLine()
		ac.P.EndRequest(id)
		if err != nil {
			log.Print(fmt.Sprintf("Error: %s", err))
			return
		}
		ac.P.StartResponse(id)
		result := parseCommand(cmd[4:], ael)
		if result == "END" {
			ac.P.PrintfLine("END aelita " + PROTOV)
			ac.P.EndResponse(id)
			break
		}
		d := 1+strings.Count(result,"\n")
		ac.P.PrintfLine("DAT %v",d)
		ac.P.PrintfLine(result)
		ac.P.EndResponse(id)
	}
}

func checkHeader(header string) (bool, string) {
	headerFields := strings.Fields(string(header))
	if len(headerFields) != 3 {
		log.Print("Error: Received bad header len")
		return false, "ERR: Invalid header len"
	}
	if headerFields[0] != "NEW" {
		log.Print("Error: Not new connection")
		return false, "ERR: Not new connection"
	}
	if headerFields[1] != "aelita" {
		log.Print("Error: Received bad header server")
		return false, "ERR: Invalid header server"
	}
	if headerFields[2] != PROTOV {
		log.Print("Error: Received bad header version")
		return false, "ERR: Protocol mismatch: server accepts " + PROTOV
	}
	return true, "OK aelita " + PROTOV
}

func CleanUpConnection(a *AelConn) {
	log.Print("Breaking connection")
	a.P.PrintfLine("END aelita " + PROTOV)
}

func CleanUpListener(listener *net.TCPListener) {
	log.Println("stopping listening on", listener.Addr())
	listener.Close()
}
