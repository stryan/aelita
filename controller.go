package main

type Controller struct {
	inputs  map[string]Command
	outputs map[string]Command
	actions map[string]Command
}

func NewController() *Controller {
	inputs := make(map[string]Command)
	outputs := make(map[string]Command)
	actions := make(map[string]Command)
	c := Controller{inputs, outputs, actions}
	return &c
}

func (c *Controller) AddCommand(com Command) {
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
