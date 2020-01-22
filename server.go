package main

import (
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/chewxy/sexp"
)

type AelConn struct {
	C      net.Conn
	P      *textproto.Conn
	Active bool
	Closed bool
}

func NewAelConn(c net.Conn, p *textproto.Conn) *AelConn {
	a := &AelConn{
		C:      c,
		P:      p,
		Active: false,
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
	openConns []*AelConn
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
	if waitTimeout(s.waitGroup, 10*time.Second) {
		//timed out, kill them all
		log.Print("Closing remaining connections")
		for _, c := range s.openConns {
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
	ln, err := net.ListenTCP("tcp", laddr)
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
		a := NewAelConn(conn, textproto.NewConn(conn))
		s.waitGroup.Add(1)
		s.openConns = append(s.openConns, a)
		go s.serve(a, ael)
	}
}

func (s *Service) serve(ac *AelConn, ael *Controller) {
	defer ac.Close()
	defer s.waitGroup.Done()
	//conn.SetDeadline(time.Now().Add(1e9))
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
		sexp, sres := sexp.ParseString(cmd)
		if sres != nil {
			ac.P.PrintfLine(fmt.Sprint(newErr("Invalid s-expression")))
			ac.P.EndResponse(id)
			break
		}
		result := parseCommand(sexp, ael)
		output_result := fmt.Sprint(result)
		log.Println(output_result)
		if output_result == "(END)" {
			ac.P.PrintfLine("(ACK)")
			ac.P.EndResponse(id)
			break
		}
		if output_result == "(ACTIVE)" {
			ac.Active = true
			ac.P.PrintfLine("(ACK)")
			ac.P.EndResponse(id)
			continue
		}
		if ac.Active {
			ac.P.PrintfLine(output_result)
			ac.P.EndResponse(id)
		} else {
			ac.P.PrintfLine(fmt.Sprint(newErr("No header exchange")))
			ac.P.EndResponse(id)
		}
	}
}

// Example: (NEW (aelita 0.2))
func checkHeader(head sexp.List) (bool, sexp.Sexp) {
	bad_header := "Bad header"

	if len(head) <= 1 {
		log.Printf("Error 0: %v\n", len(head))
		return false, newErr(bad_header)
	}
	header := sexp.List(head)
	log.Printf("len: %v\n value: %v\n", len(header), header)
	if header.LeafCount() < 3 || header.Head().Head() != sexp.Symbol("NEW") {
		log.Println("Error 1")
		return false, newErr(bad_header)
	}
	tail := header.Tail()
	fmt.Printf("%v \n", tail)
	if tail.LeafCount() < 2 || tail.Head().Head() != sexp.Symbol("aelita") {
		log.Printf("Error 2")
		return false, newErr(bad_header)
	}
	if tail.LeafCount() < 2 || tail.Head().Tail().Head() != sexp.Symbol(PROTOV) {
		log.Printf("Error 3")
		return false, newErr("Protocol version mismatch")
	}
	return true, newAck()
}

func CleanUpConnection(a *AelConn) {
	log.Print("Breaking connection")
	a.P.PrintfLine("(END)")
}

func CleanUpListener(listener *net.TCPListener) {
	log.Println("stopping listening on", listener.Addr())
	listener.Close()
}
