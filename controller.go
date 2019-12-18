package main

import (
	"strings"
	"log"
)

type Controller struct {
	inputs  map[string]Command
	outputs map[string]Command
	actions map[string]Command
	broadcasts map[int]string
	bid int
}

func NewController() *Controller {
	inputs := make(map[string]Command)
	outputs := make(map[string]Command)
	actions := make(map[string]Command)
	bc := make(map[int]string)
	c := Controller{inputs, outputs, actions, bc,1}
	return &c
}

func (c *Controller) LogCommands() {
	inputs := make([]string, len(c.inputs))
	i := 0
	for k := range c.inputs {
	    inputs[i] = k
	    i++
	}
	outputs := make([]string, len(c.outputs))
	i = 0
	for k := range c.outputs {
	    outputs[i] = k
	    i++
	}
	actions := make([]string, len(c.actions))
	i = 0
	for k := range c.actions {
	    actions[i] = k
	    i++
	}
	log.Printf("Inputs Available: %v",inputs)
	log.Printf("Outputs Available: %v",outputs)
	log.Printf("Actions Available: %v",actions)
}


func (c *Controller) AddCommand(com Command) {
	if com == nil || com.GetType() == NIL {
		return
	}
	for _, i := range com.GetInputs() {
		c.inputs[i] = com
	}
	for _, o := range com.GetOutputs() {
		c.outputs[o] = com
	}
	c.actions[com.GetName()] = com

}

func (c *Controller) FindInput(i string) Command {
	if val, ok := c.inputs[i]; ok {
		return val
	} else {
		return NewNilCommand()
	}
}

func (c *Controller) FindOutput(o string) Command {
	if val, ok := c.outputs[o]; ok {
		return val
	} else {
		return NewNilCommand()
	}
}

func (c *Controller) FindAction(a string) Command {
	if val, ok := c.actions[a]; ok {
		return val
	} else {
		return NewNilCommand()
	}
}

func (c *Controller) NewBroadcast(b string) {
	c.broadcasts[c.bid] = strings.TrimSpace(b)
	c.bid++
}

func (c *Controller) GetBroadcast(bid int) string {
	if bid > 0 {
		return c.broadcasts[bid]
	} else {
		bs := make([]string, 0, len(c.broadcasts))
		for _,v := range c.broadcasts {
			bs = append(bs,v)
		}
		results := strings.Join(bs, "\n")
		return results
	}
}
