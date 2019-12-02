package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func parseCommand(cmd string) string {
	switch cmd {
	case "get_ip":
		return getIP()
	case "close":
		return "END"
	case "ping":
		return "pong"
	default:
		return "ERROR: Invalid command"
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
