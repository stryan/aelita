package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/chewxy/sexp"
	"gopkg.in/yaml.v2"
)

func parseCommand(cmd []sexp.Sexp, ael *Controller) string {
	if len(cmd) <= 0 {
		log.Printf("%v %v\n", len(cmd), cmd)
		panic("Bad command")
	}
	switch fmt.Sprint(cmd[0].Head()) {
	case "CMD":
		scmd := cmd[0].Tail().Head()
		scmd_args := scmd.Tail()
		com := ael.FindAction(fmt.Sprint(scmd.Head()))
		if com.GetName() != "nil" {
			return com.Run(ael, scmd_args)
		}
		return "(ERR \"Invalid Command\")"
	case "DAT":
		return "(ERR \"Didn't ask for data\")"
	case "NEW":
		if ParseHeader(cmd[0]) == true {
			return "ACTIVE"
		} else {
			return "(ERR \"Bad header\")"
		}
	case "ERR":
		return "(ERR \"Only aelita can send errors\")"
	case "END":
		return "(END)"
	case "ACK":
		return ""
	}
	panic("Failed to parse properly!")
}

type YAMLCommand struct {
	Inputs  []string `yaml:",flow"`
	Outputs []string `yaml:",flow"`
	Action  string   `yaml:",flow"`
}

func parseYAMLCommand(filename string) *ExternalCommand {
	y := YAMLCommand{}
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading external file %s: %v", filename, err)
		return nil
	}
	err = yaml.Unmarshal(fileData, &y)
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}
	return NewExternalCommand(y.Inputs, y.Outputs, y.Action)
}
