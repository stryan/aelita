package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/chewxy/sexp"
	"gopkg.in/yaml.v2"
)

func parseCommand(cmd []sexp.Sexp, ael *Controller) sexp.Sexp {
	if len(cmd) <= 0 {
		log.Printf("%v %v\n", len(cmd), cmd)
		panic("Bad command")
	}
	switch fmt.Sprint(cmd[0].Head()) {
	case "CMD":
		scmd := cmd[0].Tail().Head()
		var scmd_args sexp.List
		if scmd.Tail() != nil {
			scmd_args = scmd.Tail().(sexp.List)
		} else {
			scmd_args = sexp.List{}
		}
		com := ael.FindAction(fmt.Sprint(scmd.Head()))
		if com.GetName() != "nil" {
			return com.Run(ael, scmd_args)
		}
		return newErr("Invalid Command")
	case "DAT":
		return newErr("Didn't ask for data")
	case "NEW":
		if ParseHeader(cmd[0]) == true {
			res, _ := sexp.SymbolReader("ACTIVE")
			return sexp.List{res}
		} else {
			return newErr("Bad header")
		}
	case "ERR":
		return newErr("Only aelita can send errors")
	case "END":
		return newEnd()
	case "ACK":
		return newEmpty()
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

func newErr(msg string) sexp.Sexp {
	res, _ := sexp.SymbolReader(msg)
	return sexp.List{sexp.Symbol("ERR"), res}
}

func newAck() sexp.Sexp {
	res, _ := sexp.SymbolReader("ACK")
	return sexp.List{res}
}

func newEmpty() sexp.Sexp {
	res, _ := sexp.SymbolReader("")
	return sexp.List{res}
}

func newEnd() sexp.Sexp {
	res, _ := sexp.SymbolReader("END")
	return sexp.List{res}
}
