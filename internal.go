package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
	"log"
)

func RegisterInternal(ael *Controller) {
	ael.AddCommand(NewInternalCommand([]string{}, []string{}, "get", Get))
	ael.AddCommand(NewInternalCommand([]string{}, []string{"ip"}, "GetExternalIP", GetIP))
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
		log.Printf("TODO: Handle error %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("TODO: Handle body read error")
	}
	return string(body)
}

func Poll(ael *Controller, args ...string) string {
	if len(args) == 0 {
		return ael.GetBroadcast(0)
	} else {
		results := make([]string,len(args))
		for k,i := range args {
			bid,_ := strconv.Atoi(i)
			results[k] = strings.TrimSpace(ael.GetBroadcast(bid))
		}
		if len(results) <= 1 {
			return strings.TrimSpace(results[0])
		} else {
			return strings.Join(results,"\n")
		}
	}
}

