package main

import (
	"strings"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func parseCommand(cmd string, ael *Controller) string {
	cmd_s := strings.Split(cmd, " ")
	com := ael.FindAction(cmd_s[0])
	if com.GetName() != "nil" {
		return com.Run(ael, cmd_s[1:]...)
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

type YAMLCommand struct {
	Inputs []string `yaml:",flow"`
	Outputs []string `yaml:",flow"`
	Action string `yaml:",flow"`
}


func parseYAMLCommand(filename string) *ExternalCommand{
	y := YAMLCommand{}
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading external file %s: %v",filename,err)
		return nil
	}
	err = yaml.Unmarshal(fileData, &y)
	if err != nil {
                log.Printf("error: %v", err)
		return nil
	}
	return NewExternalCommand(y.Inputs,y.Outputs,y.Action)
}
