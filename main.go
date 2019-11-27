package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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

	l, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		panic(fmt.Errorf("Failed to start listening server: %s \n", err))
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Print("Error: Failed to accept connection")
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	log.Print("Serving %s\n", c.RemoteAddr().String())
	header, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		log.Print("Error: %s", err)
	}
	headerFields := strings.Fields(string(header))
	if len(headerFields) != 2 {
		log.Print("Error: Received bad header len")
		c.Write([]byte("ERR: Invalid header len"))
		c.Close()
		return
	}
	if headerFields[0] != "aelita" {
		log.Print("Error: Received bad header server")
		c.Write([]byte("ERR: Invalid header server"))
		c.Close()
		return
	}
	if headerFields[1] != PROTOV {
		log.Print("Error: Received bad header version")
		msg := fmt.Sprintf("ERR: Protocol mismatch: server accepts %s", PROTOV)
		c.Write([]byte(msg))
		c.Close()
		return
	}
	msg := fmt.Sprintf("OK aelita %s\n", PROTOV)
	c.Write([]byte(msg))
	for {
		cmd, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Print(fmt.Sprintf("Error: %s", err))
			return
		}
		temp := strings.TrimSpace(string(cmd))
		result := parseCommand(temp)
		if result == "END" {
			break
		}
		c.Write([]byte(result))
	}
	c.Write([]byte("END\n"))
	c.Close()
	return
}
