package main

import (
	"bytes"
	"os"
	"os/exec"
	"log"
)

type CommandAction func(*Controller, ...string) string
type CommandType int

const (
	INTERNAL = 0
	EXTERNAL = 1
	NIL      = 2
)

type Command interface {
	GetType() CommandType
	GetInputs() []string
	GetOutputs() []string
	GetName() string
	Run(*Controller, ...string) string
}

type InternalCommand struct {
	Inputs     []string
	Outputs    []string
	ActionName string
	Action     CommandAction
}

type ExternalCommand struct {
	Inputs     []string
	Outputs    []string
	ActionName string
}

type NilCommand struct {
	Inputs     []string
	Outputs    []string
	ActionName string
}

func NewNilCommand() *NilCommand {
	return &NilCommand{[]string{}, []string{}, "nil"}
}

func NewInternalCommand(inputs []string, outputs []string, name string, action CommandAction) *InternalCommand {
	return &InternalCommand{inputs, outputs, name, action}
}

func NewExternalCommand(inputs []string, outputs []string, name string) *ExternalCommand {
	return &ExternalCommand{inputs, outputs, name}
}

func (i *InternalCommand) GetType() CommandType {
	return INTERNAL
}

func (e *ExternalCommand) GetType() CommandType {
	return EXTERNAL
}

func (n *NilCommand) GetType() CommandType {
	return NIL
}

func (i *InternalCommand) GetInputs() []string {
	return i.Inputs
}

func (e *ExternalCommand) GetInputs() []string {
	return e.Inputs
}

func (n *NilCommand) GetInputs() []string {
	return n.Inputs
}

func (i *InternalCommand) GetOutputs() []string {
	return i.Outputs
}

func (e *ExternalCommand) GetOutputs() []string {
	return e.Outputs
}

func (n *NilCommand) GetOutputs() []string {
	return n.Outputs
}

func (i *InternalCommand) GetName() string {
	return i.ActionName
}

func (e *ExternalCommand) GetName() string {
	return e.ActionName
}

func (n *NilCommand) GetName() string {
	return n.ActionName
}

func (i *InternalCommand) Run(ael *Controller, args ...string) string {
	return i.Action(ael,args...)
}

func (e *ExternalCommand) Run(ael *Controller, args ...string) string {
	cmd := exec.Command(e.ActionName, args...)
	cmd.Env = os.Environ()
	var out,cmdErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &cmdErr
	err := cmd.Run()
	if err != nil {
		log.Printf("Unable to run external command: '%v'", err)
		log.Printf("Command '%v' error: %v",e.ActionName,cmdErr.String())
		return "nil"
	}
	return out.String()
}

func (n *NilCommand) Run(ael *Controller, args ...string) string {
	return "nil"
}
