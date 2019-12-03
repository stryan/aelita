package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func RegisterInternal(ael *Controller) {
	get_c := NewInternalCommand([]string{}, []string{}, "get", Get)
	get_ip := NewInternalCommand([]string{}, []string{"ip"}, "GetExternalIP", GetIP)
	ael.AddCommand(get_c)
	ael.AddCommand(get_ip)
}

func Get(ael *Controller, args ...string) string {
	if len(args) == 0 {
		return "No arguments provided"
	}
	var output []string
	for _, a := range args {
		get_command := ael.FindOutput(a)
		if get_command.GetName() == "nil" {
			return "N/A"
		}
		output = append(output, strings.TrimSpace(get_command.Run(ael, args[1:]...)))
	}
	result := strings.Join(output, "; ")
	return result
}

func GetIP(ael *Controller, args ...string) string {
	resp, err := http.Get("https://ifconfig.co")
	if err != nil {
		fmt.Println("TODO: Handle error %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("TODO: Handle body read error")
	}
	return string(body)
}
