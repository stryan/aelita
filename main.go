package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	PROTOV = "0.1"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/aelita/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.SetConfigType("yaml")
	port := viper.GetString("port")

	Listen(":" + port)
}

func Listen(addr string) {
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalf("Listen failed: %v\n", err)
		os.Exit(1)
	}
	defer ln.Close()
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Print("Error: Failed to accept connection")
		}
		go handleConnection(textproto.NewConn(c))
	}
}

func handleConnection(p *textproto.Conn) {
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
		id := p.Next()
		p.StartRequest(id)
		cmd, err := p.ReadLine()
		p.EndRequest(id)
		if err != nil {
			log.Print(fmt.Sprintf("Error: %s", err))
			return
		}
		p.StartResponse(id)
		result := parseCommand(cmd)
		if result == "END" {
			p.PrintfLine("END")
			p.EndResponse(id)
			break
		}
		p.PrintfLine(result)
		p.EndResponse(id)

	}
	p.Close()
	return
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
