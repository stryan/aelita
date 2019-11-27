package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func parseCommand(cmd string) string {
	fmt.Println(cmd)
	switch cmd {
	case "get_ip":
		return getIP()
	case "close":
		return "END\n"
	case "ping":
		fmt.Println("Responding pong")
		return "pong\n"
	default:
		return "ERROR: Invalid command\n"
	}
}

func getIP() string {
	resp, err := http.Get("https://ifconfig.co")
	if err != nil {
		fmt.Println("TODO: Handle error %s",err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("TODO: Handle body read error")
	}
	return string(body)
}
