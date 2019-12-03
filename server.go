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

type Service struct {
	ch        chan bool
	waitGroup *sync.WaitGroup
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
	s.waitGroup.Wait()
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
		conn, err := listener.Accept()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}
		log.Println(conn.RemoteAddr(), "connected")
		s.waitGroup.Add(1)
		go s.serve(textproto.NewConn(conn), ael)
	}
}

func (s *Service) serve(p *textproto.Conn, ael *Controller) {
	defer p.Close()
	defer s.waitGroup.Done()
	id := p.Next()
	p.StartRequest(id)
	header, err := p.ReadLine()
	if err == io.EOF {
		return
	}
	if err != nil {
		log.Printf("reading request failed: %v\n", err)
		return
	}
	p.EndRequest(id)
	hcheck, res := checkHeader(header)
	p.StartResponse(id)
	p.PrintfLine(res)
	p.EndResponse(id)
	if hcheck == false {
		p.Close()
		return
	}
	for {
		select {
		case <-s.ch:
			CleanUpConnection(p)
			return
		default:
		}
		id := p.Next()
		p.StartRequest(id)
		cmd, err := p.ReadLine()
		p.EndRequest(id)
		if err != nil {
			log.Print(fmt.Sprintf("Error: %s", err))
			return
		}
		p.StartResponse(id)
		result := parseCommand(cmd, ael)
		if result == "END" {
			p.PrintfLine("END")
			p.EndResponse(id)
			break
		}
		p.PrintfLine(result)
		p.EndResponse(id)
	}
}

func checkHeader(header string) (bool, string) {
	headerFields := strings.Fields(string(header))
	if len(headerFields) != 2 {
		log.Print("Error: Received bad header len")
		return false, "ERR: Invalid header len"
	}
	if headerFields[0] != "aelita" {
		log.Print("Error: Received bad header server")
		return false, "ERR: Invalid header server"
	}
	if headerFields[1] != PROTOV {
		log.Print("Error: Received bad header version")
		return false, "ERR: Protocol mismatch: server accepts " + PROTOV
	}
	return true, "OK aelita " + PROTOV
}

func CleanUpConnection(p *textproto.Conn) {
	log.Print("Breaking connection")
	p.PrintfLine("END")
}

func CleanUpListener(listener *net.TCPListener) {
	log.Println("stopping listening on", listener.Addr())
	listener.Close()
}
