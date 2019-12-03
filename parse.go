package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func parseCommand(cmd string, ael *Controller) string {
	cmd_s := strings.Split(cmd," ")
	com := ael.FindAction(cmd_s[0])
	if com.GetName() != "nil" {
		return com.Run(ael,cmd_s[1:]...)
	}
	// Not in action list, check built in
	switch cmd_s[0] {
	case "close":
		return "END"
	case "ping":
		return "pong"
	default:
		return "ERROR: Invalid command: " + cmd_s[0]
	}
}

func GetIP(ael *Controller, args ...string) string {
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
